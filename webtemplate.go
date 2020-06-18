package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
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

	// tell pkger to include the templates folder
	pkger.Include("/templates")

	// get all the template file paths
	templatePaths, err := web.globPkger("/templates")
	if err != nil {
		return tmpl, err
	}
	// get the template helper file paths (to be included while compiling each template)
	templateHelperPaths, err := web.globPkger("/templates/helpers")
	if err != nil {
		return tmpl, err
	}

	// load the helper files once.
	hb, err := web.readPkgerFiles(templateHelperPaths)

	// parse each template
	for _, filePath := range templatePaths {
		if filePath == "/templates/helpers" {
			// Don't try to load the helper directory
			continue
		}
		name := path.Base(filePath)
		t := template.New(name).Funcs(funcMap)
		b, err := web.readPkgerFiles([]string{filePath})
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

// globPkger is a simple (non pattern matching) version of glob for Pkger
func (*web) globPkger(dir string) (m []string, err error) {
	fi, err := pkger.Stat(dir)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return
	}
	d, err := pkger.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdir(-1)
	if err != nil {
		return
	}

	for _, n := range names {
		m = append(m, path.Join(dir, n.Name()))
	}

	return
}

// readPkgerFiles reads all the files in the specified slice and returns a
// single []byte containing the full contents
func (*web) readPkgerFiles(fns []string) (b []byte, err error) {

	for _, fn := range fns {
		var f pkging.File
		f, err = pkger.Open(fn)
		if err != nil {
			return
		}
		var fb []byte
		fb, err = ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			return
		}
		b = append(b, fb...)
	}

	return
}
