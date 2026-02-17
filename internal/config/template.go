package config

import (
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer struct untuk Echo
type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
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
			t, _ := time.Parse("2006-01-02", date)
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
	}

	tmpl := template.New("").Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseGlob("web/templates/layouts/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/auth/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/user/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/admin/*.html"))

	return &TemplateRenderer{
		templates: tmpl,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Clone the template set and parse the specific template file
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}
	
	// Parse the specific template file
	filePath := "web/templates/" + name
	tmpl, err = tmpl.ParseFiles(filePath)
	if err != nil {
		return err
	}
	
	// Execute the base template
	return tmpl.ExecuteTemplate(w, "base.html", data)
}
