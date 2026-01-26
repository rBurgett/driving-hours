package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"driving-hours/internal/middleware"
	"driving-hours/internal/models"
)

type Data map[string]interface{}

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer(templatesDir string) (*Renderer, error) {
	templates := make(map[string]*template.Template)
	funcMap := createFuncMap()

	// Parse base layout
	baseLayout := filepath.Join(templatesDir, "layouts", "base.html")

	// Parse partials
	partialsGlob := filepath.Join(templatesDir, "partials", "*.html")
	partials, err := filepath.Glob(partialsGlob)
	if err != nil {
		return nil, fmt.Errorf("failed to glob partials: %w", err)
	}

	// Walk through page templates
	err = filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Skip layouts and partials
		rel, _ := filepath.Rel(templatesDir, path)
		if strings.HasPrefix(rel, "layouts") || strings.HasPrefix(rel, "partials") {
			return nil
		}

		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Create template with base layout, partials, and page
		files := append([]string{baseLayout}, partials...)
		files = append(files, path)

		tmpl, err := template.New(filepath.Base(baseLayout)).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		// Use relative path as key (e.g., "auth/login.html")
		templates[rel] = tmpl
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Renderer{templates: templates}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, req *http.Request, name string, data Data) {
	tmpl, ok := r.templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("template %s not found", name), http.StatusInternalServerError)
		return
	}

	// Add CSRF token to all templates
	if data == nil {
		data = make(Data)
	}
	data["CSRFField"] = template.HTML(middleware.CSRFTemplateField(req))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("Jan 2, 2006 3:04 PM")
		},
		"formatHours": func(h float64) string {
			hours := int(h)
			minutes := int((h - float64(hours)) * 60)
			if minutes == 0 {
				return fmt.Sprintf("%dh", hours)
			}
			return fmt.Sprintf("%dh %dm", hours, minutes)
		},
		"formatDecimal": func(h float64) string {
			return fmt.Sprintf("%.1f", h)
		},
		"percentage": func(current, required float64) float64 {
			if required <= 0 {
				return 0
			}
			p := (current / required) * 100
			if p > 100 {
				return 100
			}
			return p
		},
		"add": func(a, b int) int {
			return a + b
		},
		"addFloat": func(a, b float64) float64 {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"isAdmin": func(user *models.User) bool {
			return user != nil && user.IsAdmin()
		},
		"isDriver": func(user *models.User) bool {
			return user != nil && user.IsDriver()
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			d := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil
				}
				d[key] = values[i+1]
			}
			return d
		},
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
	}
}
