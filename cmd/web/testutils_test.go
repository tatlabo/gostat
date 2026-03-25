package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *Application {
	app := &Application{
		Logger: slog.New(slog.DiscardHandler),
	}

	return app
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

type testResponse struct {
	Status  int
	Header  http.Header
	Body    string
	cookies []*http.Cookie
}

func (ts *testServer) get(t *testing.T, urlPath string) testResponse {

	req, err := http.NewRequest("GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		Status:  res.StatusCode,
		Header:  res.Header,
		Body:    string(bytes.TrimSpace(body)),
		cookies: res.Cookies(),
	}

}
