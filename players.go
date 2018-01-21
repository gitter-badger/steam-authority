package main

import (
	"net/http"
)

func playersHandler(w http.ResponseWriter, r *http.Request) {

	template := playersTemplate{}

	returnTemplate(w, "players", template)
}

func playerHandler(w http.ResponseWriter, r *http.Request) {

	template := playerTemplate{}

	returnTemplate(w, "player", template)
}

type playersTemplate struct {
	GlobalTemplate
}

type playerTemplate struct {
	GlobalTemplate
}
