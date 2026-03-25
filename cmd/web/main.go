package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"gostats/cmd/internal/database"
	"gostats/cmd/internal/models"
	"gostats/static"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"

	_ "github.com/lib/pq"
)

type Application struct {
	DB             *sql.DB
	Snippets       *models.SnippetModel
	Snippet        *models.Snippet
	RenderHTML     func(wr io.Writer, name string, data any) error
	SessionManager *scs.SessionManager
	User           *models.UserModel
	Logger         *slog.Logger
	Validator      models.Validator
}

type flashMsg struct {
	Flash string
	Time  time.Time
}

func (app *Application) newTemplateData(r *http.Request) flashMsg {
	return flashMsg{
		Flash: app.SessionManager.PopString(r.Context(), "flash"),
		Time:  time.Now(),
	}
}

func main() {

	customTemplate, err := customTemplate()
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
	customTemplateExecute := customTemplate.ExecuteTemplate

	DB := database.New()
	// app instance with dependencies
	app := &Application{
		DB: DB,
		Snippets: &models.SnippetModel{
			DB: DB,
		},
		Snippet:        &models.Snippet{},
		RenderHTML:     customTemplateExecute,
		SessionManager: scs.New(),
		User: &models.UserModel{
			DB: DB,
		},
		Logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Validator: models.Validator{},
	}

	app.SessionManager.Store = postgresstore.New(DB)
	app.SessionManager.Lifetime = 1 * time.Hour
	//
	// default addr
	var addr = flag.String("addr", ":5500", "HTTP network address")
	//

	app.Logger.Info("HTTPS", "port", *addr)

	flag.Parse()

	//
	routes := app.Routes()
	//
	//listen and serve
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	//
	server := &http.Server{
		Addr:        *addr,
		Handler:     routes,
		IdleTimeout: time.Minute,
		TLSConfig:   tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	go app.HttpServer80(":5000")

	if err := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
		app.Logger.Error("server error", "error", err)
		os.Exit(1)
	}

}

func (app *Application) HttpServer80(port string) {

	server := &http.Server{
		Addr: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://localhost:5500", http.StatusSeeOther)
		}),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	app.Logger.Info("HTTP", "port", port)

	err := server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		app.Logger.Error("server error", "error", err)
		os.Exit(1)
	}

	<-done
	app.Logger.Info("Graceful shutdown complete.")
}

func (app *Application) Routes() http.Handler {

	mux := http.NewServeMux()

	sessions := app.SessionManager.LoadAndSave

	media := http.FileServer(http.Dir("media"))
	mux.Handle("GET /media/", http.StripPrefix("/media/", media))

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static.StaticFiles)))

	mux.Handle("GET /{$}", setHeaderFunc(hello))
	mux.Handle("GET /ping", setHeaderFunc(ping))

	// Snippet routes

	mux.Handle("GET /snippet", sessions(setHeaderFunc(app.snippet)))

	mux.Handle("GET /snippet/all", sessions(setHeaderFunc(app.snippetList)))

	mux.Handle("POST /snippet/create", sessions(setHeaderFunc(app.snippetCreate)))

	mux.Handle("POST /snippet/delete", sessions(setHeaderFunc(app.snippetDelete)))

	// User

	mux.Handle("GET /user/signup", sessions(setHeaderFunc(app.userSignup)))

	mux.Handle("POST /user/signup", sessions(setHeaderFunc(app.userSignupPost)))

	mux.Handle("GET /user/login", sessions(setHeaderFunc(app.userLogin)))

	mux.Handle("POST /user/login", sessions(setHeaderFunc(app.userLoginPost)))

	mux.Handle("POST /user/logout", sessions(setHeaderFunc(app.userLogoutPost)))

	// Default route for 404

	mux.Handle("/", setHeaderFunc(app.notFound))

	return app.recoverPanic(app.logRequest(mux))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func setHeaderFunc(next http.HandlerFunc) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		for key, value := range responseHeadersMap() {
			w.Header().Set(key, value)
		}

		next.ServeHTTP(w, r)
	}

	return fn
}

func setHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for key, value := range responseHeadersMap() {
			w.Header().Set(key, value)
		}
		next.ServeHTTP(w, r)
	})

}

func responseHeadersMap() map[string]string {

	return map[string]string{
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
}

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Execute our middleware logic here...
		next.ServeHTTP(w, r)
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			path   = r.URL.RequestURI()
		)

		app.Logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "path", path)

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

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown with error", "error", err)
	}

	logger.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
