package main

import (
	"bytes"
	"gostats/cmd/internal/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	app := &Application{
		Logger: slog.New(slog.DiscardHandler),
	}

	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
