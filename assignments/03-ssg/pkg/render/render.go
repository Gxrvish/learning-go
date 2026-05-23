// Package render is the public API for rendering pages to HTML.
package render

import (
	"html/template"
	"io"
	"time"
)

// Page is one input document with parsed front matter and rendered body.
type Page struct {
	Title string
	Date  time.Time
	Tags  []string
	Body  template.HTML
	Path  string
}

// Renderer turns a Page into HTML written to w.
type Renderer interface {
	Render(w io.Writer, p Page) error
}
