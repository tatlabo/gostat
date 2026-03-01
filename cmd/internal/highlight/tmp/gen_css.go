package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	var css strings.Builder
	formatter := html.New(html.WithClasses(true))
	style := styles.Get("friendly")
	formatter.WriteCSS(&css, style)
	fmt.Println(css.String())
}
