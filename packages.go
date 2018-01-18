package main

import (
	"net/http"

	"github.com/steam-authority/steam-authority/datastore"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	template := packagesTemplate{}

	returnTemplate(w, "apps", template)
}

type packagesTemplate struct {
	Packages []*datastore.DsPackage
}
