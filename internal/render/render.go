package render

import (
	"BnBManagementSystem/internal/config"
	"BnBManagementSystem/internal/models"
	"bytes"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templateCache map[string]*template.Template

func init() {
	templateCache = make(map[string]*template.Template)
}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Template not found in cache:", tmpl)
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	td = AddDefaultData(td, r)
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	pages, err := filepath.Glob("templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}

	for _, p := range pages {
		name := filepath.Base(p)
		ts, err := template.New(name).ParseFiles(p)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob("templates/*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("templates/*.layout.tmpl")
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

func init() {
	var err error
	templateCache, err = CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache:", err)
	}
}
