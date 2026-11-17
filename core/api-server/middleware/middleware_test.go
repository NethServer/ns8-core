/*
 * Copyright (C) 2026 Nethesis S.r.l.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimiter(1, 2))
	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	request := func(ip string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip + ":12345"
		router.ServeHTTP(w, req)
		return w
	}

	// the burst is allowed, then requests are throttled
	for i := 0; i < 2; i++ {
		if w := request("192.0.2.1"); w.Code != http.StatusOK {
			t.Fatalf("request %d: status %d", i, w.Code)
		}
	}
	w := request("192.0.2.1")
	if w.Code != http.StatusTooManyRequests || w.Header().Get("Retry-After") != "1" {
		t.Fatalf("throttled request: status %d, Retry-After %q", w.Code, w.Header().Get("Retry-After"))
	}

	// another client has its own bucket
	if w := request("192.0.2.2"); w.Code != http.StatusOK {
		t.Fatalf("other client: status %d", w.Code)
	}
}

func TestRateLimiterZeroBurst(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimiter(1, 0))
	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("a zero burst must be clamped to 1: status %d", w.Code)
	}
}

func TestBodyLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/", BodyLimit(16), func(c *gin.Context) {
		if _, err := io.ReadAll(c.Request.Body); err != nil {
			c.String(http.StatusRequestEntityTooLarge, "too large")
			return
		}
		c.String(http.StatusOK, "ok")
	})

	for _, tc := range []struct {
		body string
		want int
	}{
		{strings.Repeat("x", 16), http.StatusOK},
		{strings.Repeat("x", 17), http.StatusRequestEntityTooLarge},
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body)))
		if w.Code != tc.want {
			t.Errorf("body of %d bytes: status %d, want %d", len(tc.body), w.Code, tc.want)
		}
	}
}
