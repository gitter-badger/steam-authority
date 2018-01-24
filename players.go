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

	// Normalise the order
	var sort string
	switch chi.URLParam(r, "id") {
	case "level":
		sort = "level_rank"
	case "games":
		sort = "games_rank"
	case "badges":
		sort = "badges_rank"
	case "playtime":
		sort = "play_time_rank"
	case "steamtime":
		sort = "-time_created_rank"
	default:
		sort = "level_rank"
	}

	players, err := datastore.GetPlayers(sort, 10)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	// Set rank field
	for k := range players {
		players[k].Rank = players[k].LevelRank // todo, set depending on the sort
	}

	template := playersTemplate{}
	template.Players = players

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
			dsPlayer.PersonaName = summary.Response.Players[0].PersonaName

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
