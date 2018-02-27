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

type staticTemplate struct {
	GlobalTemplate
}
