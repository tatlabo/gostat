package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var templates = template.Must(template.ParseGlob("./cmd/ui/html/*.html"))
var tpl = templates.ExecuteTemplate

func hello(w http.ResponseWriter, r *http.Request) {

	var msg msg

	msg.Title = "Snippet Page"

	msg.Message = "Hello, World! Everyone loves Go!"

	tpl(w, "home.html", msg)

}

type msg struct {
	Message string
	Title   string
}

func snippet(w http.ResponseWriter, r *http.Request) {

	var id int

	idUrl := r.URL.Query().Get("id")
	id, _ = strconv.Atoi(idUrl)

	var msg msg
	msg.Title = "Snippet Page"

	if id == 0 {
		msg.Message = "Snnippet page"
		tpl(w, "home.html", msg)
		return
	}

	msg.Message = "Snnippet page nr " + strconv.Itoa(id)
	tpl(w, "home.html", msg)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	var s = app.Snippet
	s.Title = r.FormValue("title")
	s.Content = r.FormValue("content")
	s.Expires = r.FormValue("expires")

	res, err := app.Snippets.Insert(s)
	if err != nil {
		fmt.Printf("ERROR inserting snippet: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to create snippet: %v", err), http.StatusInternalServerError)
		return
	}

	var msg msg
	msg.Title = "Snippet Page"
	msg.Message = fmt.Sprintf("ID: %d, Title: %s, Content: %s, Expires: %v\n", *res, s.Title, s.Content, s.Expires)

	tpl(w, "home.html", msg)
}
