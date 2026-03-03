package models

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
	NotBlank    func(string, int) bool
	MaxChar     func(string, int) bool
}

// check is there any error
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// add if not exists
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, ok := v.FieldErrors[key]; !ok {
		v.FieldErrors[key] = message
	}
}

// if check is NOT ok => add field
func (v *Validator) CheckField(ok bool, key, val string) {
	if !ok {
		v.AddFieldError(key, val)
	}
}

// check for empty string
func NotBlank(value string, n int) bool {
	return strings.TrimSpace(value) != ""
}

// max strings
func MaxChar(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

// check if value is in the specific list
func PermittedValues[T comparable](value T, permitted ...T) bool {
	return slices.Contains(permitted, value)
}
