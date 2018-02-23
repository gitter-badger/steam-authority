package main

import "net/http"

func creditsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, "credits", template)
}

func donateHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, "donate", template)
}

func faqsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.SetSession(r)

	returnTemplate(w, "faqs", template)
}

type staticTemplate struct {
	GlobalTemplate
}
