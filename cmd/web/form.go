package main

import (
	"html/template"
	"strings"
	"time"
	"unicode/utf8"
)

/// snippet form

type FormSnippet struct {
	Title      string            `form:"title"`
	Content    string            `form:"content"`
	Expires    time.Time         `form:"expire"`
	Html       template.HTML     `form:"-"`
	FieldError map[string]string `form:"-"`
}

func (f *FormSnippet) AddError(k, v string) {
	if f.FieldError == nil {
		f.FieldError = make(map[string]string)
	}

	if _, ok := f.FieldError[k]; !ok {
		f.FieldError[k] = v
	}
}

func (v *FormSnippet) CheckField(ok bool, key, val string) {
	if !ok {
		v.AddError(key, val)
	}
}

func NotBlank(field *string) bool {
	clean := strings.TrimSpace(*field)
	return len(clean) > 0
}

func MaxChar(field *string, n int) bool {
	clean := strings.TrimSpace(*field)
	return utf8.RuneCountInString(clean) <= n
}

func CheckDate(t time.Time) bool {
	return !t.IsZero()
}

func (fs *FormSnippet) CheckForm() {
	fs.CheckField(NotBlank(&fs.Title), "ErrorTitle", "Title empty")
	fs.CheckField(NotBlank(&fs.Content), "ErrorContent", "Content empty")
	fs.CheckField(MaxChar(&fs.Title, 150), "ErrorTitle", "Title too long")
	fs.CheckField(MaxChar(&fs.Content, 1500), "ErrorContent", "Content too long")
	fs.CheckField(CheckDate(fs.Expires), "ErrorDate", "Date is required")
}

func (fs *FormSnippet) FormSendBack() {
	fs.FieldError["Error"] = "Error parsing form"
	fs.FieldError["BackTitle"] = template.HTMLEscapeString(fs.Title)
	fs.FieldError["BackContent"] = fs.Content
	fs.FieldError["BackDate"] = fs.Expires.Format("2006-01-02")
}
