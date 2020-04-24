package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonHelloWorld(w http.ResponseWriter, r *http.Request) {

	x := struct {
		Hello string `json:"hello"`
	}{
		Hello: "World",
	}

	b, err := json.Marshal(x)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Server Error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
