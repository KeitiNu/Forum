package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/data"
	"git.01.kood.tech/roosarula/forum/pkg/forms"
)

type templateData struct {
	AuthenticatedUser *data.User
	Form              *forms.Form
	Post              *data.Post
	Posts             []*data.Post
	Comments          []*data.Comment
	User              *data.User
	Users             []*data.User
	Categories        []*data.Category
	UserVotes         [][]int
	Sort              string
}

var functions = template.FuncMap{
	"timeAgo":   timeAgo,
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("Mon 02 Jan 2006 15:04:05 MST")
}

// Converts content creation time to a string "x (minutes/hours/...) ago"
func timeAgo(t time.Time) string {
	s := time.Now()
	timeDiff := int(s.Sub(t).Minutes())
	switch {
	case timeDiff >= 525600:
		return fmt.Sprintf("%d year ago", timeDiff/525600)
	case timeDiff >= 2880:
		return fmt.Sprintf("%d days ago", timeDiff/1440)
	case timeDiff >= 1440:
		return fmt.Sprintf("%d day ago", timeDiff/1440)
	case timeDiff > 120:
		return fmt.Sprintf("%d hours ago", timeDiff/60)
	case timeDiff > 60:
		return fmt.Sprintf("%d hour ago", timeDiff/60)
	case timeDiff > 1:
		return fmt.Sprintf("%d minutes ago", timeDiff)
	default:
		return fmt.Sprintf("%d minute ago", timeDiff)

	}

}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
