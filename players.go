package main

import (
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/Jleagle/go-helpers/logger"

	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/steam"
)

func playersHandler(w http.ResponseWriter, r *http.Request) {

	template := playersTemplate{}

	returnTemplate(w, "players", template)
}

func playerHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	dsPlayer, err := datastore.GetPlayer(id)
	if err != nil {
		if err.Error() == "datastore: no such entity" || err.Error() == "expired" {

			//Get summary
			summary, _ := steam.GetPlayerSummaries([]int{idx})

			id64, _ := strconv.Atoi(summary.Response.Players[0].SteamID)

			dsPlayer.ID64 = id64
			dsPlayer.Avatar = summary.Response.Players[0].AvatarFull
			dsPlayer.ValintyURL = path.Base(summary.Response.Players[0].ProfileURL)
			dsPlayer.RealName = summary.Response.Players[0].RealName
			dsPlayer.LastUpdated = time.Now().Unix()
			dsPlayer.CountryCode = summary.Response.Players[0].LOCCountryCode
			dsPlayer.StateCode = summary.Response.Players[0].LOCStateCode

			// todo, get friends, player bans, groups

			datastore.SavePlayer(dsPlayer)
		} else {
			logger.Error(err)
			returnErrorTemplate(w, 404, err.Error())
			return
		}
	}

	template := playerTemplate{}
	template.Player = dsPlayer

	returnTemplate(w, "player", template)
}

func playerIDHandler(w http.ResponseWriter, r *http.Request) {

	// api := os.Getenv("STEAM_API_KEY")
	r.ParseForm()

	http.Redirect(w, r, "/players/76561197968626192", 302) // Temp redirect
}

type playersTemplate struct {
	GlobalTemplate
	Players []datastore.DsPlayer
}

type playerTemplate struct {
	GlobalTemplate
	Player datastore.DsPlayer
}
