package web

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/players", 302)
	return

	template := homeTemplate{}
	returnTemplate(w, r, "home", template)
}

type homeTemplate struct {
}
