/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package middleware

import (
	"errors"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
)

const testSecret = "test-secret"

func signToken(t *testing.T, method gojwt.SigningMethod, key any, claims gojwt.MapClaims) string {
	t.Helper()
	token, err := gojwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

// TestTokenErrorMessage checks that token validation errors keep the
// texts of gin-jwt v2: clients match them, e.g. the UI shows the
// "session expired" notice for "Token is expired".
func TestTokenErrorMessage(t *testing.T) {
	configuration.Config.Secret = testSecret
	mw := InitJWT()
	now := time.Now()
	valid := func() gojwt.MapClaims {
		return gojwt.MapClaims{"id": "admin", "exp": now.Add(time.Hour).Unix(), "orig_iat": now.Unix()}
	}

	expired := valid()
	expired["exp"] = now.Add(-time.Minute).Unix()
	notYet := valid()
	notYet["nbf"] = now.Add(time.Hour).Unix()

	cases := []struct {
		name  string
		token string
		want  string
	}{
		{"expired", signToken(t, gojwt.SigningMethodHS256, []byte(testSecret), expired), "Token is expired"},
		{"not valid yet", signToken(t, gojwt.SigningMethodHS256, []byte(testSecret), notYet), "Token is not valid yet"},
		{"wrong signature", signToken(t, gojwt.SigningMethodHS256, []byte("another-secret"), valid()), "signature is invalid"},
		{"other algorithm", signToken(t, gojwt.SigningMethodHS512, []byte(testSecret), valid()), "invalid signing algorithm"},
		{"malformed", "not-a-token", "token contains an invalid number of segments"},
	}
	for _, tc := range cases {
		_, err := mw.ParseTokenString(tc.token)
		if err == nil {
			t.Errorf("%s: token accepted", tc.name)
			continue
		}
		if got := TokenErrorMessage(err); got != tc.want {
			t.Errorf("%s: TokenErrorMessage(%q) = %q, want %q", tc.name, err, got, tc.want)
		}
	}

	if _, err := mw.ParseTokenString(signToken(t, gojwt.SigningMethodHS256, []byte(testSecret), valid())); err != nil {
		t.Errorf("valid token refused: %v", err)
	}
}

func TestTokenErrorMessageOther(t *testing.T) {
	if got := TokenErrorMessage(errors.New("some error")); got != "some error" {
		t.Errorf("unknown errors keep their text, got %q", got)
	}
}
