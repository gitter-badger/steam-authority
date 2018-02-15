package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
)

func newsHandler(w http.ResponseWriter, r *http.Request) {

	articles, err := datastore.GetArticles(0, 100)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 500, "Error getting articles")
		return
	}

	// Filter out artciles with no app id
	var filteredArticles []datastore.Article
	var appIDs []int

	for _, v := range articles {
		if v.AppID != 0 {
			filteredArticles = append(filteredArticles, v)
			appIDs = append(appIDs, v.AppID)
		}
	}

	// Get app info
	apps, err := mysql.GetApps(appIDs)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 500, "Error getting articles")
		return
	}

	// Make map of apps
	appMap := make(map[int]mysql.App)
	for _, v := range apps {
		appMap[v.ID] = v
	}

	fmt.Println(appMap[30].Name)

	// Template
	template := articlesTemplate{}
	template.Articles = filteredArticles
	template.Apps = appMap

	returnTemplate(w, "news", template)
	return
}

type articlesTemplate struct {
	Articles []datastore.Article
	Apps     map[int]mysql.App
}
