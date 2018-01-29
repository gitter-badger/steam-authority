package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/steam"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := datastore.SearchApps(r.URL.Query(), 96)
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

	id := chi.URLParam(r, "id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	// Get app
	dsApp, err := datastore.GetApp(id)
	if err != nil {
		if err.Error() == "datastore: no such entity" {

			dsApp.AppID = idx

			// Get app details
			details, err := steam.GetAppDetails(id)
			if err != nil {
				if err.Error() == "no app with id" {
					returnErrorTemplate(w, 404, "Sorry but there is no app with this ID")
					return
				}
				logger.Error(err)
			}
			dsApp.FillFromAppDetails(details)

			dsApp.Tidy()
			_, err = datastore.SaveKind(dsApp.GetKey(), &dsApp)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Error(err)
			returnErrorTemplate(w, 404, err.Error())
			return
		}
	}
	dsApp.FillFromJSON()

	// Get packages
	packages, err := datastore.GetPackagesAppIsIn(dsApp.AppID)
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := appTemplate{}
	template.App = dsApp
	template.Packages = packages

	returnTemplate(w, "app", template)
}

type appTemplate struct {
	GlobalTemplate
	App      datastore.DsApp
	Packages []datastore.DsPackage
}
