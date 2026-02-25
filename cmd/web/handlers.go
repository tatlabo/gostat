package main

import (
	"errors"
	"fmt"
	"gostats/cmd/internal/models"
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

func (app *application) snippet(w http.ResponseWriter, r *http.Request) {

	var id int
	var msg msg
	msg.Title = "Snippet Page"

	idUrl := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idUrl)

	if id == 0 || err != nil {
		msg.Message = "Snnippet page"
		tpl(w, "home.html", msg)
		return
	}

	res, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
			return
		}
		fmt.Printf("ERROR fetching snippet: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to fetch snippet: %v", err), http.StatusInternalServerError)
		return
	}

	type SnippetData struct {
		Snippet models.Snippet
		Title   string
		Msg     string
	}

	var data SnippetData
	data.Snippet = *res
	data.Title = "Snippet Page"

	data.Msg = "Snnippet page nr " + strconv.Itoa(id)
	tpl(w, "snippet.html", data)

}

func (app *application) snippetList(w http.ResponseWriter, r *http.Request) {

	var msg msg
	msg.Title = "Snippet Page"
	msg.Message = "Snnippet list page"

	type Render struct {
		Snippets []models.Snippet
		Title    string
		Msg      string
	}

	var render Render

	res, err := app.Snippets.Latest()
	if err != nil {
		fmt.Printf("ERROR fetching latest snippets: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to fetch latest snippets: %v", err), http.StatusInternalServerError)
		return
	}

	render.Snippets = res
	render.Title = "Snippet List Page"
	render.Msg = "Snippet list page"
	tpl(w, "snippet_list.html", render)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	var s = app.Snippet
	s.Title = r.FormValue("title")
	s.Content = r.FormValue("content")
	s.Expires = r.FormValue("expires")

	res, err := app.Snippets.Insert(&s)
	if err != nil {
		fmt.Printf("ERROR inserting snippet: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to create snippet: %v", err), http.StatusInternalServerError)
		return
	}

	var msg msg
	msg.Title = "Snippet Page"
	msg.Message = fmt.Sprintf("ID: %d, Title: %s, Content: %s, Expires: %v\n", *res, s.Title, s.Content, s.Expires)

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", *res), http.StatusSeeOther)

}
