// internal/server/helpers.go
package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/logger"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"go.uber.org/zap"
)

func RenderTemplate(w http.ResponseWriter, templatePath string, templateName string, data interface{}) {
	// Define common template functions
	funcMap := template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"GetSavingThrowModifiers": charRules.GetSavingThrowModifiers,
		"add": func(a, b interface{}) int64 {
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v + int64(w)
				case int64:
					return v + w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) + w
				case int:
					return int64(v + w)
				}
			}
			return 0
		},
		"mul": func(a, b interface{}) int64 {
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v * int64(w)
				case int64:
					return v * w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) * w
				case int:
					return int64(v * w)
				}
			}
			return 0
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"sub": func(a, b interface{}) int64 {
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v - int64(w)
				case int64:
					return v - w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) - w
				case int:
					return int64(v - w)
				}
			}
			return 0
		},
		"abs": func(x int) int {
			if x < 0 {
				return -x
			}
			return x
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("January 2, 2006 3:04 PM")
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"formatModifier": func(mod int) string {
			if mod > 0 {
				return "+" + strconv.Itoa(mod)
			}
			return strconv.Itoa(mod)
		},
		"contains": containsString,
	}

	var tmpl *template.Template
	var err error

	// Special handling for character details template which requires multiple nested templates
	if templatePath == "templates/characters/details.html" {
		// Parse all required templates for character details
		tmpl, err = template.New("base.html").Funcs(funcMap).ParseFiles(
			"templates/layout/base.html",
			templatePath,
			"templates/characters/_character_header.html",
			"templates/characters/_inventory.html",
			"templates/characters/_ability_scores.html",
			"templates/characters/_class_features.html",
			"templates/characters/_combat_stats.html",
			"templates/characters/_saving_throws.html",
			"templates/characters/_hp_display.html",
			"templates/characters/_hp_section.html",
			"templates/characters/_currency_section.html",
			"templates/characters/inventory_modal.html",
		)
	} else if strings.HasSuffix(templatePath, "base.html") {
		// For base templates, parse just the single file
		tmpl, err = template.New("base.html").Funcs(funcMap).ParseFiles(templatePath)
	} else {
		// For templates that extend base.html, parse both files together
		tmpl, err = template.New("base.html").Funcs(funcMap).ParseFiles("templates/layout/base.html", templatePath)
	}

	if err != nil {
		logger.Error("Template parsing error", zap.Error(err), zap.String("template", templatePath))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err), zap.String("template", templateName))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RespondWithError logs and returns an HTTP error
func RespondWithError(w http.ResponseWriter, message string, status int, err error) {
	logger.Error(message, zap.Error(err))
	http.Error(w, message, status)
}

// ParseIntField extracts and parses an integer from a form field
func ParseIntField(r *http.Request, fieldName string) (int64, error) {
	return strconv.ParseInt(r.FormValue(fieldName), 10, 64)
}

// RedirectWithMessage redirects with a query parameter message
func RedirectWithMessage(w http.ResponseWriter, r *http.Request, path, message string, status int) {
	http.Redirect(w, r, fmt.Sprintf("%s?message=%s", path, url.QueryEscape(message)), status)
}
