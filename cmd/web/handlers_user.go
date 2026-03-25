package main

import (
	"context"
	"fmt"
	"gostats/cmd/internal/models"
	"net/http"
)

/// user handlers

type RenderUser struct {
	Msg      map[string]string
	MsgAlert map[string]string
	Error    map[string]string
	Flash    string
	SignupUserForm
	LoginUserForm
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

type SignupUserForm struct {
	Name             string `form:"name"`
	Email            string `form:"email"`
	Password         string `form:"password"`
	models.Validator `form:"-"`
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}
	render.Make()

	render.Msg["Title"] = "User Signup Page"

	flash := app.SessionManager.Pop(r.Context(), "Flash")

	if flash != nil {
		render.Flash = flash.(string)
	}

	fieldsError, ok := r.Context().Value("error").(map[string]string)
	if ok && fieldsError != nil {
		fmt.Printf("%#v", fieldsError)
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
	err := form.Resolve(r)

	form.CheckForm()
	if len(form.FieldError) > 0 {
		app.SessionManager.Put(r.Context(), "Flash", "Error in form submission")
		form.AddError("Error", fmt.Sprintf("Form errors: %v\n", form.FieldError))
		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", form.FieldError)
		r = r.WithContext(ctx)
		app.userSignup(w, r)
		return
	}

	_, err = app.User.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		app.SessionManager.Put(r.Context(), "Flash", "Error creating user: "+err.Error())
		app.userSignup(w, r)
		return
	}

	Message := "User " + form.Name + " created successfully"

	app.SessionManager.Put(r.Context(), "flash", Message)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	//
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}
	render.Make()

	fieldsError, ok := r.Context().Value("error").(map[string]string)
	if ok && fieldsError != nil {
		fmt.Printf("%#v", fieldsError)
		render.Error = fieldsError
	}

	render.AddMsg("Message", "User Login Page")
	render.AddMsg("Title", "User Login Page")

	flash := app.SessionManager.Pop(r.Context(), "Flash")
	if flash != nil {
		render.AddMsg("Message", flash.(string))
	}

	err := app.RenderHTML(w, "login.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	form := LoginUserForm{}
	form.Resolve(r)
	form.CheckForm()

	if len(form.FieldError) > 0 {
		app.SessionManager.Put(r.Context(), "Flash", "Error in form submission")
		form.AddError("Error", fmt.Sprintf("Form errors: %v\n", form.FieldError))
		ctx := r.Context()
		ctx = context.WithValue(ctx, "error", form.FieldError)
		r = r.WithContext(ctx)
		app.userLogin(w, r)
		return
	}

	id, err := app.User.Authenticate(form.Email, form.Password)
	if id < 0 {
		app.SessionManager.Put(r.Context(), "Flash", "Authentication failed!")
		app.userLogin(w, r)
		return

	} else if err != nil {
		app.SessionManager.Put(r.Context(), "Flash", "Wrong email/password")
		app.userLogin(w, r)
		return
	}

	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.error500(w, r)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.SessionManager.Put(r.Context(), "Flash", "User OK!")

	app.userLogin(w, r)

	//
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	render := RenderUser{}

	// app.SessionManager.Delete(r.Context(), "flash", "Error creating user: "+err.Error())

	render.AddMsg("Title", "User Logout Page")
	render.AddMsg("Message", "User logout successfully")

	err := app.RenderHTML(w, "user.html", render)
	//
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
		return
	}
}
