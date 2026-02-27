package main

import (
	"context"
	"fmt"
	"gostats/cmd/internal/models"
	"net/http"
	"strconv"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {

	var render Render
	render.Msg = map[string]string{
		"Title":   "Snippet Page",
		"Message": "Hello, World! Everyone loves Go!",
	}

	customTpl, err := customTemplate()
	tmpl := customTpl.ExecuteTemplate

	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl(w, "home.html", render)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

type Render struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
	Msg      map[string]string
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Page not found"
	}
	w.WriteHeader(http.StatusNotFound)
	app.Template(w, "404", Render{Msg: map[string]string{"Message": msg}})
}

func (app *Application) error500(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Internal Server Error"

	}
	w.WriteHeader(http.StatusInternalServerError)
	app.Template(w, "500", Render{Msg: map[string]string{"Message": msg}})
}

func (app *Application) snippet(w http.ResponseWriter, r *http.Request) {

	m := map[string]string{
		"Title": "Snnippet page",
	}

	idUrl := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idUrl)

	if id == 0 || err != nil {
		app.Template(w, "home.html", Render{Msg: m})
		return
	}

	res, err := app.Snippets.Get(id)
	if err != nil {

		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", fmt.Sprintf("Unable to find snippet %d: %v", id, err))
		r = r.WithContext(ctx)
		app.notFound(w, r)

		return
	}

	m["Message"] = "Snnippet page nr " + strconv.Itoa(id)
	//

	err = app.Template(w, "snippet.html", Render{Msg: m, Snippet: *res})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

func (app *Application) snippetList(w http.ResponseWriter, r *http.Request) {

	m := map[string]string{
		"Title":   "Snippet List Page",
		"Message": "Snippet list page",
	}

	res, err := app.Snippets.Latest()
	if err != nil {
		fmt.Printf("ERROR fetching latest snippets: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to fetch latest snippets: %v", err), http.StatusInternalServerError)
		return
	}

	//
	err = app.Template(w, "snippet_list.html", Render{Msg: m, Snippets: res})
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	s := app.Snippet
	//
	s.Title = r.FormValue("title")
	s.Content = r.FormValue("content")
	expires, err := time.Parse("2006-01-02", r.FormValue("expires"))
	//
	if err != nil {

		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", fmt.Sprintf("Unable to inseet:  %v", err))
		r = r.WithContext(ctx)
		app.error500(w, r)

		return

	}
	s.Expires = expires

	res, err := app.Snippets.Insert(&s)
	if err != nil {
		msg := fmt.Sprintf("ERROR inserting snippet: %v\n%v\n", *res, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", *res), http.StatusSeeOther)

}
