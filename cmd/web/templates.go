package main

import "github.com/petrostrak/code-snippet/pkg/models"

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
	Snippet *models.Snippet
}
