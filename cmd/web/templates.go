package main

import (
	"html/template"
	"path/filepath"
	"toramanomer/snippetbox/internal/models"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		templateSet, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		cache[name] = templateSet
	}

	return cache, nil
}

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
