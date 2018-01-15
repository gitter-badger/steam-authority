package main

import (
	"net/http"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	template := appsTemplate{}

	returnTemplate(w, "apps", template)
}

type appsTemplate struct {
	Apps []*dsApp
}
