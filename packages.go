package main

import (
	"net/http"

	"github.com/steam-authority/steam-authority/datastore"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	template := packagesTemplate{}

	returnTemplate(w, "packages", template)
}

type packagesTemplate struct {
	Packages []*datastore.DsPackage
}
