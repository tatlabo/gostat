package main

import "net/http"

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
