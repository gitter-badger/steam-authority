package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := datastore.GetLatestUpdatedApps(96)
	if err != nil {
		logger.Error(err)
	}

	// Get apps count
	count, err := datastore.CountApps()
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := appsTemplate{}
	template.Apps = apps
	template.Count = count

	returnTemplate(w, "apps", template)
}

type appsTemplate struct {
	GlobalTemplate
	Apps  []datastore.DsApp
	Count int
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	// Get app
	app, err := datastore.GetApp(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)
		if err.Error() == "datastore: no such entity" {
			returnErrorTemplate(w, 404, "We can't find this app in our database, there may not be one with this ID.")
			return
		}
	}

	// Get packages
	packages, err := datastore.GetPackagesAppIsIn(app.AppID)
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := appTemplate{}
	template.App = app
	template.Packages = packages

	returnTemplate(w, "app", template)
}

type appTemplate struct {
	GlobalTemplate
	App      datastore.DsApp
	Packages []datastore.DsPackage
}
