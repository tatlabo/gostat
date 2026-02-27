package main

import (
	"flag"
	"gostats/cmd/internal/database"
	"gostats/cmd/internal/models"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Application struct {
	Snippets *models.SnippetModel
	Snippet  models.Snippet
	Template func(wr io.Writer, name string, data any) error
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

func (app *Application) Routes() *http.ServeMux {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))

	media := http.FileServer(http.Dir("media"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/media/", http.StripPrefix("/media/", media))

	mux.HandleFunc("GET /{$}", setHeaders(hello))

	mux.HandleFunc("GET /snippet", (app.snippet))

	mux.HandleFunc("GET /snippet/all", (app.snippetList))

	mux.HandleFunc("POST /snippet/create", (app.snippetCreate))

	mux.HandleFunc("/", app.notFound)

	return mux
}

func setHeaders(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Server", "GO-Server/1.0")

		w.Header().Set("Transfer-Encoding", "chunked")
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
