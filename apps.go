package main

import (
	"net/http"

	"github.com/steam-authority/steam-authority/datastore"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	template := appsTemplate{}

	returnTemplate(w, "apps", template)
}

type appsTemplate struct {
	Apps []*datastore.DsApp
}
