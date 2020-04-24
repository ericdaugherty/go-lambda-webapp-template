package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/markbates/pkger"
)

func handler(devMode bool) http.Handler {
	web := web{devMode: devMode}

	// setup static FileServer to use pkger on AWS or local files when running locally
	var dir http.FileSystem = pkger.Dir("/public")
	if devMode {
		dir = http.Dir("./public")
	}
	public := http.FileServer(dir)

	r := chi.NewRouter()

	// add basic middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// configure application specific handlers.
	r.Get("/json/hello", jsonHelloWorld)
	r.Get("/tmpl/hello", web.home)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		public.ServeHTTP(w, r)
	})

	return r
}

func main() {

	// check and see if we are running within AWS.
	aws := len(os.Getenv("AWS_REGION")) > 0

	http.Handle("/", handler(!aws))

	// run using apex gateway on Lambda, or just plain net/http locally
	if aws {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	} else {
		log.Println("Starting listener http://localhost:3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}
}
