/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package socket

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"golang.org/x/crypto/ssh"

	"github.com/NethServer/ns8-core/core/api-server/redis"
)

const terminalHandshakeTimeout = 15 * time.Second

var (
	errTerminalControl        = errors.New("terminal: malformed control frame")
	errTerminalAddressChanged = errors.New("the browser address changed since the terminal was requested, open it again")
)

// sshUsernamePattern matches cluster/validator-definitions.json
// strict-username-string, the only shape cluster/actions/add-user accepts.
var sshUsernamePattern = regexp.MustCompile(`^[a-z_][a-z0-9_-]*\$?$`)

func errTerminalBusy(holder string) error {
	return fmt.Errorf("a terminal is already open on this node by %s", holder)
}

func errTerminalTicket(reason string) error {
	return errors.New(reason)
}

type terminalState struct {
	mu     sync.Mutex
	phase  string
	opened time.Time

	nodeID   string
	holderID uint64
	user     string
	sshUser  string
	clientIP string
	exp      int64

	// Kept outside the SSH session: the browser sizes the pane as soon as it is
	// mounted, which is now before the shell exists.
	rows int
	cols int

	// Keystrokes waiting to be consumed by the login prompts. Buffered so a
	// paste does not block melody's read pump.
	loginInput chan []byte

	// What followed the newline in the frame the previous prompt stopped on. A
	// paste can carry the user name and the password in one frame, and those
	// trailing bytes belong to the next prompt.
	loginRemainder []byte

	// Closed once by closeTerminal. Without it a login goroutine would sit on
	// its prompt until the login deadline after the browser has gone, and could
	// still bring up a shell on a session nobody is reading.
	done chan struct{}

	client    *ssh.Client
	session   *ssh.Session
	stdin     io.WriteCloser
	localPort string
	lastIO    time.Time
	closed    bool
}

func (state *terminalState) currentPhase() string {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.phase
}

func (state *terminalState) setPhase(phase string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.phase = phase
}

// adopt binds a consumed ticket to the session and enters the login phase in a
// single step. It reports false when the session was closed meanwhile: melody
// runs closeTerminal from its write pump, so a close can land between consuming
// the ticket and installing it. Installing it anyway would revive a closed
// session whose later closeTerminal calls all return at the closed guard,
// leaving the node reserved for good.
func (state *terminalState) adopt(ticket *terminalTicket) bool {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.closed {
		return false
	}
	state.nodeID = ticket.nodeID
	state.holderID = ticket.holderID
	state.user = ticket.user
	state.clientIP = ticket.clientIP
	state.exp = ticket.exp
	state.phase = phaseLogin
	return true
}

func (state *terminalState) touch() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.lastIO = time.Now()
}

func (state *terminalState) idleFor() time.Duration {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.lastIO.IsZero() {
		return 0
	}
	return time.Since(state.lastIO)
}

func (state *terminalState) write(data []byte) {
	state.mu.Lock()
	stdin := state.stdin
	state.lastIO = time.Now()
	state.mu.Unlock()
	if stdin == nil {
		return
	}
	_, _ = stdin.Write(data)
}

// resize records the size in every phase and forwards it only once a shell
// exists: the pane is mounted for the login prompts, so the first frames land
// before there is any session to notify.
func (state *terminalState) resize(rows int, cols int) {
	if rows <= 0 || cols <= 0 {
		return
	}
	state.mu.Lock()
	state.rows = rows
	state.cols = cols
	session := state.session
	state.mu.Unlock()
	if session == nil {
		return
	}
	_ = session.WindowChange(rows, cols)
}

func (state *terminalState) currentSize() (int, int) {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.rows, state.cols
}

func (state *terminalState) takeLoginRemainder() []byte {
	state.mu.Lock()
	defer state.mu.Unlock()
	rest := state.loginRemainder
	state.loginRemainder = nil
	return rest
}

func (state *terminalState) stashLoginRemainder(data []byte) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.loginRemainder = data
}

// dropLoginRemainder throws away buffered input after a failed attempt. A paste
// that went wrong must not spill its second line into the next prompt, where
// the user name echoes and a password would land in the scrollback.
func (state *terminalState) dropLoginRemainder() {
	state.mu.Lock()
	defer state.mu.Unlock()
	for i := range state.loginRemainder {
		state.loginRemainder[i] = 0
	}
	state.loginRemainder = nil
}

// offerLoginInput queues keystrokes for the login prompts. It reports false
// when the queue is saturated, which no human typing can do.
func (state *terminalState) offerLoginInput(keys []byte) bool {
	select {
	case state.loginInput <- keys:
		return true
	default:
		return false
	}
}

/*
 * Handshake
 */

// dialNode authenticates against a node and returns the live client. It neither
// closes the WebSocket nor records anything: only the caller can tell a refused
// password from an operator who pressed Ctrl-C at the prompt.
func dialNode(target string, hostKeys []ssh.PublicKey, sshUser string, ask ssh.KeyboardInteractiveChallenge) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: sshUser,
		// One method only: x/crypto/ssh tries every method listed and each
		// failure counts against the server's MaxAuthTries. Every retry is a
		// fresh Dial, so each one starts from a clean budget on the node.
		Auth:              []ssh.AuthMethod{ssh.KeyboardInteractive(ask)},
		HostKeyCallback:   publishedHostKeyCallback(hostKeys),
		HostKeyAlgorithms: publishedHostKeyAlgorithms(hostKeys),
		Timeout:           terminalHandshakeTimeout,
	}

	return ssh.Dial("tcp", target, config)
}

// attachShell takes an authenticated client and puts the session in service.
func attachShell(s *melody.Session, state *terminalState, sshUser string, client *ssh.Client) {
	// The browser can vanish while the handshake is in flight. Opening a shell
	// then would resurrect a closed session and leak the SSH client.
	if state.currentPhase() != phaseLogin {
		_ = client.Close()
		return
	}

	rows, cols := state.currentSize()

	localPort := ""
	if _, port, splitErr := net.SplitHostPort(client.LocalAddr().String()); splitErr == nil {
		localPort = port
	}

	session, stdout, stdin, err := openTerminalShell(client, rows, cols)
	if err != nil {
		_ = client.Close()
		writeTerminalControl(s, gin.H{"type": "auth-error", "message": "could not open a shell: " + err.Error()})
		closeTerminal(s, "shell refused")
		return
	}

	// Best effort, and after the shell is up: the operating system
	// administrator has no access to audit.db, so this is the only place the
	// cluster-admin identity shows up on the node itself, but it must not delay
	// a session that is otherwise ready.
	correlated := writeCorrelationRecord(client, state.user, state.clientIP) == nil

	state.mu.Lock()
	state.sshUser = sshUser
	state.client = client
	state.session = session
	state.stdin = stdin
	state.localPort = localPort
	state.lastIO = time.Now()
	state.phase = phaseRunning
	state.mu.Unlock()

	auditTerminal(state.user, "terminal-open", gin.H{
		"node":       state.nodeID,
		"ssh_user":   sshUser,
		"client_ip":  state.clientIP,
		"local_port": localPort,
		"correlated": correlated,
	})

	writeTerminalControl(s, gin.H{"type": "ready"})

	go pumpTerminalOutput(s, state, stdout)
	go terminalKeepaliveLoop(s, state)
}

func openTerminalShell(client *ssh.Client, rows int, cols int) (*ssh.Session, io.Reader, io.WriteCloser, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, nil, err
	}

	if rows <= 0 {
		rows = 24
	}
	if cols <= 0 {
		cols = 80
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 38400,
		ssh.TTY_OP_OSPEED: 38400,
	}
	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		_ = session.Close()
		return nil, nil, nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, nil, err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, nil, err
	}
	if err := session.Shell(); err != nil {
		_ = session.Close()
		return nil, nil, nil, err
	}
	return session, stdout, stdin, nil
}

/*
 * Target resolution. The browser only ever sends a node identifier: taking a
 * host and port from the client would turn api-server into an arbitrary SSH
 * proxy.
 */

type publishedHostKey struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

func nodeSSHTarget(nodeID string) (string, []ssh.PublicKey, error) {
	redisConnection := redis.Instance()
	defer redisConnection.Close()

	address, err := redisConnection.HGet(terminalCtx, "node/"+nodeID+"/vpn", "ip_address").Result()
	if err != nil || address == "" {
		return "", nil, fmt.Errorf("no VPN address published for node/%s", nodeID)
	}

	port, _ := redisConnection.HGet(terminalCtx, "node/"+nodeID+"/ssh", "port").Result()
	if port == "" {
		port = "22"
	}

	raw, _ := redisConnection.HGet(terminalCtx, "node/"+nodeID+"/ssh", "host_keys").Result()
	keys, err := parsePublishedHostKeys(raw)
	if err != nil {
		return "", nil, err
	}

	return net.JoinHostPort(address, port), keys, nil
}

func parsePublishedHostKeys(raw string) ([]ssh.PublicKey, error) {
	if raw == "" {
		return nil, errors.New("this node published no SSH host key, run probe-terminal-access on it")
	}

	var published []publishedHostKey
	if err := json.Unmarshal([]byte(raw), &published); err != nil {
		return nil, errors.New("the SSH host keys published by this node are unreadable")
	}

	var keys []ssh.PublicKey
	for _, entry := range published {
		parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(entry.Type + " " + entry.Key))
		if err != nil {
			continue
		}
		keys = append(keys, parsed)
	}
	if len(keys) == 0 {
		return nil, errors.New("this node published no usable SSH host key")
	}
	return keys, nil
}

/*
 * publishedHostKeyCallback accepts any key the node published over the
 * cluster's authenticated channel, Redis ACL plus WireGuard. This is not a
 * pinning: a legitimately regenerated key is accepted as soon as it is
 * republished, and so would a substituted one if an attacker controlled the
 * node -- a case where the node is lost anyway.
 */
func publishedHostKeyCallback(keys []ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, offered ssh.PublicKey) error {
		offeredBytes := offered.Marshal()
		for _, known := range keys {
			if subtle.ConstantTimeCompare(offeredBytes, known.Marshal()) == 1 {
				return nil
			}
		}
		return errors.New("the SSH host key does not match the keys published by this node")
	}
}

// publishedHostKeyAlgorithms restricts negotiation to the published key types.
// Without it, a server offering RSA before a published Ed25519 key produces a
// refusal that looks like a host key mismatch.
func publishedHostKeyAlgorithms(keys []ssh.PublicKey) []string {
	var algorithms []string
	seen := map[string]bool{}
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			algorithms = append(algorithms, name)
		}
	}
	for _, key := range keys {
		if key.Type() == ssh.KeyAlgoRSA {
			// An ssh-rsa host key also serves the SHA-2 signature algorithms,
			// and modern servers refuse the SHA-1 one.
			add(ssh.KeyAlgoRSASHA512)
			add(ssh.KeyAlgoRSASHA256)
		}
		add(key.Type())
	}
	return algorithms
}

/*
 * Relay
 */

func pumpTerminalOutput(s *melody.Session, state *terminalState, stdout io.Reader) {
	chunks := make(chan []byte, 64)

	go func() {
		defer close(chunks)
		buffer := make([]byte, 32<<10)
		for {
			read, err := stdout.Read(buffer)
			if read > 0 {
				chunk := make([]byte, read)
				copy(chunk, buffer[:read])
				chunks <- chunk
			}
			if err != nil {
				return
			}
		}
	}()

	var pending []byte
	flush := func() {
		if len(pending) == 0 {
			return
		}
		_ = s.WriteBinary(pending)
		pending = nil
	}

	// Coalesce so that a burst of small reads does not fill melody's outbound
	// buffer, which drops frames rather than blocking.
	ticker := time.NewTicker(terminalCoalesceWindow)
	defer ticker.Stop()

	for {
		select {
		case chunk, open := <-chunks:
			if !open {
				flush()
				closeTerminal(s, "remote shell closed")
				return
			}
			pending = append(pending, chunk...)
			state.touch()
			if len(pending) >= terminalMaxMessageSize/2 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func terminalKeepaliveLoop(s *melody.Session, state *terminalState) {
	ticker := time.NewTicker(terminalKeepalive)
	defer ticker.Stop()

	for range ticker.C {
		if state.currentPhase() != phaseRunning {
			return
		}

		state.mu.Lock()
		client := state.client
		exp := state.exp
		nodeID := state.nodeID
		opened := state.opened
		state.mu.Unlock()

		// x/crypto/ssh has no keepalive of its own. Bounded because it opens
		// this loop, and the loop is what enforces the expiry, idle, duration
		// and disabled checks below: a black-holed node would otherwise leave a
		// live shell that none of them can ever reach.
		alive := make(chan error, 1)
		go func() {
			_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
			alive <- err
		}()

		select {
		case err := <-alive:
			if err != nil {
				closeTerminal(s, "lost the connection to the node")
				return
			}
		case <-time.After(terminalKeepaliveTimeout):
			closeTerminal(s, "the node stopped answering")
			return
		}

		if exp > 0 && time.Now().Unix() > exp {
			closeTerminal(s, "session token expired")
			return
		}
		if state.idleFor() > terminalIdleTimeout {
			closeTerminal(s, "idle for too long")
			return
		}
		if time.Since(opened) > terminalMaxDuration {
			closeTerminal(s, "maximum session duration reached")
			return
		}
		// Re-read the flag so that disable-node-terminal, a cluster restore or
		// a manual change closes live sessions. A cluster action cannot reach
		// these connections, which live in this process.
		if terminalDisabled(nodeID) {
			closeTerminal(s, "the terminal was disabled on this node")
			return
		}
	}
}

// terminalDisabled only reports a disabled terminal on positive information: a
// Redis hiccup must not tear down every established session.
func terminalDisabled(nodeID string) bool {
	redisConnection := redis.Instance()
	defer redisConnection.Close()

	enabled, err := redisConnection.HGet(terminalCtx, "node/"+nodeID+"/terminal", "enabled").Result()
	if err != nil {
		return false
	}
	return enabled != "1"
}

func closeTerminal(s *melody.Session, reason string) {
	state := terminalStateOf(s)
	if state == nil {
		_ = s.CloseWithMsg(melody.FormatCloseMessage(1000, "Bye"))
		return
	}

	state.mu.Lock()
	if state.closed {
		state.mu.Unlock()
		return
	}
	state.closed = true
	state.phase = phaseClosed
	if state.done != nil {
		close(state.done)
	}
	client := state.client
	session := state.session
	nodeID := state.nodeID
	holderID := state.holderID
	user := state.user
	sshUser := state.sshUser
	localPort := state.localPort
	opened := state.opened
	started := client != nil
	state.mu.Unlock()

	if session != nil {
		_ = session.Close()
	}
	if client != nil {
		_ = client.Close()
	}
	if nodeID != "" {
		releaseTerminal(nodeID, holderID)
	}

	if started {
		auditTerminal(user, "terminal-close", gin.H{
			"node":       nodeID,
			"ssh_user":   sshUser,
			"local_port": localPort,
			"reason":     reason,
			"seconds":    int(time.Since(opened).Seconds()),
		})
	}

	writeTerminalControl(s, gin.H{"type": "closed", "reason": reason})
	// melody writes the close frame straight to the socket while queued
	// messages are still in the session buffer, so closing immediately would
	// drop the reason the browser needs to display.
	time.Sleep(50 * time.Millisecond)
	_ = s.CloseWithMsg(melody.FormatCloseMessage(1000, "Bye"))
}

/*
 * Correlation record
 */

func writeCorrelationRecord(client *ssh.Client, adminUser string, clientIP string) error {
	// SSH exec carries no argv: sshd hands the string to the user's shell, so
	// the identity is interpolated there. add-user constrains cluster-admin
	// names today, but that guarantee lives in another component's schema, so
	// validate here too and quote.
	if !sshUsernamePattern.MatchString(adminUser) {
		return errors.New("cluster-admin name outside the expected character set")
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	message := fmt.Sprintf("session opened by %s from %s", adminUser, clientIP)

	// Bounded on purpose. This runs a command through the account's login
	// shell, and a shell that never returns would otherwise hang the caller
	// with no timeout of its own watching over it. The deferred Close unblocks
	// the goroutine, and the buffered channel keeps it from leaking.
	done := make(chan error, 1)
	go func() {
		done <- session.Run("logger -t ns8-terminal -- " + shellSingleQuote(message))
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(terminalCorrelationTimeout):
		return errors.New("timed out writing the correlation record")
	}
}

func shellSingleQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

/*
 * Throttle. Resistance to password guessing rests entirely here: pam_faillock
 * has even_deny_root disabled by default, so root is not locked out node-side,
 * and faillock is not in Debian's default stack.
 *
 * Counted per node and SSH account, and per cluster-admin user. Kept in this
 * map rather than on the WebSocket session so that reconnecting does not reset
 * the counter.
 */

type throttleEntry struct {
	failures  int
	notBefore time.Time
}

var throttleMu sync.Mutex
var throttleByAccount = map[string]*throttleEntry{}
var throttleByUser = map[string]*throttleEntry{}

func throttleDelay(failures int) time.Duration {
	if failures <= 1 {
		return 0
	}
	delay := time.Duration(1<<min(failures-1, 9)) * time.Second
	if delay > terminalThrottleCeiling {
		return terminalThrottleCeiling
	}
	return delay
}

func throttleRetryAfter(nodeID string, sshUser string, clusterUser string) time.Duration {
	throttleMu.Lock()
	defer throttleMu.Unlock()

	pruneThrottleLocked()

	longest := time.Duration(0)
	for _, entry := range []*throttleEntry{
		throttleByAccount[nodeID+"/"+sshUser],
		throttleByUser[clusterUser],
	} {
		if entry == nil {
			continue
		}
		if wait := time.Until(entry.notBefore); wait > longest {
			longest = wait
		}
	}
	return longest
}

func throttleRecordFailure(nodeID string, sshUser string, clusterUser string) {
	throttleMu.Lock()
	defer throttleMu.Unlock()

	for _, bucket := range []struct {
		table map[string]*throttleEntry
		key   string
	}{
		{throttleByAccount, nodeID + "/" + sshUser},
		{throttleByUser, clusterUser},
	} {
		entry := bucket.table[bucket.key]
		if entry == nil {
			entry = &throttleEntry{}
			bucket.table[bucket.key] = entry
		}
		entry.failures++
		entry.notBefore = time.Now().Add(throttleDelay(entry.failures))
	}
}

func throttleRecordSuccess(nodeID string, sshUser string, clusterUser string) {
	throttleMu.Lock()
	defer throttleMu.Unlock()
	delete(throttleByAccount, nodeID+"/"+sshUser)
	delete(throttleByUser, clusterUser)
}

func pruneThrottleLocked() {
	horizon := time.Now().Add(-terminalThrottleCeiling)
	for _, table := range []map[string]*throttleEntry{throttleByAccount, throttleByUser} {
		for key, entry := range table {
			if entry.notBefore.Before(horizon) {
				delete(table, key)
			}
		}
	}
}
