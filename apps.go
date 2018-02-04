package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	slugify "github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/mysql"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := mysql.SearchApps(r.URL.Query())
	if err != nil {
		logger.Error(err)
	}

	// Get apps count
	count, err := mysql.CountTable("apps")
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
	Apps  []mysql.App
	Count uint
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	slug := chi.URLParam(r, "slug")

	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	// Get app
	app, err := mysql.GetApp(uint(idx))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {

			// Get app details
			app.FillFromAppDetails()
			if err != nil {
				if err.Error() == "no app with id" {
					returnErrorTemplate(w, 404, "Sorry but there is no app with this ID")
					return
				}
				logger.Error(err)
			}

			err = app.Save(idx)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Error(err)
			returnErrorTemplate(w, 404, err.Error())
			return
		}
	}

	// Redirect to correct slug
	correctSLug := slugify.Make(app.Name)
	if slug != "" && slug != correctSLug {
		http.Redirect(w, r, "/apps/"+id+"/"+correctSLug, 302)
		return
	}

	// Get packages
	packages, err := mysql.GetPackagesAppIsIn(app.ID)
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
	App      mysql.App
	Packages []mysql.Package
}
