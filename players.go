package main

import (
	"net/http"
	"os"
	"path"

	"github.com/Acidic9/steam"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
	"github.com/steam-authority/steam-authority/datastore"
)

func playersHandler(w http.ResponseWriter, r *http.Request) {

	template := playersTemplate{}

	returnTemplate(w, "players", template)
}

func playerHandler(w http.ResponseWriter, r *http.Request) {

	template := playerTemplate{}

	returnTemplate(w, "player", template)
}

func playerIDHandler(w http.ResponseWriter, r *http.Request) {

	api := os.Getenv("STEAM_API_KEY")
	r.ParseForm()
	steam64 := steam.SearchForID(r.Form.Get("id"), api)
	pretty.Print(steam64)
	summary, err := steam.GetPlayerSummaries(api, steam64)
	if err != nil {
		logger.Error(err)
	}

	http.Redirect(w, r, "/players/"+path.Base(summary.ProfileURL), 302) // Temp redirect
}

type playersTemplate struct {
	GlobalTemplate
	Players []datastore.DsPlayer
}

type playerTemplate struct {
	GlobalTemplate
}
