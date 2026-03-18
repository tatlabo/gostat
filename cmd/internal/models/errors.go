package models

import "errors"

var ErrNoRecord = errors.New("No matching record found")
var ErrInvalidCredentials = errors.New("Invalid credentials")
var ErrDuplicateEmail = errors.New("Email already registered")
