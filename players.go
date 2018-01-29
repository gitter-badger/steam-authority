package main

import (
	"net/http"
	"strconv"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/steam"
	slugify "github.com/gosimple/slug"
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

			dsPlayer.ID64 = idx
			
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
	Player datastore.DsPlayer
}

func reRankHandler(w http.ResponseWriter, r *http.Request) {

	var playersToRank = 500

	// Get keys, will delete any that are not removed from this map
	oldKeys, err := datastore.GetRankKeys()

	newRanks := make(map[int]*datastore.DsRank)

	// Get players by level
	players, err := datastore.GetPlayers("-level", playersToRank)
	if err != nil {
		logger.Error(err)
		return
	}

	for k, v := range players {

		_, ok := newRanks[v.ID64]
		if !ok {

			rank := &datastore.DsRank{}
			rank.FillFromPlayer(v)

			newRanks[v.ID64] = rank
		}
		newRanks[v.ID64].LevelRank = k + 1

		_, ok = oldKeys[strconv.Itoa(v.ID64)]
		if ok {
			delete(oldKeys, strconv.Itoa(v.ID64))
		}
	}

	// Convert new ranks to slice
	var ranks []*datastore.DsRank
	for _, v := range newRanks {
		ranks = append(ranks, v)
	}

	// Bulk save ranks
	err = datastore.BulkSaveRanks(ranks)
	if err != nil {
		logger.Error(err)
		return
	}

	// Delete leftover keys
	datastore.BulkDeleteRanks(oldKeys)

	w.Write([]byte("OK"))
}
