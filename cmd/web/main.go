package main

import (
	"context"
	"flag"
	"fmt"
	"gostats/cmd/internal/database"
	"gostats/cmd/internal/models"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

// Key should be 32 or 64 bytes for AES-256
var store = sessions.NewCookieStore([]byte("your-secret-key-change-in-production"))

type Application struct {
	Snippets *models.SnippetModel
	Snippet  models.Snippet
	Template func(wr io.Writer, name string, data any) error
	Session  *sessions.CookieStore
}

func main() {

	customTemplate, err := customTemplate()
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
	customTemplateExecute := customTemplate.ExecuteTemplate

	// app instance with dependencies
	app := &Application{
		Snippets: &models.SnippetModel{
			DB: database.New(),
		},
		Snippet:  models.Snippet{},
		Template: customTemplateExecute,
		Session:  store,
	}
	//

	var addr = flag.String("addr", ":5000", "HTTP network address")

	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.Parse()

	logger.Info("starting server", "addr", *addr)
	//
	routes := app.Routes()
	//listen and serve
	//
	if err := http.ListenAndServe(*addr, routes); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

}

func (app *Application) Routes() http.Handler {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))

	media := http.FileServer(http.Dir("media"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/media/", http.StripPrefix("/media/", media))

	mux.Handle("GET /{$}", setHeaderFunc(hello))

	mux.Handle("GET /snippet", setHeaderFunc(app.snippet))

	mux.Handle("GET /snippet/all", setHeaderFunc(app.snippetList))

	mux.Handle("POST /snippet/create", setHeaderFunc(app.snippetCreate))

	mux.Handle("POST /snippet/delete", setHeaderFunc(app.snippetDelete))

	mux.Handle("/", setHeaderFunc(app.notFound))

	return app.recoverPanic(app.logRequest(mux))
}

func setHeaderFunc(next http.HandlerFunc) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		for key, value := range ResponseHeaders {
			w.Header().Set(key, value)
		}

		next.ServeHTTP(w, r)
	}

	return fn
}

func setHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for key, value := range ResponseHeaders {
			w.Header().Set(key, value)
		}
		next.ServeHTTP(w, r)
	})

}

var ResponseHeaders = map[string]string{
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

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Execute our middleware logic here...
		next.ServeHTTP(w, r)
	})
}

var funcMap = func() template.FuncMap {

	return template.FuncMap{
		"mod": func(i, j int) int {
			return i % j
		},
		"sub": func(i, j int) int {
			return i - j
		},
		"CurrentYear": func() int {
			return time.Now().Year()
		},
		"CurrentDay": func() string {
			return time.Now().Format("2006-01-02")
		},
	}
}

func customTemplate() (*template.Template, error) {

	parse, err := template.New("").Funcs(funcMap()).ParseGlob("./cmd/ui/html/*.html")
	if err != nil {
		return nil, err
	}

	return parse, nil
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			path   = r.URL.RequestURI()
		)

		slog.Info("Received request", "ip", ip, "proto", proto, "method", method, "path", path)

		next.ServeHTTP(w, r)

	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%v", err)
				log.Printf("%v", err)
				ctx := r.Context()
				ctx = context.WithValue(ctx, "error", msg)
				r = r.WithContext(ctx)
				app.error500(w, r)
			}
		}()

		next.ServeHTTP(w, r)

	})
}
