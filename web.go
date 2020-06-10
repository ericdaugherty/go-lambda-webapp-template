package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
)

type web struct {
	devMode   bool
	tmpl      *template.Template
	tmplOnce  sync.Once
	tmplMux   sync.RWMutex
	templates map[string]*template.Template
}

func (web *web) home(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	templateData := make(map[string]interface{})
	templateData["timestamp"] = time.Now()
	templateData["name"] = name

	web.renderTemplate(w, r, "hello.html", templateData)
}

func (web *web) errorHandler(w http.ResponseWriter, r *http.Request, errorDesc string) {
	reqID := middleware.GetReqID(r.Context())
	log.Printf("[%v] Server Error: %v", reqID, errorDesc)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if web.devMode {
		fmt.Fprintf(w, "Server Error: %v", errorDesc)
	} else {
		fmt.Fprintf(w, "Server Error")
	}
}
