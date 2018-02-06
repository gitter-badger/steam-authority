package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	slugify "github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/steam"
)

func playersHandler(w http.ResponseWriter, r *http.Request) {

	// Normalise the order
	var ranks []datastore.Rank
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
	Ranks []datastore.Rank
}

func playerHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	slug := chi.URLParam(r, "slug")

	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 404, err.Error())
		return
	}

	dsPlayer, err := datastore.GetPlayer(id)
	if err != nil {
		if err.Error() == "datastore: no such entity" || err.Error() == "expired" {

			dsPlayer.PlayerID = idx

			//Get summary
			summary, err := steam.GetPlayerSummaries([]int{idx})
			if err != nil {
				logger.Error(err)
				returnErrorTemplate(w, 404, err.Error())
				return
			}
			dsPlayer.FillFromSummary(summary)

			//Get friends
			friends, err := steam.GetFriendList(id)
			if err != nil {
				logger.Error(err)
				returnErrorTemplate(w, 404, err.Error())
				return
			}
			dsPlayer.FillFromFriends(friends)

			// todo, get player bans, groups
			// todo, clear latest players cache
			dsPlayer.Tidy()
			_, err = datastore.SaveKind(dsPlayer.GetKey(), &dsPlayer)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Error(err)
			returnErrorTemplate(w, 404, err.Error())
			return
		}
	}

	// Redirect to correct slug
	correctSLug := slugify.Make(dsPlayer.PersonaName)
	if slug != "" && slug != correctSLug {
		http.Redirect(w, r, "/players/"+id+"/"+correctSLug, 302)
		return
	}

	// Template
	template := playerTemplate{}
	template.Player = dsPlayer

	returnTemplate(w, "player", template)
}

func playerIDHandler(w http.ResponseWriter, r *http.Request) {

	post := r.PostFormValue("id")

	id, err := steam.GetID(post)
	if err != nil {
		logger.Info(err.Error() + ": " + post)
		returnErrorTemplate(w, 404, "Can't find user: "+post)
		return
	}

	http.Redirect(w, r, "/players/"+id, 302)
}

type playerTemplate struct {
	GlobalTemplate
	Player datastore.Player
}
