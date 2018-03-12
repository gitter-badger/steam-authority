package web

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	template := homeTemplate{}
	template.Fill(r)

	returnTemplate(w, r, "home", template)
}

type homeTemplate struct {
	GlobalTemplate
}
