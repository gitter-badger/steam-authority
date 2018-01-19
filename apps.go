package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	apps, err := datastore.GetLatestUpdatedApps(10)
	if err != nil {
		logger.Error(err)
	}

	template := appsTemplate{}
	template.Apps = apps

	returnTemplate(w, "apps", template)
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	app, err := datastore.GetApp(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)
		if err.Error() == "datastore: no such entity" {
			returnErrorTemplate(w, 404, "We can't find this change in our database, there may not have been one with this ID.")
			return
		}
	}

	template := appTemplate{}
	template.App = app

	returnTemplate(w, "change", template)
}

type appsTemplate struct {
	Apps []datastore.DsApp
}

type appTemplate struct {
	App datastore.DsApp
}
