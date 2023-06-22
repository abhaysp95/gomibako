package main

import (
	"html/template"
	"path/filepath"

	"github.com/abhaysp95/gomibako/pkg/models"
)

type templateData struct {
	Gomi *models.Gomi
	GomiList []*models.Gomi
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all the pages in "dir"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parse page template
		ts, err := template.ParseFiles(name)
		if err != nil {
			return nil, err
		}

		// add layouts to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// add partial to template sets
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
