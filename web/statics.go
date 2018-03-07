package web

import "net/http"

func CreditsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.Fill(r)

	returnTemplate(w, r, "credits", template)
}

func DonateHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.Fill(r)

	returnTemplate(w, r, "donate", template)
}

func FAQsHandler(w http.ResponseWriter, r *http.Request) {

	template := staticTemplate{}
	template.Fill(r)

	returnTemplate(w, r, "faqs", template)
}

func Error404Handler(w http.ResponseWriter, r *http.Request) {

	returnErrorTemplate(w, r, 404, "Page not found")
}

type staticTemplate struct {
	GlobalTemplate
}
