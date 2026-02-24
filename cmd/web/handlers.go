package main

import (
	"net/http"
	"strconv"
	"text/template"
)

var templates = template.Must(template.ParseGlob("./cmd/ui/html/*.html"))
var tpl = templates.ExecuteTemplate

func hello(w http.ResponseWriter, r *http.Request) {

	var msg struct {
		Message string
		Title   string
	}

	msg.Title = "Snippet Page"

	msg.Message = "Hello, World! Everyone loves Go!"

	tpl(w, "home.html", msg)

}

func snippet(w http.ResponseWriter, r *http.Request) {

	var id int

	idUrl := r.URL.Query().Get("id")
	id, _ = strconv.Atoi(idUrl)

	var msg struct {
		Message string
		Title   string
	}

	msg.Title = "Snippet Page"

	if id == 0 {
		msg.Message = "Snnippet page"
		tpl(w, "home.html", msg)
		return
	}

	msg.Message = "Snnippet page"
	tpl(w, "home.html", msg)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	var snippet = app.snippet

	_, err := app.snippets.Insert(&snippet)
	if err != nil {
		http.Error(w, "Unable to create snippet", http.StatusInternalServerError)
		return
	}

}
