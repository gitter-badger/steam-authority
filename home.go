package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	template := homeTemplate{}

	returnTemplate(w, "home", template)
}

type homeTemplate struct {
}
