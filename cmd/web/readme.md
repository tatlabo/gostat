



type Application struct {
    Models 
    Model
    Template func(wr io.Writer, name string, data any) error
}



main()
    http.ListenAndServe(*address, routes)






type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (*int, error) {
    ....
}






*Application Routes() http.Handler 
    mux := http.NewServeMux()
    fileServer = http.FileServer(http.Dir("static"))
    mux.Handle("/static", http.StripPrefix("/static/"), fileServer)

    mux.Handle("GET /path", http.HandlerFunc)






func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Execute our middleware logic here...
		next.ServeHTTP(w, r)
	})
}



func customTemplate() (*template.Template, error) {

	parse, err := template.New("").Funcs(funcMap()).ParseGlob("./cmd/ui/html/*.html")
	if err != nil {
		return nil, err
	}

	return parse, nil
}




func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%v", err)
				log.Printf("%v", err)
				ctx := r.Context()
				ctx = context.WithValue(ctx, "error", msg)
				r = r.WithContext(ctx)
				app.error500(w, r)
			}
		}()

		next.ServeHTTP(w, r)

	})
}
