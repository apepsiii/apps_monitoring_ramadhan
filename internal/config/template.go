package config

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer struct untuk Echo
type TemplateRenderer struct {
	templates *template.Template
	embedFS   fs.FS
}

func NewTemplateRenderer(embedFS fs.FS) *TemplateRenderer {
	// Parse all templates with custom functions
	funcMap := template.FuncMap{
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatDateShort": func(date string) string {
			// Parse date and return short format
			t, _ := time.Parse("2006-01-02", date)
			days := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
			return days[int(t.Weekday())]
		},
		"formatDateLong": func(date string) string {
			if date == "" {
				return "-"
			}
			// Try parsing multiple formats
			formats := []string{"2006-01-02", "2006-01-02T15:04:05Z07:00", "2006-01-02 15:04:05"}
			var t time.Time
			var err error
			
			for _, format := range formats {
				t, err = time.Parse(format, date)
				if err == nil {
					break
				}
			}
			
			if err != nil {
				return date // Return original string if parse fails
			}
			
			months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
			days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
			return days[int(t.Weekday())] + ", " + strconv.Itoa(t.Day()) + " " + months[t.Month()-1] + " " + strconv.Itoa(t.Year())
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},
		"divide": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"subtract": func(a, b float64) float64 {
			return a - b
		},
		"calcPercentage": func(current, total int) float64 {
			if total == 0 {
				return 0
			}
			return float64(current) / float64(total) * 100
		},
	}

	tmpl := template.New("").Funcs(funcMap)
	
	// Parse ONLY layouts and shared components from embedded FS
	// Specific pages will be parsed in Render()
	patterns := []string{
		"templates/layouts/*.html",
		"templates/partials/*.html",
		// "templates/*.html", // Don't parse pages yet to avoid block conflicts
		// "templates/auth/*.html",
		// "templates/user/*.html",
		// "templates/admin/*.html",
	}

	for _, pattern := range patterns {
		t, err := tmpl.ParseFS(embedFS, pattern)
		if err != nil {
			fmt.Printf("Error parsing pattern %s: %v\n", pattern, err)
			continue
		}
		tmpl = t
	}

	return &TemplateRenderer{
		templates: tmpl,
		embedFS:   embedFS,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Clone the template set (which has layouts)
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	// Parse the specific template file from FS
	filePath := "templates/" + name
	
	tmpl, err = tmpl.ParseFS(t.embedFS, filePath)
	if err != nil {
		return err
	}

	// Execute the base template
	return tmpl.ExecuteTemplate(w, "base.html", data)
}
