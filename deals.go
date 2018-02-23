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

	search := url.Values{}
	search.Set("is_free", "1")
	search.Set("name", "-")

	// Types not in this list will show first
	apps, err := mysql.SearchApps(search, 1000, "FIELD(`type`,'game','dlc','demo','mod','video','movie','series','episode','application','tool','advertising'), name ASC")
	if err != nil {
		logger.Error(err)
	}

	template := dealsTemplate{}
	template.SetSession(r)
	template.Apps = apps
	template.Tab = tab

	returnTemplate(w, "deals", template)
	return
}

type dealsTemplate struct {
	GlobalTemplate
	Apps []mysql.App
	Tab  string
}
