package main

import "net/http"

func creditsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, r, "credits", template)
}

func donateHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, r, "donate", template)
}

func faqsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, r, "faqs", template)
}

type staticTemplate struct {
	GlobalTemplate
}
