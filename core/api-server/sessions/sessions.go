/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of NethServer project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with NethServer.  If not, see COPYING.
 */

// Package sessions keeps the login sessions of cluster administrators in
// a SQLite database local to the node running the API server. Session
// data are never written to Redis, which is replicated to worker nodes.
//
// A session starts at login and is identified by a random ID, carried
// by access tokens in the "sid" claim. A refresh token, stored only as a
// SHA-256 hash, obtains a new access token and is replaced at each use.
package sessions

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// AccessTimeout is the lifetime of an access token
	AccessTimeout = 15 * time.Minute
	// IdleTimeout ends a session that is not refreshed for this time
	IdleTimeout = 12 * time.Hour
	// MaxLifetime ends a session this time after login
	MaxLifetime = 14 * 24 * time.Hour
	// rotationGrace keeps a replaced refresh token usable for a short
	// time, so concurrent refresh requests (e.g. from two browser tabs)
	// do not end the session
	rotationGrace = 30 * time.Second
	// cleanupInterval is the period of expired sessions removal
	cleanupInterval = 10 * time.Minute
)

// ErrInvalid is returned when a session or refresh token is unknown,
// expired or revoked
var ErrInvalid = errors.New("session is not valid")

// ErrUnavailable is returned when the session database cannot be used
var ErrUnavailable = errors.New("session store is not available")

// Session describes a login session
type Session struct {
	ID        string
	Username  string
	OtpPass   bool
	CreatedAt time.Time
}

var (
	db  *sql.DB
	mu  sync.Mutex
	now = time.Now
)

// Init opens (or creates) the session database at path
func Init(path string) error {
	if path == "" {
		return errors.New("session database path is not set")
	}

	// create the file readable only by the API server
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	f.Close()

	conn, err := sql.Open("sqlite3", "file:"+path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return err
	}
	// a single connection serializes access: the load is a few queries
	// for each API request
	conn.SetMaxOpenConns(1)

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT NOT NULL PRIMARY KEY,
			username TEXT NOT NULL,
			otp_pass INTEGER NOT NULL,
			refresh_hash TEXT NOT NULL UNIQUE,
			prev_hash TEXT NOT NULL DEFAULT '',
			prev_until INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			idle_until INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS sessions_username ON sessions (username);
		CREATE INDEX IF NOT EXISTS sessions_prev_hash ON sessions (prev_hash);
	`)
	if err != nil {
		conn.Close()
		return err
	}

	db = conn
	go cleanupLoop()
	return nil
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Create starts a new session and returns its ID and refresh token
func Create(username string, otpPass bool) (string, string, error) {
	if db == nil {
		return "", "", ErrUnavailable
	}
	id, err := randomString(16)
	if err != nil {
		return "", "", err
	}
	refresh, err := randomString(32)
	if err != nil {
		return "", "", err
	}
	t := now()

	mu.Lock()
	defer mu.Unlock()
	_, err = db.Exec(`INSERT INTO sessions (id, username, otp_pass, refresh_hash, created_at, idle_until)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id, username, otpPass, hashToken(refresh), t.UnixNano(), t.Add(IdleTimeout).UnixNano())
	if err != nil {
		return "", "", err
	}
	return id, refresh, nil
}

// Refresh validates a refresh token, replaces it with a new one and
// extends the session idle timeout. A replaced refresh token presented
// after the grace period ends the session: it may have been stolen.
func Refresh(refresh string) (Session, string, error) {
	var s Session
	if db == nil {
		return s, "", ErrUnavailable
	}
	if refresh == "" {
		return s, "", ErrInvalid
	}
	newRefresh, err := randomString(32)
	if err != nil {
		return s, "", err
	}
	hash := hashToken(refresh)
	t := now()

	mu.Lock()
	defer mu.Unlock()

	var refreshHash, prevHash string
	var prevUntil, createdAt, idleUntil int64
	row := db.QueryRow(`SELECT id, username, otp_pass, refresh_hash, prev_hash, prev_until, created_at, idle_until
		FROM sessions WHERE refresh_hash = ? OR prev_hash = ?`, hash, hash)
	err = row.Scan(&s.ID, &s.Username, &s.OtpPass, &refreshHash, &prevHash, &prevUntil, &createdAt, &idleUntil)
	if err == sql.ErrNoRows {
		return s, "", ErrInvalid
	} else if err != nil {
		return s, "", err
	}
	s.CreatedAt = time.Unix(0, createdAt)

	if hash != refreshHash && t.UnixNano() > prevUntil {
		// replaced refresh token used again: revoke the session
		db.Exec(`DELETE FROM sessions WHERE id = ?`, s.ID)
		return s, "", ErrInvalid
	}
	if t.UnixNano() > idleUntil || t.After(s.CreatedAt.Add(MaxLifetime)) {
		db.Exec(`DELETE FROM sessions WHERE id = ?`, s.ID)
		return s, "", ErrInvalid
	}

	// The token being replaced stays usable for the grace period. If a
	// replaced token is used within the grace period, the current one
	// becomes the replaced token instead.
	_, err = db.Exec(`UPDATE sessions SET refresh_hash = ?, prev_hash = ?, prev_until = ?, idle_until = ? WHERE id = ?`,
		hashToken(newRefresh), refreshHash, t.Add(rotationGrace).UnixNano(), t.Add(IdleTimeout).UnixNano(), s.ID)
	if err != nil {
		return s, "", err
	}
	return s, newRefresh, nil
}

// Valid tells if the session exists and has not expired
func Valid(id string) bool {
	if db == nil || id == "" {
		return false
	}
	t := now()

	mu.Lock()
	defer mu.Unlock()
	var createdAt, idleUntil int64
	err := db.QueryRow(`SELECT created_at, idle_until FROM sessions WHERE id = ?`, id).Scan(&createdAt, &idleUntil)
	if err != nil {
		return false
	}
	return t.UnixNano() <= idleUntil && !t.After(time.Unix(0, createdAt).Add(MaxLifetime))
}

// Revoke ends the session with the given ID
func Revoke(id string) error {
	if db == nil {
		return ErrUnavailable
	}
	mu.Lock()
	defer mu.Unlock()
	_, err := db.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

// RevokeRefresh ends the session of the given refresh token
func RevokeRefresh(refresh string) error {
	if db == nil {
		return ErrUnavailable
	}
	hash := hashToken(refresh)
	mu.Lock()
	defer mu.Unlock()
	_, err := db.Exec(`DELETE FROM sessions WHERE refresh_hash = ? OR prev_hash = ?`, hash, hash)
	return err
}

// RevokeUser ends the sessions of username started before notBefore
func RevokeUser(username string, notBefore time.Time) (int64, error) {
	if db == nil {
		return 0, ErrUnavailable
	}
	mu.Lock()
	defer mu.Unlock()
	res, err := db.Exec(`DELETE FROM sessions WHERE username = ? AND created_at < ?`, username, notBefore.UnixNano())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// RevokeOthers ends the sessions of username, except the one with the given ID
func RevokeOthers(username string, id string) (int64, error) {
	if db == nil {
		return 0, ErrUnavailable
	}
	mu.Lock()
	defer mu.Unlock()
	res, err := db.Exec(`DELETE FROM sessions WHERE username = ? AND id != ?`, username, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func cleanup() {
	t := now()
	mu.Lock()
	defer mu.Unlock()
	if db == nil {
		return
	}
	db.Exec(`DELETE FROM sessions WHERE idle_until < ? OR created_at < ?`,
		t.UnixNano(), t.Add(-MaxLifetime).UnixNano())
}

func cleanupLoop() {
	for {
		cleanup()
		time.Sleep(cleanupInterval)
	}
}
