package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	slugify "github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/steam-authority/steam-authority/queue"
)

func appsHandler(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := mysql.SearchApps(r.URL.Query(), 96, "id DESC")
	if err != nil {
		logger.Error(err)
	}

	// Get apps count
	count, err := mysql.CountApps()
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := appsTemplate{}
	template.SetSession(r)
	template.Apps = apps
	template.Count = count

	returnTemplate(w, r, "apps", template)
}

type appsTemplate struct {
	GlobalTemplate
	Apps  []mysql.App
	Count int
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	slug := chi.URLParam(r, "slug")

	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	queue.AppProducer(idx, 0)

	// Get app
	app, err := mysql.GetApp(idx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			returnErrorTemplate(w, r, 404, err.Error())
		} else {
			logger.Error(err)
			returnErrorTemplate(w, r, 500, err.Error())
		}
		return
	}

	// Make banners
	banners := make(map[string][]string)
	var primary []string

	if app.ReleaseState == "prerelease" {
		primary = append(primary, "This game is not released yet!")
	}
	if app.Type == "movie" {
		primary = append(primary, "This listing is for a movie")
	}
	banners["primary"] = primary

	//if err != nil {
	//	if err.Error() == "sql: no rows in result set" {
	//
	//		// Create the app
	//		app, err = mysql.CreateApp(idx)
	//		if err != nil {
	//			logger.Error(err)
	//			returnErrorTemplate(w, 404, err.Error())
	//			return
	//		}
	//
	//		// Get app articles
	//		_, err = datastore.GetArticlesFromSteam(idx)
	//		if err != nil {
	//			logger.Error(err)
	//		}
	//
	//	} else {
	//		logger.Error(err)
	//		returnErrorTemplate(w, 500, err.Error())
	//		return
	//	}
	//}

	// Get news
	news, err := datastore.GetArticles(idx, 1000)

	// Redirect to correct slug
	correctSLug := slugify.Make(app.Name)
	if slug != "" && app.Name != "" && slug != correctSLug {
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
	template.SetSession(r)
	template.App = app
	template.Packages = packages
	template.Articles = news
	template.Banners = banners

	returnTemplate(w, r, "app", template)
}

type appTemplate struct {
	GlobalTemplate
	App      mysql.App
	Packages []mysql.Package
	Articles []datastore.Article
	Banners  map[string][]string
}
