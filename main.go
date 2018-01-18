package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/pics"
	"github.com/steam-authority/steam-authority/websockets"
)

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("STEAM_GOOGLE_APPLICATION_CREDENTIALS"))

	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/apps", appsHandler)
	r.Get("/apps/{id}", appsHandler)
	r.Get("/apps/mine", appsHandler)
	r.Get("/packages", packagesHandler)
	r.Get("/packages/{id}", packagesHandler)
	r.Get("/changes", changesHandler)
	r.Get("/changes/{id}", changeHandler)
	r.Get("/websocket", websockets.Handler)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	go pics.Run()

	http.ListenAndServe(":8085", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		logger.Info("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
