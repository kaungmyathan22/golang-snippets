package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/kaungmyathan22/golang-sinppets/pkg/forms"
	"github.com/kaungmyathan22/golang-sinppets/pkg/models"
)

type templateData struct {
	AuthenticatedUser int
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	User              *models.User
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	CSRFToken         string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
