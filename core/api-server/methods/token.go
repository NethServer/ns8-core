/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package methods

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
)

const (
	moduleTokenIssuer = "ns8-api-server"
	moduleTokenTTL    = 60 * time.Second
)

// SignModuleToken returns a short-lived EdDSA JWT asserting that username
// was authenticated for the given module. MapClaims is used on purpose, so that "aud"
// is marshalled as a plain string, as expected by Grafana expect_claims.
func SignModuleToken(key ed25519.PrivateKey, username, moduleID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": moduleTokenIssuer,
		"sub": username,
		"aud": "module/" + moduleID,
		"iat": now.Unix(),
		"exp": now.Add(moduleTokenTTL).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(key)
}

// ModuleTokenPubKeyPEM returns the PEM public key matching the signing key.
func ModuleTokenPubKeyPEM(key ed25519.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

// ModuleTokenPubKey serves the public key to verify X-Auth-Token.
func ModuleTokenPubKey(c *gin.Context) {
	data, err := ModuleTokenPubKeyPEM(configuration.Config.TokenSigningKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "application/x-pem-file", data)
}
