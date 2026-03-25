package main

import (
	"gostats/cmd/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication()

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	res := ts.get(t, "/ping")

	assert.Equal(t, res.Status, http.StatusOK)
	assert.Equal(t, res.Body, "OK")

}
