package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
)

func changesHandler(w http.ResponseWriter, r *http.Request) {

	template := changesTemplate{}

	// Get changes
	changes, err := datastore.GetLatestChanges(100)
	if err != nil {
		logger.Error(err)
	}
	template.Changes = changes

	// // Get apps/packages
	// apps := make([]int, 0, 500)
	// packages := make([]int, 0, 500)
	// for _, v := range changes {
	// 	apps = append(apps, v.Apps...)
	// 	packages = append(apps, v.Packages...)
	// }

	// dsApps, err := datastore.GetMultiAppsByKey(apps)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// template.Apps = make(map[int]datastore.DsApp)
	// for _, v := range dsApps {
	// 	template.Apps[v.AppID] = v
	// }

	// dsPackages, err := datastore.GetMultiPackagesByKey(packages)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// template.Packages = make(map[int]datastore.DsPackage)
	// for _, v := range dsPackages {
	// 	template.Packages[v.PackageID] = v
	// }

	// pretty.Print(template)

	returnTemplate(w, "changes", template)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {

	change, err := datastore.GetChange(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)
		if err.Error() == "datastore: no such entity" {
			returnErrorTemplate(w, 404, "We can't find this change in our database, there may not be one with this ID.")
			return
		}
	}

	template := changeTemplate{}
	template.Change = change

	returnTemplate(w, "change", template)
}

type changesTemplate struct {
	GlobalTemplate
	Changes  []datastore.Change
	Apps     map[int]mysql.App
	Packages map[int]mysql.Package
}

type changeTemplate struct {
	GlobalTemplate
	Change *datastore.Change
}
