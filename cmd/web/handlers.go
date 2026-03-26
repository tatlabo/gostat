package main

import (
	"context"
	"fmt"
	"gostats/cmd/internal/models"
	"html/template"
	"maps"
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
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Msg             map[string]string
	MsgAlert        map[string]string
	Error           map[string]string
	Flash           string
	IsAuthenticated bool
	AuthName        string
	SignupUserForm
	LoginUserForm
}

func (r *Render) Make() {
	r.Msg = make(map[string]string)
	r.MsgAlert = make(map[string]string)
	r.Error = make(map[string]string)
}

func (r *Render) AddError(k, v string) {
	if r.Error == nil {
		r.Error = make(map[string]string)
	}
	r.Error[k] = v
}

func (r *Render) AddMsg(k, v string) {
	if r.Msg == nil {
		r.Msg = make(map[string]string)
	}
	r.Msg[k] = v
}

func (r *Render) AddMsgAlert(k, v string) {
	if r.MsgAlert == nil {
		r.MsgAlert = make(map[string]string)
	}
	r.MsgAlert[k] = v
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Page not found"
	}
	w.WriteHeader(http.StatusNotFound)
	app.RenderHTML(w, "404", Render{Msg: map[string]string{"Message": msg}})
}

func (app *Application) error500(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Internal Server Error"

	}
	w.WriteHeader(http.StatusInternalServerError)
	app.RenderHTML(w, "500", Render{Msg: map[string]string{"Message": msg}})
}

func (app *Application) error422(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	msg, ok := ctx.Value("error").(string)
	if !ok || msg == "" {
		msg = "Unprocessable Entity"

	}
	w.WriteHeader(http.StatusUnprocessableEntity)
	app.RenderHTML(w, "422", Render{Msg: map[string]string{"Message": msg}})
}

func (app *Application) snippet(w http.ResponseWriter, r *http.Request) {

	m := map[string]string{
		"Title": "Snnippet page",
	}

	idUrl := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idUrl)

	if id == 0 || err != nil {
		app.RenderHTML(w, "home.html", Render{Msg: m})
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

	err = app.RenderHTML(w, "snippet.html", Render{Msg: m, Snippet: s})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

func (app *Application) snippetList(w http.ResponseWriter, r *http.Request) {

	var render Render
	render.Make()
	render.IsAuthenticated = app.IsAuthenticated(r)

	render.AddMsg("Title", "Snippet List Page")
	render.AddMsg("Message", "Snippet list page")

	ctx := r.Context()
	msgError, ok := ctx.Value("error").(map[string]string)
	if ok && msgError != nil {
		maps.Copy(render.Msg, msgError)
	}

	// Get flash message if any
	// if flash := app.GetFlash(w, r, "success"); flash != "" {
	// 	m["Deleted"] = flash
	// }

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
		// temporary
		// s.Html = template.HTML(s.Content)
	}

	render.Snippets = append(render.Snippets, res...)

	fmt.Printf("/n%#v/n", render.IsAuthenticated)

	//
	err = app.RenderHTML(w, "snippet_list.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

// end of snippet form

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	var render Render
	render.AddMsg("Title", "Snippet Page")
	render.IsAuthenticated = app.IsAuthenticated(r)

	switch {
	case r.Method == http.MethodPost && render.IsAuthenticated:
		var fs FormSnippet
		err := r.ParseForm()
		if err != nil {
			fs.AddError("Error", fmt.Sprintf("Error parsing form: %v", err))
			ctx := r.Context()
			ctx = context.WithValue(ctx, "error", fs.FieldError)
			r = r.WithContext(ctx)
			app.snippetList(w, r)
			return
		}

		fs.Title = r.PostForm.Get("title")
		fs.Content = r.PostForm.Get("content")
		fs.Expires, _ = time.Parse("2006-01-02", r.PostForm.Get("expires"))

		fs.CheckForm()
		s := app.Snippet

		if len(fs.FieldError) > 0 {
			fs.FormSendBack()

			ctx := r.Context()
			ctx = context.WithValue(ctx, "error", fs.FieldError)
			r = r.WithContext(ctx)
			app.snippetList(w, r)
			return
		}

		s.Title = fs.Title
		s.Content = fs.Content
		s.Expires = fs.Expires

		res, err := app.Snippets.Insert(s)
		if err != nil {
			fs.FieldError["Error"] = fmt.Sprintf("ERROR inserting snippet: %v\n%v\n", *res, err)
			http.Error(w, fs.FieldError["Error"], http.StatusInternalServerError)
			return
		}

		okMsg := strings.Builder{}
		okMsg.WriteString("Snippet ")
		okMsg.WriteString(s.Title)
		okMsg.WriteString(" created successfully")

		app.SessionManager.Put(r.Context(), "flash", okMsg.String())

		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

	default:

		if !app.IsAuthenticated(r) {
			render.Msg = map[string]string{
				"Message": "Login to create snippet.",
			}

			err := app.RenderHTML(w, "home.html", render)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		app.RenderHTML(w, "snippet-create.html", render)

	}

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

	sDeleteted, err := app.Snippets.Delete(s)
	if err != nil {
		msg := fmt.Sprintf("ERROR deleting snippet: %v\n%v\n", s.ID, err)
		ctx = context.WithValue(ctx, "error", msg)
		r = r.WithContext(ctx)
		app.error500(w, r)
		return
	}

	// Set flash message for redirect
	app.SessionManager.Put(r.Context(), "flash", "Snippet "+sDeleteted.Title+" deleted successfully")

	http.Redirect(w, r, "/snippet/all", http.StatusSeeOther)

}
