package main

import (
	"html/template"
	"path/filepath"

	"github.com/petrostrak/code-snippet/pkg/models"
)

// An important thing to explain is that Go’s html/template package allows
// you to pass in one — and only one — item of dynamic data when
// rendering a template. But in a real-world application there are often
// multiple pieces of dynamic data that you want to display in the same
// page.
//
// A lightweight and type-safe way to acheive this is to wrap your dynamic
// data in a struct which acts like a single ‘holding structure’ for your data.

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Each and every time we render a web page, our application must read
// the template files from disk. This could be speeded up by caching the
// templates in memory. Create an in-memory map with type
// map[string]*template.Template to cache the templates.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() to get a slice of all filepaths with the extension
	// '.page.tmpl'. This essentially gives us a slice of all the 'page' templates
	// for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one
	for _, page := range pages {

		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page templage file in to a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob() to add any 'layout' templates to the template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
