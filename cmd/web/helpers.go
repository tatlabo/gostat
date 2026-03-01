package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Flash message helpers
const sessionName = "flash-session"

func (app *Application) SetFlash(w http.ResponseWriter, r *http.Request, key, value string) error {
	session, err := app.Session.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.AddFlash(value, key)
	return session.Save(r, w)
}

func (app *Application) GetFlash(w http.ResponseWriter, r *http.Request, key string) string {
	session, err := app.Session.Get(r, sessionName)
	if err != nil {
		return ""
	}
	flashes := session.Flashes(key)
	if len(flashes) == 0 {
		return ""
	}
	// Save to clear the flash
	session.Save(r, w)
	return fmt.Sprintf("%v", flashes[0])
}

func Highlight(someSourceCode string) (string, error) {

	var s strings.Builder

	lexer := lexers.Get("go")
	if lexer == nil {
		return "", fmt.Errorf("lexer not found for 'go'")
	}

	tokens, err := lexer.Tokenise(nil, someSourceCode)
	if err != nil {
		return "", fmt.Errorf("error tokenizing: %v", err)
	}

	formatter := html.New(html.WithLineNumbers(false), html.WithClasses(true), html.TabWidth(4))

	style := styles.Get("friendly")
	if style == nil {
		return "", fmt.Errorf("style 'friendly' not found")
	}

	err = formatter.Format(&s, style, tokens)
	if err != nil {
		return "", fmt.Errorf("error formatting: %v", err)
	}

	return s.String(), nil
}

