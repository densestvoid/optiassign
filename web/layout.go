package web

import (
	"html/template"
	"net/http"
	"optiassign/domain"
)

// LayoutData holds data for the main layout
type LayoutData struct {
	Title   string
	Content template.HTML
	Session *domain.Session
}

// PageData holds data for individual pages
type PageData struct {
	Title   string
	Content template.HTML
	Session *domain.Session
	Data    interface{}
}

// Templates holds all parsed templates
var Templates *template.Template

// InitTemplates initializes and parses all templates
func InitTemplates() error {
	var err error
	Templates, err = template.ParseGlob("web/templates/*.html")
	return err
}

// RenderTemplate renders a template with the given data
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	return Templates.ExecuteTemplate(w, tmpl, data)
}