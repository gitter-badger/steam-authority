package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"

	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/{url}/list", homeRoute)

	http.ListenAndServe(":8085", r)
}

func homeRoute(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("http://localhost:8086/info?apps=440,441,730&packages=75330&prettyprint=1")
	if err != nil {
		// handle err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// Handle err
	}

	fmt.Printf("%s\n", string(contents))

	returnTemplate(w, "home", nil)
}

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
