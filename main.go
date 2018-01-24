package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/websockets"
)

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("STEAM_GOOGLE_APPLICATION_CREDENTIALS"))

	logger.SetRollbarKey(os.Getenv("STEAM_ROLLBAR_PRIVATE"))

	r := chi.NewRouter()
	r.Get("/", homeHandler)

	// Apps
	r.Get("/apps", appsHandler)
	r.Get("/apps/{id}", appHandler)
	r.Get("/apps/{id}/{slug}", appHandler)

	// Packages
	r.Get("/packages", packagesHandler)
	r.Get("/packages/{id}", packageHandler)

	// Players
	r.Get("/players", playersHandler)
	r.Post("/players", playerIDHandler)
	r.Get("/players/{id:[a-z]+}", playersHandler)
	r.Get("/players/{id:[0-9]+}", playerHandler)
	r.Get("/players/{id:[0-9]+}/{slug}", playerHandler)

	// Changes
	r.Get("/changes", changesHandler)
	r.Get("/changes/{id}", changeHandler)

	// Experience
	r.Get("/experience", experienceHandler)
	r.Get("/experience/{id}", experienceHandler)

	// Static pages
	r.Get("/contact", contactHandler)
	r.Get("/donate", donateHandler)
	r.Get("/faqs", faqsHandler)
	r.Get("/credits", creditsHandler)

	// Other
	r.Get("/websocket", websockets.Handler)
	r.Get("/changelog", changelogHandler)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	// go pics.Run()

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
