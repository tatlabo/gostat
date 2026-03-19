package main

import (
	"html/template"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

/// snippet form

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func MinChar(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func NotBlank(field *string) bool {
	clean := strings.TrimSpace(*field)
	return len(clean) > 0
}

func Blank(field *string) bool {
	clean := strings.TrimSpace(*field)
	return len(clean) == 0
}

func MaxChar(field *string, n int) bool {
	clean := strings.TrimSpace(*field)
	return utf8.RuneCountInString(clean) <= n
}

func CheckDate(t time.Time) bool {
	return !t.IsZero()
}

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

func (fs *FormSnippet) CheckForm() {
	fs.CheckField(NotBlank(&fs.Title), "ErrorTitle", "Title empty")
	fs.CheckField(NotBlank(&fs.Content), "ErrorContent", "Content empty")
	fs.CheckField(MaxChar(&fs.Title, 150), "ErrorTitle", "Title too long")
	fs.CheckField(MaxChar(&fs.Content, 1500), "ErrorContent", "Content too long")
	fs.CheckField(CheckDate(fs.Expires), "ErrorDate", "Date is required")
}

type userSignupForm struct {
	Name       string            `form:"name"`
	Email      string            `form:"email"`
	Password   string            `form:"password"`
	FieldError map[string]string `form:"-"`
}

type userLoginForm struct {
	Name       string            `form:"name"`
	Email      string            `form:"email"`
	Password   string            `form:"password"`
	FieldError map[string]string `form:"-"`
}

func (f *userSignupForm) AddError(k, v string) {
	if f.FieldError == nil {
		f.FieldError = make(map[string]string)
	}

	if _, ok := f.FieldError[k]; !ok {
		f.FieldError[k] = v
	}
}

func (v *userSignupForm) CheckField(ok bool, key, val string) {
	if !ok {
		v.AddError(key, val)
	}
}

func (fs *userSignupForm) CheckForm() {
	fs.CheckField(!Blank(&fs.Name), "ErrorName", "Name empty")
	fs.CheckField(MaxChar(&fs.Name, 32), "ErrorName", "Name too long, 32 characters maximum")
	fs.CheckField(!Blank(&fs.Email), "ErrorEmail", "Email empty")
	//
	fs.CheckField(Matches(fs.Email, EmailRX), "ErrorEmail", "Email invalid")
	//
	fs.CheckField(!Blank(&fs.Password), "ErrorPassword", "Password empty")
	fs.CheckField(MinChar(fs.Password, 8), "ErrorPassword", "Password too short, 8 characters minimum")
	fs.CheckField(MaxChar(&fs.Password, 32), "ErrorPassword", "Password too long, 32 characters maximum")
}

func (fs *FormSnippet) FormSendBack() {
	fs.FieldError["Error"] = "Error parsing form"
	fs.FieldError["BackTitle"] = template.HTMLEscapeString(fs.Title)
	fs.FieldError["BackContent"] = fs.Content
	fs.FieldError["BackDate"] = fs.Expires.Format("2006-01-02")
}
