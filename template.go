package main

import (
	"html/template"
	"net/http"
	"path"
	"runtime"

	"github.com/Jleagle/go-helpers/logger"
)

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) (err error) {

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		logger.Info("Failed to get path")
	}
	folder := path.Dir(file)

	// Load templates needed
	t, err := template.ParseFiles(folder+"/templates/header.html", folder+"/templates/footer.html", folder+"/templates/"+page+".html")
	if err != nil {
		logger.Info("x")
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return err
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		logger.Info("y")
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return err
	}

	return nil
}

func returnErrorTemplate(w http.ResponseWriter, code int, message string) {

	template := errorTemplate{
		Code:    code,
		Message: message,
	}

	returnTemplate(w, "error", template)
}

type errorTemplate struct {
	GlobalTemplate
	Code    int
	Message string
}

// GlobalTemplate is added to every other template
type GlobalTemplate struct {
	Env string
}
