package main

import "net/http"

func experienceHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "experience", nil)
}
