package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/mysql"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	packages, err := mysql.GetLatestPackages()
	if err != nil {
		logger.Error(err)
	}

	template := packagesTemplate{}
	template.SetSession(r)
	template.Packages = packages

	returnTemplate(w, r, "packages", template)
}

type packagesTemplate struct {
	GlobalTemplate
	Packages []mysql.Package
}

func packageHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)
	}

	pack, err := mysql.GetPackage(id)
	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			returnErrorTemplate(w, r, 404, "We can't find this package in our database, there may not be one with this ID.")
			return
		}

		logger.Error(err)

		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	appIDs, err := pack.GetApps()
	if err != nil {
		logger.Error(err)
	}

	apps, err := mysql.GetApps(appIDs, []string{})
	if err != nil {
		logger.Error(err)
	}

	template := packageTemplate{}
	template.SetSession(r)
	template.Package = pack
	template.Apps = apps

	returnTemplate(w, r, "package", template)
}

type packageTemplate struct {
	GlobalTemplate
	Package mysql.Package
	Apps    []mysql.App
}
