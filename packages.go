package main

import (
	"net/http"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	template := packagesTemplate{}

	returnTemplate(w, "apps", template)
}

type packagesTemplate struct {
	Packages []*dsPackage
}
