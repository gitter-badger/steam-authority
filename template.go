package main

import (
	"html/template"
	"net/http"
	"path"
	"runtime"
)

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) {

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		// Handle err
	}
	folder := path.Dir(file)

	// Load templates needed
	t, err := template.ParseFiles(folder+"/templates/header.html", folder+"/templates/footer.html", folder+"/templates/"+page+".html")
	if err != nil {
		// Handle err
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		// Handle err
	}
}
