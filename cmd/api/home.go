package main

import (
	"fmt"
	"net/http"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "welcome home")
}
