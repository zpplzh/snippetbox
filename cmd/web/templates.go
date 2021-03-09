package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/zappel/snippetbox/pkg/forms"
	"github.com/zappel/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Flash       string
	//FormData    url.Values
	//FormErrors  map[string]string
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

/*func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	return td

}*/

func newTemplateCache(dir string) (map[string]*template.Template, error) { //return nya map dkk dan error
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
