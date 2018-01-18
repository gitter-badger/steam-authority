package main

import (
	"net/http"

	"github.com/steam-authority/steam-authority/websockets"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	template := HomeTemplate{}
	template.Test = "xx"

	websockets.Send(template)

	returnTemplate(w, "home", template)
}

type HomeTemplate struct {
	Test string `json:"test"`
}
