package main

import (
	"fmt"
	"gostats/cmd/ui"
	"html/template"
	"time"
)

var funcMap = func() template.FuncMap {

	return template.FuncMap{
		"mod":         mod,
		"sub":         sub,
		"CurrentYear": currentYear,
		"CurrentDay":  currentDate,
	}
}

func customTemplate() (*template.Template, error) {

	parse, err := template.New("").Funcs(funcMap()).ParseFS(ui.Html, "html/*.html")
	if err != nil {
		return nil, err
	}

	return parse, nil
}

func mod(i, j int) int {
	return i % j
}
func sub(i, j int) int {
	return i - j
}

func currentDay() int {
	return time.Now().Day()
}

func currentYear() int {
	return time.Now().Year()
}

func currentMonth() time.Month {
	return time.Now().Month()
}

func currentDate() string {
	return time.Now().Format("2006-01-02")
}

func humanDate(t time.Time) string {

	var polisMonths = [12]string{
		"sty",
		"lut",
		"mar",
		"kwi",
		"maj",
		"cze",
		"lip",
		"sie",
		"wrz",
		"paź",
		"lis",
		"gru",
	}

	return fmt.Sprintf("%d-%02s-%d", t.Year(), polisMonths[int(t.Month())-1], t.Day())
}
