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
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		templateSet, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
