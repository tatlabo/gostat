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
	Snippet  models.Snippet
	Snippets []models.Snippet
	Msg      map[string]string
}

type RenderUser struct {
	Msg      map[string]string
	MsgAlert map[string]string
	Error    map[string]string
}

func (ru *RenderUser) Make() {
	ru.Msg = make(map[string]string)
	ru.MsgAlert = make(map[string]string)
	ru.Error = make(map[string]string)
}

func (ru *RenderUser) AddError(k, v string) {
	if ru.Error == nil {
		ru.Error = make(map[string]string)
	}
	ru.Error[k] = v
}

func (ru *RenderUser) AddMsg(k, v string) {
	if ru.Msg == nil {
		ru.Msg = make(map[string]string)
	}
	ru.Msg[k] = v
}

func (ru *RenderUser) AddMsgAlert(k, v string) {
	if ru.MsgAlert == nil {
		ru.MsgAlert = make(map[string]string)
	}
	ru.MsgAlert[k] = v
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

	m := map[string]string{
		"Title":   "Snippet List Page",
		"Message": "Snippet list page",
	}

	ctx := r.Context()
	msgError, ok := ctx.Value("error").(map[string]string)
	if ok && msgError != nil {
		maps.Copy(m, msgError)
	}

	// Get flash message if any
	// if flash := app.GetFlash(w, r, "success"); flash != "" {
	// 	m["Deleted"] = flash
	// }

	flash := app.newTemplateData(r)
	// flash := app.sessionManager.GetString(r.Context(), "flash")

	if flash.Flash != "" {
		m["Flash"] = flash.Flash
		m["FlashTime"] = flash.Time.Format("2006-01-02 15:04:05")
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
		// temporary
		// s.Html = template.HTML(s.Content)
	}
	//
	err = app.RenderHTML(w, "snippet_list.html", Render{Msg: m, Snippets: res})
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}

}

// end of snippet form

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {

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

	res, err := app.Snippets.Insert(&s)
	if err != nil {
		fs.FieldError["Error"] = fmt.Sprintf("ERROR inserting snippet: %v\n%v\n", *res, err)
		http.Error(w, fs.FieldError["Error"], http.StatusInternalServerError)
		return
	}

	okMsg := strings.Builder{}
	okMsg.WriteString("Snippet ")
	okMsg.WriteString(s.Title)
	okMsg.WriteString(" created successfully")

	app.sessionManager.Put(r.Context(), "flash", okMsg.String())

	http.Redirect(w, r, "/snippet/all", http.StatusSeeOther)

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
	app.sessionManager.Put(r.Context(), "flash", "Snippet "+sDeleteted.Title+" deleted successfully")

	http.Redirect(w, r, "/snippet/all", http.StatusSeeOther)

}

/// user handlers

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}
	render.Make()

	render.Msg["Title"] = "User Signup Page"

	flash := app.sessionManager.Pop(r.Context(), "flash")

	if flash != nil {
		render.Msg["Message"] = flash.(string)
		render.MsgAlert["Message"] = "alert alert-warning"
	}

	fieldsError, ok := r.Context().Value("error").(map[string]string)
	if ok && fieldsError != nil {
		render.Error = fieldsError
	}

	err := app.RenderHTML(w, "signup.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	var form userSignupForm
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	form.Name = r.PostForm.Get("name")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	form.CheckForm()
	if len(form.FieldError) > 0 {
		app.sessionManager.Put(r.Context(), "flash", "Error in form submission")
		form.AddError("Error", fmt.Sprintf("Form errors: %v\n", form.FieldError))
		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", form.FieldError)
		r = r.WithContext(ctx)
		app.userSignup(w, r)
		return
	}

	id, err := app.User.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		app.sessionManager.Put(r.Context(), "flash", "Error creating user: "+err.Error())
		app.userSignup(w, r)
		return
	}

	Message := "User " + strconv.Itoa(*id) + " created successfully"

	app.sessionManager.Put(r.Context(), "flash", Message)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	//
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}
	render.Make()

	render.AddMsg("Message", "User Login Page")

	flash := app.sessionManager.Pop(r.Context(), "flash")
	if flash != nil {
		render.AddMsg("Message", flash.(string))
		render.AddMsgAlert("Message", "alert alert-success")
	}

	err := app.RenderHTML(w, "user.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}

	render.AddMsg("Msg", "User Login Page")
	render.AddMsg("Message", "User login page")

	err := app.RenderHTML(w, "user.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}

	// app.sessionManager.Delete(r.Context(), "flash", "Error creating user: "+err.Error())

	render.AddMsg("Title", "User Logout Page")
	render.AddMsg("Message", "User logout successfully")

	err := app.RenderHTML(w, "user.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}
