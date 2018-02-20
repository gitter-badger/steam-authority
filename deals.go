package main

import (
	"net/http"
	"net/url"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/mysql"
)

const (
	FREE      = "free"
	CHANGES   = "changes"
	DISCOUNTS = "discounts"
)

func dealsHandler(w http.ResponseWriter, r *http.Request) {

	tab := chi.URLParam(r, "id")
	if tab == "" {
		tab = "free"
	}

	template := dealsTemplate{}

	search := url.Values{}
	search.Set("is_free", "1")
	search.Set("name", "-")

	apps, err := mysql.SearchApps(search, 1000, "name ASC")
	if err != nil {
		logger.Error(err)
	}

	template.Apps = apps
	template.Tab = tab

	returnTemplate(w, "deals", template)
	return
}

type dealsTemplate struct {
	Apps []mysql.App
	Tab  string
}
