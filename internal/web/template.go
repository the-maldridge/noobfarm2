package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

type renderer struct {
	log hclog.Logger

	templates *template.Template
}

// NewRenderer constructs a new echo compatible renderer.
func NewRenderer(l hclog.Logger) *renderer {
	x := new(renderer)
	x.log = l.Named("template")
	return x
}

// Reload loads all templates again.
func (r *renderer) Reload() {
	newTmpl := template.New("base")
	if _, err := newTmpl.ParseGlob("web/fragments/*.tpl"); err != nil {
		r.log.Error("Error parsing fragments", "error", err)
	}

	if _, err := newTmpl.ParseGlob("web/layouts/*.tpl"); err != nil {
		r.log.Error("Error parsing layouts", "error", err)
	}
	r.templates = newTmpl
}

func (r renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var page bytes.Buffer
	err := r.templates.ExecuteTemplate(&page, name, data)
	if err != nil {
		r.log.Error("Template runtime error", "error", err)
		return err
	}

	html := strings.Replace(page.String(), "\\n", "<br />", -1)
	fmt.Fprint(w, html)

	return nil
}
