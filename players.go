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

const (
	playersToRank = 500
)

func playersHandler(w http.ResponseWriter, r *http.Request) {

	// Normalise the order
	var ranks []datastore.DsRank
	var err error

	switch chi.URLParam(r, "id") {
	case "level":
		ranks, err = datastore.GetRanksBy("level_rank")

		for k := range ranks {
			ranks[k].Rank = ranks[k].LevelRank
		}
		//case "games":
		//	sort = "games_rank"
		//case "badges":
		//	sort = "badges_rank"
		//case "playtime":
		//	sort = "play_time_rank"
		//case "steamtime":
		//	sort = "-time_created_rank"
	default:
		ranks, err = datastore.GetRanksBy("level_rank")

		for k := range ranks {
			ranks[k].Rank = ranks[k].LevelRank
		}
	}

	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	template := playersTemplate{}
	template.Ranks = ranks

	returnTemplate(w, "players", template)
}

type playersTemplate struct {
	GlobalTemplate
	Ranks []datastore.DsRank
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
			dsPlayer.TimeUpdated = time.Now().Unix()
			dsPlayer.CountryCode = summary.Response.Players[0].LOCCountryCode
			dsPlayer.StateCode = summary.Response.Players[0].LOCStateCode
			dsPlayer.PersonaName = summary.Response.Players[0].PersonaName

			// todo, get friends, player bans, groups

			// todo, clear latest players cache
			datastore.SaveKind(dsPlayer.GetKey(), dsPlayer)
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

	id, err := steam.GetID(r.PostFormValue("id"))
	if err != nil {
		logger.Error(err)
		// todo error page
	}

	http.Redirect(w, r, "/players/"+id, 302) // Temp redirect
}

type playerTemplate struct {
	GlobalTemplate
	Player datastore.DsPlayer
}

func reRankHandler(w http.ResponseWriter, r *http.Request) {

	reRank("level")
	//reRank("games")
	//reRank("badhes")

	w.Write([]byte("OK"))
}

func reRank(order string) {

	// todo!!
	// - key new ranks table by ID64 - done
	// - get keys only for current ranked players
	// - read players table and get top 500
	// - insert into ranks table
	// - delete any keys that have not just been updated
	// todo!!

	//keys, err := datastore.GetRankKeys()

	//players, err := datastore.GetPlayers("-level", playersToRank)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//
	//var bulk []*datastore.DsPlayer
	//
	//for index := 0; index < len(players); index++ {
	//	player := players[index]
	//	//player.LevelRank = index + 1
	//	bulk = append(bulk, &player)
	//}

	//err = datastore.BulkSavePlayers(bulk)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
}
