package main

import (
	"net/http"
	"net/url"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/mysql"
)

func dealsHandler(w http.ResponseWriter, r *http.Request) {

	tab := chi.URLParam(r, "id")
	if tab == "" {
		tab = "free"
	}

	search := url.Values{}

	apps, err := mysql.SearchApps(search)
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := dealsTemplate{}
	template.Apps = apps
	template.Tab = tab

	returnTemplate(w, "deals", template)
	return
}

type dealsTemplate struct {
	Apps []mysql.App
	Tab  string
}
