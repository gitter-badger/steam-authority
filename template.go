package main

import (
	"html/template"
	"net/http"
	"path"
	"runtime"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/dustin/go-humanize"
	"github.com/gosimple/slug"
)

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) (err error) {

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		logger.Info("Failed to get path")
	}
	folder := path.Dir(file)

	// Load templates needed
	t, err := template.New("t").Funcs(getTemplateFuncMap()).ParseFiles(folder+"/templates/header.html", folder+"/templates/footer.html", folder+"/templates/"+page+".html")
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return err
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return err
	}

	return nil
}

func returnErrorTemplate(w http.ResponseWriter, code int, message string) {

	tmpl := errorTemplate{
		Code:    code,
		Message: message,
	}

	returnTemplate(w, "error", tmpl)
}

func getTemplateFuncMap() map[string]interface{} {
	return template.FuncMap{
		"join":  func(a []string) string { return strings.Join(a, ", ") },
		"title": func(a string) string { return strings.Title(a) },
		"comma": func(a int) string { return humanize.Comma(int64(a)) },
		"slug":  func(a string) string { return slug.Make(a) },
	}
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
