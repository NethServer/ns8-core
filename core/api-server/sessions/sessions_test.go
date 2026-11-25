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

package sessions

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setup opens a fresh database and returns a function moving the clock
func setup(t *testing.T) func(time.Duration) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sessions.db")
	clock := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time { return clock }
	t.Cleanup(func() {
		db.Close()
		db = nil
		now = time.Now
	})
	if err := Init(path); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(path); err != nil || fi.Mode().Perm() != 0600 {
		t.Fatalf("database file mode: %v %v", fi.Mode().Perm(), err)
	}
	return func(d time.Duration) { clock = clock.Add(d) }
}

func TestRefreshRotation(t *testing.T) {
	advance := setup(t)
	id, r0, err := Create("admin", true)
	if err != nil || !Valid(id) {
		t.Fatalf("create: %v", err)
	}

	s, r1, err := Refresh(r0)
	if err != nil || s.ID != id || s.Username != "admin" || !s.OtpPass || r1 == r0 {
		t.Fatalf("refresh: %+v %v", s, err)
	}

	// the replaced token works within the grace period: r1 becomes
	// the replaced token
	advance(10 * time.Second)
	if _, _, err := Refresh(r0); err != nil {
		t.Fatalf("refresh within grace: %v", err)
	}

	// an older token is unknown
	advance(time.Minute)
	if _, _, err := Refresh(r0); err != ErrInvalid || !Valid(id) {
		t.Fatalf("unknown token: %v", err)
	}

	// the replaced token used after the grace period revokes the session
	if _, _, err := Refresh(r1); err != ErrInvalid {
		t.Fatalf("reuse after grace: %v", err)
	}
	if Valid(id) {
		t.Fatal("session still valid after refresh token reuse")
	}
}

func TestIdleAndMaxLifetime(t *testing.T) {
	advance := setup(t)
	id, r, _ := Create("admin", false)

	// keep refreshing: the session lives until MaxLifetime
	var err error
	for elapsed := time.Duration(0); elapsed < MaxLifetime-IdleTimeout; elapsed += 11 * time.Hour {
		advance(11 * time.Hour)
		if _, r, err = Refresh(r); err != nil {
			t.Fatalf("refresh at %v: %v", elapsed, err)
		}
	}
	advance(IdleTimeout - time.Minute)
	if Valid(id) {
		t.Fatal("session valid after MaxLifetime")
	}
	if _, _, err := Refresh(r); err != ErrInvalid {
		t.Fatalf("refresh after MaxLifetime: %v", err)
	}

	// idle timeout
	id, r, _ = Create("admin", false)
	advance(IdleTimeout + time.Second)
	if Valid(id) {
		t.Fatal("session valid after IdleTimeout")
	}
	if _, _, err := Refresh(r); err != ErrInvalid {
		t.Fatalf("refresh after IdleTimeout: %v", err)
	}
}

func TestRevoke(t *testing.T) {
	advance := setup(t)
	a1, r1, _ := Create("alice", false)
	a2, _, _ := Create("alice", false)
	b1, _, _ := Create("bob", false)

	if err := RevokeRefresh(r1); err != nil || Valid(a1) {
		t.Fatalf("revoke by refresh token: %v", err)
	}

	advance(time.Second)
	a3, _, _ := Create("alice", false)
	if n, err := RevokeUser("alice", now().Add(-time.Millisecond)); err != nil || n != 1 || Valid(a2) || !Valid(a3) {
		t.Fatalf("revoke user before: n=%d %v", n, err)
	}

	a4, _, _ := Create("alice", false)
	if n, err := RevokeOthers("alice", a4); err != nil || n != 1 || Valid(a3) || !Valid(a4) {
		t.Fatalf("revoke others: n=%d %v", n, err)
	}

	if err := Revoke(a4); err != nil || Valid(a4) || !Valid(b1) {
		t.Fatalf("revoke: %v", err)
	}
}
