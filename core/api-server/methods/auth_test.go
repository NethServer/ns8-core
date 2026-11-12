/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package methods

import "testing"

func TestCheckIPAllowed(t *testing.T) {
	cases := []struct {
		ip      string
		allowed string
		want    bool
	}{
		{"192.0.2.10", "", true},
		{"192.0.2.10", "192.0.2.0/24", true},
		{"198.51.100.1", "192.0.2.0/24", false},
		{"198.51.100.1", "192.0.2.0/24, 198.51.100.0/24", true},
		{"192.0.2.10", "192.0.2.10", true},
		{"192.0.2.11", "192.0.2.10", false},
		{"127.0.0.1", "127.0.0.1,::1", true},
		{"::1", "127.0.0.1,::1", true},
		{"2001:db8::5", "2001:db8::/32", true},
		{"2001:db9::5", "2001:db8::/32", false},
		{"192.0.2.10", "not-an-ip,192.0.2.0/24", true},
		{"192.0.2.10", "not-an-ip", false},
		{"not-an-ip", "192.0.2.0/24", false},
		// a list without valid entries denies access (fail closed)
		{"192.0.2.10", ",,", false},
	}
	for _, tc := range cases {
		if got := CheckIPAllowed(tc.ip, tc.allowed); got != tc.want {
			t.Errorf("CheckIPAllowed(%q, %q) = %v, want %v", tc.ip, tc.allowed, got, tc.want)
		}
	}
}
