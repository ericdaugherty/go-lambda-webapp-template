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

// initTemplates initializes templates for use.
// templates are initialized only once, unless web.DevMode is true.
func (web *web) initTemplates() (err error) {
	if web.templates != nil && !web.devMode {
		return
	}

	if web.templates == nil {
		web.templates = make(map[string]*template.Template)
	}

	// add functions for use within the templates
	funcMap := template.FuncMap{
		"SayHello": func(name string) string { return "Hello " + name },
	}

	if web.devMode {
		err = web.initLocalTemplates(funcMap)
	} else {
		err = web.initPkgerTemplates(funcMap)
	}
	return

}

func (web *web) initLocalTemplates(funcMap template.FuncMap) error {
	// get all the template file paths
	templatePaths, err := filepath.Glob(path.Join("./templates", "*.*"))
	if err != nil {
		return err
	}
	// get the template helper file paths (to be included while compiling each template)
	templateHelperPaths, err := filepath.Glob(path.Join("./templates/helpers", "*.*"))
	if err != nil {
		return err
	}

	// parse each template
	for _, filePath := range templatePaths {
		name := path.Base(filePath)
		t := template.New(name).Funcs(funcMap)
		paths := append(templateHelperPaths, filePath)
		web.templates[name] = template.Must(t.ParseFiles(paths...))
	}
	return nil
}

func (web *web) initPkgerTemplates(funcMap template.FuncMap) error {
	// tell pkger to include the templates folder
	pkger.Include("/templates")

	// get all the template file paths
	templatePaths, err := web.globPkger("/templates")
	if err != nil {
		return err
	}
	// get the template helper file paths (to be included while compiling each template)
	templateHelperPaths, err := web.globPkger("/templates/helpers")
	if err != nil {
		return err
	}

	// load the helper files once.
	hb, err := web.readPkgerFiles(templateHelperPaths)

	// parse each template
	for _, filePath := range templatePaths {
		name := path.Base(filePath)
		t := template.New(name).Funcs(funcMap)
		b, err := web.readPkgerFiles([]string{filePath})
		if err != nil {
			return err
		}
		b = append(hb, b...)
		web.templates[name] = template.Must(t.Parse(string(b)))
	}
	return nil
}

func (web *web) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {

	tmpl, ok := web.templates[name]
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
