package main

import (
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"
)

func changesHandler(w http.ResponseWriter, r *http.Request) {

	template := changesTemplate{}

	// Get changes
	client, context := getDSClient()
	q := datastore.NewQuery("Change").Order("-change_id")
	it := client.Run(context, q)
	for {
		var change dsChange
		_, err := it.Next(&change)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		template.Changes = append(template.Changes, change)
	}

	returnTemplate(w, "changes", template)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {

	client, context := getDSClient()

	key := datastore.NameKey("Change", chi.URLParam(r, "id"), nil)

	change := &dsChange{}
	if err := client.Get(context, key, change); err != nil {
		logger.Error(err)
	}

	template := changeTemplate{}
	template.Change = change

	returnTemplate(w, "change", template)
}

type changesTemplate struct {
	Changes []dsChange
}

type changeTemplate struct {
	Change *dsChange
}
