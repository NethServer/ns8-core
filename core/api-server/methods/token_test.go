/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package methods

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignModuleToken(t *testing.T) {
	key := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{3}, ed25519.SeedSize))
	tok, err := SignModuleToken(key, "admin", "metrics1")
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(tok, claims, func(*jwt.Token) (any, error) {
		return key.Public(), nil
	}, jwt.WithValidMethods([]string{"EdDSA"}), jwt.WithIssuer("ns8-api-server"), jwt.WithAudience("module/metrics1"))
	if err != nil || parsed.Header["alg"] != "EdDSA" || claims["sub"] != "admin" {
		t.Fatalf("bad token: %v %v", err, claims)
	}
	iat, _ := claims["iat"].(float64)
	exp, _ := claims["exp"].(float64)
	if exp-iat != 60 {
		t.Errorf("iat=%v exp=%v", iat, exp)
	}
	raw, _ := base64.RawURLEncoding.DecodeString(strings.Split(tok, ".")[1])
	var m map[string]json.RawMessage
	json.Unmarshal(raw, &m)
	if string(m["aud"]) != `"module/metrics1"` {
		t.Errorf("aud raw = %s", m["aud"])
	}
}

func TestModuleTokenPubKeyVerifies(t *testing.T) {
	key := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{4}, ed25519.SeedSize))
	tok, _ := SignModuleToken(key, "admin", "metrics1")
	data, err := ModuleTokenPubKeyPEM(key)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		t.Fatalf("bad PEM: %s", data)
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := jwt.Parse(tok, func(*jwt.Token) (any, error) { return pub, nil }, jwt.WithValidMethods([]string{"EdDSA"})); err != nil {
		t.Fatalf("token does not verify: %v", err)
	}
}
