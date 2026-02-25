package main

import (
	"context"
	"flag"
	"gostats/cmd/internal/database"
	"gostats/cmd/internal/models"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	Snippets *models.SnippetModel
	Snippet  models.Snippet
}

func main() {

	app := &application{
		Snippets: &models.SnippetModel{
			DB: database.New(),
		},
		Snippet: models.Snippet{},
	}

	var addr = flag.String("addr", ":5000", "HTTP network address")

	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.Parse()

	logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, app.Routes())
	if err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

}

func (app *application) Routes() *http.ServeMux {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	media := http.FileServer(http.Dir("media"))

	mux.Handle("/media/", http.StripPrefix("/media/", media))

	mux.HandleFunc("GET /{$}", setHeaders(hello))

	mux.HandleFunc("GET /snippet", setHeaders(app.snippet))

	mux.HandleFunc("GET /snippet/all", setHeaders(app.snippetList))

	mux.HandleFunc("POST /snippet/create", setHeaders(app.snippetCreate))

	// mux.HandleFunc("GET /snippet/{id}", setHeaders(snippet))

	return mux
}

func setHeaders(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "title", "Home page")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Server", "GO-Server/1.0")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
