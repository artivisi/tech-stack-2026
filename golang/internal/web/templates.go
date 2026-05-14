package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

type Templates struct {
	tmpls map[string]*template.Template
}

func NewTemplates() (*Templates, error) {
	pages, err := fs.Glob(templatesFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	const layout = "templates/layout.html"
	tmpls := map[string]*template.Template{}
	for _, page := range pages {
		if page == layout {
			continue
		}
		name := strings.TrimSuffix(filepath.Base(page), ".html")
		t, err := template.New(filepath.Base(layout)).Funcs(funcMap()).ParseFS(templatesFS, layout, page)
		if err != nil {
			return nil, err
		}
		tmpls[name] = t
	}
	return &Templates{tmpls: tmpls}, nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"goVersion": func() string {
			return strings.TrimPrefix(runtime.Version(), "go")
		},
	}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data any) error {
	tmpl, ok := t.tmpls[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "layout", data)
}
