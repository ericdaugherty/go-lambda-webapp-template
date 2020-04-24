package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type web struct {
	devMode   bool
	tmpl      *template.Template
	templates map[string]*template.Template
}

func (web *web) home(w http.ResponseWriter, r *http.Request) {
	// called on every invocation to enable hot-reload in devMode
	err := web.initTemplates()
	if err != nil {
		web.errorHandler(w, r, err.Error())
		return
	}

	name := r.URL.Query().Get("name")

	templateData := make(map[string]interface{})
	templateData["timestamp"] = time.Now()
	templateData["name"] = name

	web.renderTemplate(w, r, "hello.html", templateData)
}

func (*web) errorHandler(w http.ResponseWriter, r *http.Request, errorDesc string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Server Error: %v", errorDesc)
}
