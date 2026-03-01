package highlight

import (
	"log/slog"
	"testing"
)

func TestHighlight(t *testing.T) {
	testStrin := "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}"
	h, err := Highlight(testStrin)

	slog.Info("highlighted code", "code", h)

	if err != nil {
		t.Fatal(err)
	}
}
