package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	template := homeTemplate{}
	template.test = "xx"

	sendWebsocket(template)

	returnTemplate(w, "home", template)
}

type homeTemplate struct {
	test string `json:"test"`
}
