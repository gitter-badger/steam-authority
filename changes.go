package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
)

func changesHandler(w http.ResponseWriter, r *http.Request) {

	template := changesTemplate{}

	changes, err := datastore.GetLatestChanges(10)
	if err != nil {
		logger.Error(err)
	}

	template.Changes = changes

	returnTemplate(w, "changes", template)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {

	change, err := datastore.GetChange(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err)
		if err.Error() == "datastore: no such entity" {
			returnErrorTemplate(w, 404, "We can't find this change in our database, there may not have been one with this ID.")
			return
		}
	}

	template := changeTemplate{}
	template.Change = change

	returnTemplate(w, "change", template)
}

type changesTemplate struct {
	Changes []datastore.DsChange
}

type changeTemplate struct {
	Change *datastore.DsChange
}
