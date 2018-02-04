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

	template.Packages = packages

	returnTemplate(w, "packages", template)
}

func packageHandler(w http.ResponseWriter, r *http.Request) {

	idx, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		logger.Error(err)
	}

	packagex, err := mysql.GetPackage(uint(idx))
	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			returnErrorTemplate(w, 404, "We can't find this package in our database, there may not be one with this ID.")
			return
		}

		logger.Error(err)

		returnErrorTemplate(w, 404, err.Error())
		return
	}

	appIDs, err := packagex.GetApps()
	if err != nil {
		logger.Error(err)
	}

	apps, err := mysql.GetApps(appIDs)
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
	Packages []mysql.Package
}

type packageTemplate struct {
	GlobalTemplate
	Package mysql.Package
	Apps    []mysql.App
}
