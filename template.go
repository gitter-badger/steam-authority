package main

import (
	"bytes"
	"html/template"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	buf := &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, page, pageData)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 500, "Something has gone wrong, the error has been logged!")
		return
	} else {
		// No error, send the content, HTTP 200 response status implied
		buf.WriteTo(w)
		return
	}

	return nil
}

func returnErrorTemplate(w http.ResponseWriter, code int, message string) {

	w.WriteHeader(code)

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
		"apps": func(a []int) template.HTML {
			var apps []string
			for _, v := range a {
				apps = append(apps, "<a href=\"/apps/"+strconv.Itoa(v)+"\">"+strconv.Itoa(v)+"</a>")
			}
			return template.HTML("Apps: " + strings.Join(apps, ", "))
		},
		"packages": func(a []int) template.HTML {
			var packages []string
			for _, v := range a {
				packages = append(packages, "<a href=\"/packages/"+strconv.Itoa(v)+"\">"+strconv.Itoa(v)+"</a>")
			}
			return template.HTML("Packages: " + strings.Join(packages, ", "))
		},
		"unix": func(t time.Time) int64 {
			return t.Unix()
		},
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
