package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
)

// template related methods for web struct.

func (web *web) templateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := web.initTemplates()
		if err != nil {
			web.errorHandler(w, r, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

// initTemplates initializes templates for use.
// templates are initialized only once, unless web.DevMode is true.
func (web *web) initTemplates() (err error) {
	if web.devMode {
		web.tmplMux.Lock()
		web.templates, err = web.initLocalTemplates()
		web.tmplMux.Unlock()
		return err
	}

	web.tmplOnce.Do(func() {
		web.templates, err = web.initPkgerTemplates()
	})
	return
}

func (*web) templateFuncs() map[string]interface{} {
	return template.FuncMap{
		"SayHello": func(name string) string { return "Hello " + name },
	}
}

func (web *web) initLocalTemplates() (map[string]*template.Template, error) {
	funcMap := web.templateFuncs()
	tmpl := make(map[string]*template.Template)

	// get all the template file paths
	templatePaths, err := filepath.Glob(path.Join("./templates", "*.*"))
	if err != nil {
		return web.templates, err
	}
	// get the template helper file paths (to be included while compiling each template)
	templateHelperPaths, err := filepath.Glob(path.Join("./templates/helpers", "*.*"))
	if err != nil {
		return web.templates, err
	}

	// parse each template
	for _, filePath := range templatePaths {
		name := path.Base(filePath)
		t := template.New(name).Funcs(funcMap)
		paths := append(templateHelperPaths, filePath)
		t, err := t.ParseFiles(paths...)
		if err != nil {
			return tmpl, err
		}
		tmpl[name] = t
	}
	return tmpl, nil
}

func (web *web) initPkgerTemplates() (map[string]*template.Template, error) {
	funcMap := web.templateFuncs()
	tmpl := make(map[string]*template.Template)

	// TODO: Not currently recursive...
	templatePaths, err := fs.Glob(embeded, "templates/*")
	if err != nil {
		return tmpl, err
	}
	// get the template helper file paths (to be included while compiling each template)
	// TODO: Not currently recursive...
	templateHelperPaths, err := fs.Glob(embeded, "templates/helpers/*")
	if err != nil {
		return tmpl, err
	}

	// load and merge helper files.
	hb := []byte{}
	for _, fn := range templateHelperPaths {
		var fb []byte
		fb, err = embeded.ReadFile(fn)
		if err != nil {
			return tmpl, err
		}
		hb = append(hb, fb...)
	}

	// parse each template
	for _, filePath := range templatePaths {
		if filePath == "templates/helpers" {
			// Don't try to load the helper directory
			continue
		}
		name := path.Base(filePath)
		t := template.New(name).Funcs(funcMap)
		b, err := embeded.ReadFile(filePath)
		if err != nil {
			return tmpl, err
		}
		b = append(hb, b...)
		tmpl[name] = template.Must(t.Parse(string(b)))
	}
	return tmpl, nil
}

func (web *web) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {

	web.tmplMux.RLock()
	tmpl, ok := web.templates[name]
	web.tmplMux.RUnlock()
	if !ok {
		web.errorHandler(w, r, fmt.Sprintf("No template found for name: %s", name))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		web.errorHandler(w, r, fmt.Sprintf("Unable to Execute Template %v. Error: %v", name, err))
	}
}
