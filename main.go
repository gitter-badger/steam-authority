package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {

	// todo, give it the path in code not env
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("STEAM_GOOGLE_APPLICATION_CREDENTIALS"))

	arguments := os.Args[1:]

	if len(arguments) > 0 {

		switch arguments[0] {
		case "check-for-changes":
			fmt.Println("Checking for changes")
			checkForChanges()
		default:
			fmt.Println("No such CLI command")
		}

		os.Exit(0)
	}

	r := chi.NewRouter()

	r.Get("/", homeHandler)

	r.Get("/apps", appsHandler)
	r.Get("/apps/mine", appsHandler)

	r.Get("/packages", packagesHandler)

	r.Get("/changes", changesHandler)
	r.Get("/changes/{id}", changeHandler)

	http.ListenAndServe(":8085", r)
}
