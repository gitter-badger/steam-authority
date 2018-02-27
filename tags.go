package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/mysql"
)

func tagsHandler(w http.ResponseWriter, r *http.Request) {

	tags, err := mysql.GetAllTags()
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 500, "Error getting tags")
		return
	}

	// Template
	template := tagsTemplate{}
	template.Fill(r)
	template.Tags = tags

	returnTemplate(w, r, "tags", template)
	return
}

type tagsTemplate struct {
	GlobalTemplate
	Tags []mysql.Tag
}
