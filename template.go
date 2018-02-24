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
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/steam-authority/steam-authority/session"
)

func returnTemplate(w http.ResponseWriter, r *http.Request, page string, pageData interface{}) (err error) {

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
		returnErrorTemplate(w, r, 404, err.Error())
		return err
	}

	// Write a respone
	buf := &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, page, pageData)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 500, "Something has gone wrong, the error has been logged!")
		return
	} else {
		// No error, send the content, HTTP 200 response status implied
		buf.WriteTo(w)
		return
	}

	return nil
}

func returnErrorTemplate(w http.ResponseWriter, r *http.Request, code int, message string) {

	w.WriteHeader(code)

	tmpl := errorTemplate{}
	tmpl.SetSession(r)
	tmpl.Code = code
	tmpl.Message = message

	returnTemplate(w, r, "error", tmpl)
}

func getTemplateFuncMap() map[string]interface{} {
	return template.FuncMap{
		"join":  func(a []string) string { return strings.Join(a, ", ") },
		"title": func(a string) string { return strings.Title(a) },
		"comma": func(a int) string { return humanize.Comma(int64(a)) },
		"slug":  func(a string) string { return slug.Make(a) },
		"apps": func(a []int, appsMap map[int]mysql.App) template.HTML {
			var apps []string
			for _, v := range a {
				apps = append(apps, "<a href=\"/apps/"+strconv.Itoa(v)+"\">"+appsMap[v].GetName()+"</a>")
			}
			return template.HTML("Apps: " + strings.Join(apps, ", "))
		},
		"packages": func(a []int, packagesMap map[int]mysql.Package) template.HTML {
			var packages []string
			for _, v := range a {
				packages = append(packages, "<a href=\"/packages/"+strconv.Itoa(v)+"\">"+packagesMap[v].GetName()+"</a>")
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
	Env    string
	ID     int
	Name   string
	Avatar string
}

func (t *GlobalTemplate) SetSession(r *http.Request) {

	id, _ := session.Read(r, session.ID)

	t.ID, _ = strconv.Atoi(id)
	t.Name, _ = session.Read(r, session.Name)
	t.Avatar, _ = session.Read(r, session.Avatar)
}

func (t GlobalTemplate) LoggedIn() bool {
	return t.ID > 0
}
