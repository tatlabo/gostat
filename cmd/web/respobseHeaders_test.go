package main

import (
	"bytes"
	"gostats/cmd/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetHeaderFunc(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		s := strings.Builder{}
		s.WriteString("OK")
		io.WriteString(w, s.String())
	})

	setHeaderFunc(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	expected := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com cdn.jsdelivr.net; font-src fonts.gstatic.com; script-src 'self' cdn.jsdelivr.net; img-src 'self' data:;",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",

		"Server": "Go",

		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "public, max-age=3600",

		"Transfer-Encoding": "chunked",
	}

	for key, value := range expected {
		assert.Equal(t, value, res.Header.Get(key))
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
