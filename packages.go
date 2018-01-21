package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	packages, err := datastore.GetLatestUpdatedPackages(10)
	if err != nil {
		logger.Error(err)
	}

	template := packagesTemplate{}

	template.Packages = packages

	returnTemplate(w, "packages", template)
}

func packageHandler(w http.ResponseWriter, r *http.Request) {

	packagex, err := datastore.GetPackage(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)

		if err.Error() == "datastore: no such entity" {
			returnErrorTemplate(w, 404, "We can't find this package in our database, there may not be one with this ID.")
			return
		}

		returnErrorTemplate(w, 404, err.Error())
		return
	}

	apps, err := datastore.GetMultiAppsByKey(packagex.Apps)
	if err != nil {
		logger.Error(err)
	}

	template := packageTemplate{}
	template.Package = packagex
	template.Apps = apps

	returnTemplate(w, "package", template)
}

type packagesTemplate struct {
	GlobalTemplate
	Packages []datastore.DsPackage
}

type packageTemplate struct {
	GlobalTemplate
	Package datastore.DsPackage
	Apps    []datastore.DsApp
}
