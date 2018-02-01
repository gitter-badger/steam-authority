package main

import (
	"net/http"
	"os"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/jmoiron/sqlx"
	"github.com/steam-authority/steam-authority/mysql"
)

func tagsHandler(w http.ResponseWriter, r *http.Request) {

	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("mysql", os.Getenv("STEAM_SQL_DSN"))
	if err != nil {
		logger.Error(err)
	}

	var tags []mysql.Tag
	err = db.Select(&tags, "SELECT * FROM tags ORDER BY games DESC")
	if err != nil {
		logger.Error(err)
	}

	// Template
	template := tagsTemplate{}
	template.Tags = tags
	returnTemplate(w, "tags", template)
	return
}

type tagsTemplate struct {
	Tags []mysql.Tag
}
