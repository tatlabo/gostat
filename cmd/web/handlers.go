package main

import (
	"context"
	"fmt"
	"gostats/cmd/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {

	// panic("panic in hello handler")

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

func (app *Application) error422(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Unprocessable Entity"

	}
	w.WriteHeader(http.StatusUnprocessableEntity)
	app.Template(w, "422", Render{Msg: map[string]string{"Message": msg}})
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

	s := *res

	html, err := Highlight(s.Content)
	if err != nil {
		fmt.Printf("ERROR Highlight Snippet snippets: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable highlight snippets: %v", err), http.StatusInternalServerError)
		return
	}
	s.Html = template.HTML(html)

	m["Message"] = "Snnippet page nr " + strconv.Itoa(id)
	//
	// content := []byte(s.Content)

	err = app.Template(w, "snippet.html", Render{Msg: m, Snippet: s})
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

	// Get flash message if any
	if flash := app.GetFlash(w, r, "success"); flash != "" {
		m["Deleted"] = flash
	}

	res, err := app.Snippets.Latest()
	if err != nil {
		fmt.Printf("ERROR fetching latest snippets: %v\n", err)
		http.Error(w, fmt.Sprintf("Unable to fetch latest snippets: %v", err), http.StatusInternalServerError)
		return
	}

	for i := range res {
		s := &res[i]
		html, err := Highlight(s.Content)
		if err != nil {
			fmt.Printf("ERROR Highlight Snippet snippets: %v\n", err)
			http.Error(w, fmt.Sprintf("Unable highlight snippets: %v", err), http.StatusInternalServerError)
			return
		}
		s.Html = template.HTML(html)
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
	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("Error parsing form: %v", err)
		log.Printf("%v", err)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", msg)
		r = r.WithContext(ctx)
		app.error500(w, r)
		return
	}

	s.Title = r.PostForm.Get("title")
	s.Content = r.PostForm.Get("content")
	expires, err := time.Parse("2006-01-02", r.PostForm.Get("expires"))
	//
	if err != nil {
		msg := ""
		if strings.HasPrefix(err.Error(), "parsing time ") {
			msg = "Invalid date format. Please use YYYY-MM-DD."
		}
		log.Printf("%v", err)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", msg)
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

func (app *Application) snippetDelete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	s := app.Snippet
	//
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		msg := fmt.Sprintf("Error parsing Id: %v", err)
		ctx = context.WithValue(ctx, "error", msg)
		r = r.WithContext(ctx)
		app.error500(w, r)
		return
	}

	s.ID = id

	sDeleteted, err := app.Snippets.Delete(&s)
	if err != nil {
		msg := fmt.Sprintf("ERROR deleting snippet: %v\n%v\n", s.ID, err)
		ctx = context.WithValue(ctx, "error", msg)
		r = r.WithContext(ctx)
		app.error500(w, r)
		return
	}

	// Set flash message for redirect
	app.SetFlash(w, r, "success", "Snippet "+sDeleteted.Title+" deleted")

	http.Redirect(w, r, "/snippet/all", http.StatusSeeOther)

}
