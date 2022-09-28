package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/justinas/nosurf"
)

var pathToTemplates = "./static/templates"

func NewRenderer(a *AppConfig) {
	app = *a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {

	td.CSRFToken = nosurf.Token(r)

	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()

	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("could not get template from template cache")

		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
		return err
	}

	return nil

}

// Map of functions available to the templates
var functions = template.FuncMap{
	"humanDate":        HumanDate,
	"formatDate":       FormatDate,
	"iterate":          Iterate,
	"add":              Add,
	"divide":           Divide,
	"resizeVolume":     ResizeVolume,
	"split":            strings.Split,
	"removeDuplicates": RemoveDuplicateStr,
}

// HumanDate returns time in yyyy-mm-dd format
func HumanDate(t time.Time) string {

	return t.Format("2006-01-02")
}

// FormatDate helps to format numeric date into string date
func FormatDate(t time.Time, f string) string {

	return t.Format(f)
}

// Add add two numbers
func Add(a, b int) int {

	return a + b
}

// Divide divides two numbers and returns the result as a string
func Divide(a, b uint64) string {

	if b == 0 {
		return "0"
	}

	floatA := float64(a)

	floatB := float64(b)

	result := floatA / floatB

	stringresult := fmt.Sprintf("%.0f", result)

	return stringresult

}

// Iterate performs a for loop
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

// Creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}
