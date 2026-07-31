/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package socket

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	"github.com/NethServer/ns8-core/core/api-server/audit"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

const (
	// Time allowed to attach the WebSocket with the ticket obtained from the
	// REST call, and then to send the credentials on that channel.
	terminalAttachTimeout      = 30 * time.Second
	terminalCredentialsTimeout = 30 * time.Second

	terminalIdleTimeout = 15 * time.Minute
	terminalMaxDuration = 8 * time.Hour
	terminalKeepalive   = 30 * time.Second

	// Bigger than the shared /ws instance, which keeps melody's 512 byte
	// default: a paste would otherwise exceed the read limit and close the
	// connection. Raising it there would raise it for that endpoint's clients.
	terminalMaxMessageSize  = 128 << 10
	terminalMessageBuffer   = 1024
	terminalCoalesceWindow  = 15 * time.Millisecond
	terminalThrottleCeiling = 5 * time.Minute
)

// Phases of a terminal WebSocket. A message that does not belong to the
// current phase closes the session instead of being ignored.
const (
	phaseTicket      = "ticket"
	phaseCredentials = "credentials"
	phaseRunning     = "running"
	phaseClosed      = "closed"
)

var terminalCtx = context.Background()

var terminalConn *melody.Melody

var terminalMu sync.Mutex

// Outstanding tickets, keyed by their hex token.
var terminalTickets = map[string]*terminalTicket{}

// Cluster-admin user currently holding a node, enforcing one session per node.
var terminalHolders = map[string]string{}

type terminalTicket struct {
	token     []byte
	nodeID    string
	user      string
	clientIP  string
	exp       int64
	expiresAt time.Time
}

// terminalControl is the typed shape of every inbound text frame. Inbound
// frames are never logged verbatim: a malformed one would otherwise put the
// password in the journal through utils.LogError, which writes err.Error() to
// stderr.
type terminalControl struct {
	Type     string `json:"type"`
	Ticket   string `json:"ticket"`
	Username string `json:"username"`
	Password string `json:"password"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
}

// TerminalInstance returns the melody instance dedicated to terminal sessions.
func TerminalInstance() *melody.Melody {
	if terminalConn != nil {
		return terminalConn
	}
	terminalConn = melody.New()
	terminalConn.Config.MaxMessageSize = terminalMaxMessageSize
	terminalConn.Config.MessageBufferSize = terminalMessageBuffer

	terminalConn.HandleConnect(onTerminalConnect)
	terminalConn.HandleMessage(onTerminalControl)
	terminalConn.HandleMessageBinary(onTerminalInput)
	terminalConn.HandleDisconnect(onTerminalDisconnect)

	// Fail closed. melody drops outbound frames when the session buffer is
	// full and Write reports nothing to the caller, so a lost frame towards
	// the node would silently truncate a command line before Enter.
	terminalConn.HandleError(func(s *melody.Session, err error) {
		closeTerminal(s, "transport error: "+err.Error())
	})

	return terminalConn
}

/*
 * CreateTerminalSession authorizes a terminal on a node and returns a one-shot
 * ticket. It carries no credentials and opens no SSH connection: the handshake
 * happens on the WebSocket channel, so a browser-held signer can replace the
 * password later without touching this route.
 *
 * The request body is {"action": "open-terminal"} so the existing Authorizator
 * matches it against the grants of the node taken from the URL.
 */
func CreateTerminalSession(c *gin.Context) {
	nodeID := c.Param("node_id")
	info := jwt.ExtractClaims(c)
	user, _ := info["id"].(string)
	exp, _ := info["exp"].(float64)

	redisConnection := redis.Instance()
	defer redisConnection.Close()

	enabled, _ := redisConnection.HGet(terminalCtx, "node/"+nodeID+"/terminal", "enabled").Result()
	if enabled != "1" {
		c.JSON(http.StatusForbidden, structs.Map(response.StatusForbidden{
			Code:    403,
			Message: "terminal is not enabled on this node",
			Data:    nil,
		}))
		return
	}

	ticket, err := reserveTerminal(nodeID, user, c.ClientIP(), int64(exp))
	if err != nil {
		c.JSON(http.StatusConflict, structs.Map(response.StatusConflict{
			Code:    409,
			Message: err.Error(),
			Data:    nil,
		}))
		return
	}

	auditTerminal(user, "terminal-open-requested", gin.H{
		"node":      nodeID,
		"client_ip": c.ClientIP(),
	})

	c.JSON(http.StatusCreated, structs.Map(response.StatusCreated{
		Code:    201,
		Message: "terminal session ticket created",
		Data:    gin.H{"ticket": ticket},
	}))
}

// reserveTerminal mints a ticket and marks the node busy. The reservation is
// taken here and must be released when the handshake fails, otherwise a single
// wrong password would lock the node until the ticket expires.
func reserveTerminal(nodeID string, user string, clientIP string, exp int64) (string, error) {
	terminalMu.Lock()
	defer terminalMu.Unlock()

	expireTicketsLocked()

	if holder, busy := terminalHolders[nodeID]; busy {
		return "", errTerminalBusy(holder)
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := hex.EncodeToString(raw)

	terminalTickets[token] = &terminalTicket{
		token:     raw,
		nodeID:    nodeID,
		user:      user,
		clientIP:  clientIP,
		exp:       exp,
		expiresAt: time.Now().Add(terminalAttachTimeout),
	}
	terminalHolders[nodeID] = user

	return token, nil
}

// consumeTicket validates and removes a ticket. The source address must match
// the one that obtained it.
func consumeTicket(token string, clientIP string) (*terminalTicket, error) {
	terminalMu.Lock()
	defer terminalMu.Unlock()

	expireTicketsLocked()

	ticket, found := terminalTickets[token]
	if !found {
		return nil, errTerminalTicket("unknown or already used ticket")
	}

	raw, err := hex.DecodeString(token)
	if err != nil || subtle.ConstantTimeCompare(raw, ticket.token) != 1 {
		return nil, errTerminalTicket("unknown or already used ticket")
	}

	delete(terminalTickets, token)

	if ticket.clientIP != clientIP {
		releaseTerminalLocked(ticket.nodeID)
		return nil, errTerminalAddressChanged
	}

	return ticket, nil
}

func expireTicketsLocked() {
	now := time.Now()
	for token, ticket := range terminalTickets {
		if now.After(ticket.expiresAt) {
			delete(terminalTickets, token)
			releaseTerminalLocked(ticket.nodeID)
		}
	}
}

func releaseTerminal(nodeID string) {
	terminalMu.Lock()
	defer terminalMu.Unlock()
	releaseTerminalLocked(nodeID)
}

func releaseTerminalLocked(nodeID string) {
	delete(terminalHolders, nodeID)
}

func auditTerminal(user string, action string, data gin.H) {
	payload, _ := json.Marshal(data)
	audit.Store(models.Audit{
		ID:        0,
		User:      user,
		Action:    action,
		Data:      string(payload),
		Timestamp: time.Now().UTC(),
	})
}

/*
 * Phase handling
 */

func onTerminalConnect(s *melody.Session) {
	state := &terminalState{phase: phaseTicket, opened: time.Now()}
	s.Set("terminal", state)

	go func() {
		time.Sleep(terminalAttachTimeout)
		if state.currentPhase() == phaseTicket {
			closeTerminal(s, "ticket not presented in time")
		}
	}()
}

func onTerminalControl(s *melody.Session, message []byte) {
	var control terminalControl
	err := json.Unmarshal(message, &control)
	// Wipe the frame before doing anything else: it may hold the password.
	for i := range message {
		message[i] = 0
	}
	if err != nil {
		// Deliberately a fixed string: the frame must never be logged.
		utils.LogError(errTerminalControl)
		closeTerminal(s, "malformed control frame")
		return
	}

	state := terminalStateOf(s)
	if state == nil {
		closeTerminal(s, "no session state")
		return
	}

	switch control.Type {
	case "ticket":
		handleTicketFrame(s, state, control)
	case "credentials":
		handleCredentialsFrame(s, state, &control)
	case "resize":
		state.resize(control.Rows, control.Cols)
	default:
		closeTerminal(s, "unexpected control frame")
	}
}

func handleTicketFrame(s *melody.Session, state *terminalState, control terminalControl) {
	if state.currentPhase() != phaseTicket {
		closeTerminal(s, "ticket already presented")
		return
	}

	ticket, err := consumeTicket(control.Ticket, clientAddress(s))
	if err != nil {
		writeTerminalControl(s, gin.H{"type": "auth-error", "message": err.Error()})
		closeTerminal(s, "ticket rejected")
		return
	}

	state.adopt(ticket)
	state.setPhase(phaseCredentials)

	go func() {
		time.Sleep(terminalCredentialsTimeout)
		if state.currentPhase() == phaseCredentials {
			closeTerminal(s, "credentials not sent in time")
		}
	}()

	writeTerminalControl(s, gin.H{"type": "ticket-accepted"})
}

func handleCredentialsFrame(s *melody.Session, state *terminalState, control *terminalControl) {
	if state.currentPhase() != phaseCredentials {
		closeTerminal(s, "credentials not expected now")
		return
	}

	// The JSON decoder materialises the password as a Go string, which is
	// immutable and may be copied by the garbage collector, so it cannot be
	// wiped. What we can do is keep a single copy, hand it straight to the
	// handshake, blank our own buffer afterwards, never log it, and set
	// LimitCORE=0 on the unit so it cannot reach a core dump.
	password := []byte(control.Password)
	control.Password = ""

	startTerminalSession(s, state, control.Username, password, control.Rows, control.Cols)

	for i := range password {
		password[i] = 0
	}
}

func onTerminalInput(s *melody.Session, message []byte) {
	state := terminalStateOf(s)
	if state == nil || state.currentPhase() != phaseRunning {
		return
	}
	state.write(message)
}

func onTerminalDisconnect(s *melody.Session) {
	closeTerminal(s, "client disconnected")
}

func terminalStateOf(s *melody.Session) *terminalState {
	value, found := s.Get("terminal")
	if !found {
		return nil
	}
	state, _ := value.(*terminalState)
	return state
}

// clientAddress returns the address gin resolved for the handshake. It is
// stored by the route through HandleRequestWithKeys so that the trusted proxy
// logic applies here too: SetTrustedProxies is limited to the loopback, so a
// forged X-Forwarded-For does not win.
func clientAddress(s *melody.Session) string {
	value, found := s.Get("client_ip")
	if !found {
		return ""
	}
	address, _ := value.(string)
	return address
}

func writeTerminalControl(s *melody.Session, payload gin.H) {
	frame, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_ = s.Write(frame)
}
