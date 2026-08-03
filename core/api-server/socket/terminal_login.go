/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package socket

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

var (
	errLoginTimeout = errors.New("no answer at the login prompt")
	errLoginAborted = errors.New("login cancelled")
)

/*
 * Login prompts
 *
 * SSH carries the user name inside the authentication request, so there is no
 * such thing as a remote "login:" prompt: only the password questions come from
 * the node, through keyboard-interactive. The user name prompt below is drawn
 * by api-server, which is why this file has to implement line editing.
 *
 * Nothing here narrows who can read the password. It still crosses this process
 * in clear on its way to the node, exactly as the credentials frame did. What
 * changes is that the browser no longer holds it in a form field.
 */

func runLoginFlow(s *melody.Session, state *terminalState) {
	deadline := time.Now().Add(terminalLoginTimeout)

	// Resolved once: it depends on the node, not on the account being tried, and
	// a node with no published host key must say so instead of looking like
	// three rejected passwords.
	target, hostKeys, err := nodeSSHTarget(state.nodeID)
	if err != nil {
		writeTerminalText(s, "\r\n"+sanitizePrompt(err.Error())+"\r\n")
		closeTerminal(s, "target unavailable")
		return
	}

	// The keepalive loop only watches the disabled flag once a shell is up, so
	// without this an operator revoking the terminal could not interrupt a
	// session parked at the password prompt.
	go watchDisabledDuringLogin(s, state)

	writeTerminalText(s, fmt.Sprintf("\r\nNS8 node %s\r\n", state.nodeID))

	for attempt := 0; attempt < terminalLoginAttempts; attempt++ {
		if state.currentPhase() != phaseLogin {
			return
		}

		writeTerminalText(s, "\r\nlogin: ")
		raw, err := readLine(s, state, deadline, true)
		if err != nil {
			closeTerminal(s, err.Error())
			return
		}

		sshUser := string(raw)
		if !sshUsernamePattern.MatchString(sshUser) {
			state.dropLoginRemainder()
			writeTerminalText(s, "\r\ninvalid user name\r\n")
			continue
		}

		// Checked per attempt and not once per session: the operator may switch
		// to an account that is already being guessed.
		if wait := throttleRetryAfter(state.nodeID, sshUser, state.user); wait > 0 {
			writeTerminalText(s, fmt.Sprintf(
				"\r\ntoo many failed attempts, try again in %d seconds\r\n", int(wait.Seconds())))
			closeTerminal(s, "throttled")
			return
		}

		// Cancelling at the prompt surfaces as a failed handshake too, so it is
		// recorded here rather than inferred from the dial error: a Ctrl-C is
		// not a wrong password and must not feed the throttle or the audit.
		var promptErr error
		client, dialErr := dialNode(target, hostKeys, sshUser,
			func(name string, instruction string, questions []string, echos []bool) ([]string, error) {
				answers, err := answerPrompts(s, state, deadline, name, instruction, questions, echos)
				if err != nil {
					promptErr = err
				}
				return answers, err
			})

		if promptErr != nil {
			closeTerminal(s, promptErr.Error())
			return
		}

		if dialErr != nil {
			throttleRecordFailure(state.nodeID, sshUser, state.user)
			// Audited with the cluster-admin identity: lastb on the node only
			// ever shows the leader address, so without this line an
			// administrator guessing root passwords leaves no attributable
			// trace.
			auditTerminal(state.user, "terminal-auth-failed", gin.H{
				"node":      state.nodeID,
				"ssh_user":  sshUser,
				"client_ip": state.clientIP,
			})
			state.dropLoginRemainder()
			// Deliberately vague and uniform, as sshd is: telling the operator
			// whether the account exists would be a gift to anyone guessing.
			writeTerminalText(s, "\r\nLogin incorrect\r\n")
			continue
		}

		throttleRecordSuccess(state.nodeID, sshUser, state.user)
		attachShell(s, state, sshUser, client)
		return
	}

	closeTerminal(s, "too many failed login attempts")
}

func watchDisabledDuringLogin(s *melody.Session, state *terminalState) {
	ticker := time.NewTicker(terminalKeepalive)
	defer ticker.Stop()

	for {
		select {
		case <-state.done:
			return
		case <-ticker.C:
			if state.currentPhase() != phaseLogin {
				return
			}
			if terminalDisabled(state.nodeID) {
				closeTerminal(s, "the terminal was disabled on this node")
				return
			}
		}
	}
}

// answerPrompts renders the questions sshd sends over keyboard-interactive and
// collects the replies. The node writes these strings, so they are stripped of
// control characters before reaching the terminal.
func answerPrompts(s *melody.Session, state *terminalState, deadline time.Time,
	name string, instruction string, questions []string, echos []bool) ([]string, error) {

	if text := sanitizePrompt(name); text != "" {
		writeTerminalText(s, "\r\n"+text+"\r\n")
	}
	if text := sanitizePrompt(instruction); text != "" {
		writeTerminalText(s, "\r\n"+text+"\r\n")
	}

	answers := make([]string, len(questions))
	for i, question := range questions {
		echo := i < len(echos) && echos[i]

		// "Password: " already ends with a space; a prompt that does not would
		// otherwise have the cursor stuck against its last character.
		prompt := sanitizePrompt(question)
		if !strings.HasSuffix(prompt, " ") {
			prompt += " "
		}

		writeTerminalText(s, "\r\n"+prompt)
		raw, err := readLine(s, state, deadline, echo)
		if err != nil {
			return nil, err
		}

		// Same limit as the credentials frame had: the string is immutable, so
		// this copy cannot be wiped. Our own buffer can, and is.
		answers[i] = string(raw)
		for j := range raw {
			raw[j] = 0
		}

		if !echo {
			writeTerminalText(s, "\r\n")
		}
	}

	return answers, nil
}

/*
 * Line editing
 *
 * The prompt runs before any shell exists, so no pty is echoing or handling
 * erase for us. Keep it to what a login prompt needs, and make sure escape
 * sequences from arrow keys never reach the buffer: a stray "\x1b[A" inside a
 * user name would be invisible on screen and rejected with a puzzling error.
 */

func readLine(s *melody.Session, state *terminalState, deadline time.Time, echo bool) ([]byte, error) {
	// Allocated once at full size so append never moves the buffer: a
	// reallocation would leave a copy of the password behind for the collector.
	line := make([]byte, 0, terminalLineLimit)

	const (
		escNone = iota
		escSeen
		escCSI
	)
	escape := escNone

	wipe := func(buffer []byte) {
		for i := range buffer {
			buffer[i] = 0
		}
	}

	for {
		// Whatever the previous prompt left behind comes first, before waiting
		// on the socket again.
		chunk := state.takeLoginRemainder()

		if len(chunk) == 0 {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				wipe(line)
				return nil, errLoginTimeout
			}

			timer := time.NewTimer(remaining)
			select {
			case chunk = <-state.loginInput:
				timer.Stop()
			case <-state.done:
				timer.Stop()
				wipe(line)
				return nil, errLoginAborted
			case <-timer.C:
				wipe(line)
				return nil, errLoginTimeout
			}
		}

		for index := 0; index < len(chunk); index++ {
			character := chunk[index]
			switch escape {
			case escSeen:
				// Only CSI and SS3 carry parameters; anything else was a lone
				// Escape followed by an ordinary key.
				if character == '[' || character == 'O' {
					escape = escCSI
				} else {
					escape = escNone
				}
				continue
			case escCSI:
				if character >= '@' && character <= '~' {
					escape = escNone
				}
				continue
			}

			switch character {
			case '\r', '\n':
				// A paste can hold the user name and the password in one
				// frame. Keep the tail for the next prompt instead of dropping
				// it, and swallow the LF of a CRLF pair so it does not read as
				// an empty answer there.
				rest := chunk[index+1:]
				if character == '\r' && len(rest) > 0 && rest[0] == '\n' {
					rest = rest[1:]
				}
				if len(rest) > 0 {
					keep := make([]byte, len(rest))
					copy(keep, rest)
					state.stashLoginRemainder(keep)
				}

				wipe(chunk)
				answer := make([]byte, len(line))
				copy(answer, line)
				wipe(line)
				return answer, nil

			case 0x7f, 0x08:
				if len(line) > 0 {
					line = line[:len(line)-1]
					if echo {
						writeTerminalText(s, "\b \b")
					}
				}

			case 0x03:
				wipe(chunk)
				wipe(line)
				return nil, errLoginAborted

			case 0x04:
				if len(line) == 0 {
					wipe(chunk)
					return nil, errLoginAborted
				}

			case 0x15:
				if echo {
					writeTerminalText(s, strings.Repeat("\b \b", len(line)))
				}
				line = line[:0]

			case 0x1b:
				escape = escSeen

			default:
				// Bytes above 0x7f are kept as they come: a UTF-8 password
				// arrives as several of them and must be reassembled untouched.
				if character < 0x20 || len(line) >= terminalLineLimit {
					continue
				}
				line = append(line, character)
				if echo {
					writeTerminalRaw(s, []byte{character})
				}
			}
		}

		wipe(chunk)
	}
}

// sanitizePrompt keeps a hostile or misconfigured node from driving the
// operator's terminal through the text of a PAM prompt.
func sanitizePrompt(text string) string {
	var clean strings.Builder
	for _, character := range text {
		if character == '\t' {
			clean.WriteRune(character)
			continue
		}
		if character < 0x20 || character == 0x7f {
			continue
		}
		clean.WriteRune(character)
	}
	return clean.String()
}

func writeTerminalRaw(s *melody.Session, data []byte) {
	_ = s.WriteBinary(data)
}

func writeTerminalText(s *melody.Session, text string) {
	_ = s.WriteBinary([]byte(text))
}
