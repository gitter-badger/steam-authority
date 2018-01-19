package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	template := HomeTemplate{}

	returnTemplate(w, "home", template)
}

type HomeTemplate struct {
}
