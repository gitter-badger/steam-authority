package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/mysql"
)

func genresHandler(w http.ResponseWriter, r *http.Request) {

	genres, err := mysql.GetAllGenres()
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 500, "Error getting genres")
		return
	}

	// Template
	template := genresTemplate{}
	template.SetSession(r)
	template.Genres = genres

	returnTemplate(w, r, "genres", template)
	return
}

type genresTemplate struct {
	GlobalTemplate
	Genres []mysql.Genre
}
