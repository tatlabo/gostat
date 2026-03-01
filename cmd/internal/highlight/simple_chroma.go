package highlight

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func Highlight(someSourceCode string) (string, error) {

	var s strings.Builder

	lexer := lexers.Get("go")
	if lexer == nil {
		return "", fmt.Errorf("lexer not found for 'go'")
	}

	style := styles.Get("friendly")
	if style == nil {
		return "", fmt.Errorf("style 'friendly' not found")
	}

	tokens, err := lexer.Tokenise(nil, someSourceCode)
	if err != nil {
		return "", fmt.Errorf("error tokenizing: %v", err)
	}

	formatter := html.New(html.WithLineNumbers(false), html.WithClasses(true), html.TabWidth(4))

	err = formatter.Format(&s, style, tokens)
	if err != nil {
		return "", fmt.Errorf("error formatting: %v", err)
	}

	return s.String(), nil
}
