package main

import "net/http"

func creditsHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "credits", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "contact", nil)
}

func donateHandler(w http.ResponseWriter, r *http.Request) {
	// todo, setup Patreon for donations
	returnTemplate(w, "donate", nil)
}

func faqsHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "faqs", nil)
}
