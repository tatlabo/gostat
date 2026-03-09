package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Flash message helpers

func Highlight(source string) (string, error) {

	var s strings.Builder

	lexer := lexers.Analyse(source)
	if lexer == nil {
		lexer = lexers.Get("go") // Use plain text fallback
	}
	// if lexer == nil {
	// 	return "", fmt.Errorf("lexer not found for 'go'")
	// }

	tokens, err := lexer.Tokenise(nil, source)
	if err != nil {
		return "", fmt.Errorf("error tokenizing: %v", err)
	}

	formatter := html.New(html.WithLineNumbers(false), html.WithClasses(true), html.TabWidth(4))

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback // Use plain text fallback
	}

	err = formatter.Format(&s, style, tokens)
	if err != nil {
		return "", fmt.Errorf("error formatting: %v", err)
	}

	return s.String(), nil
}
