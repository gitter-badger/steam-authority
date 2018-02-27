package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	slugify "github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/queue"
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
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	template := playersTemplate{}
	template.Fill(r)
	template.Ranks = ranks

	returnTemplate(w, r, "players", template)
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
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	queue.PlayerProducer(idx)

	player, err := datastore.GetPlayer(idx, true)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	// Redirect to correct slug
	correctSLug := slugify.Make(player.PersonaName)
	if slug != "" && slug != correctSLug {
		http.Redirect(w, r, "/players/"+id+"/"+correctSLug, 302)
		return
	}

	// Make friend ID slice
	var friendsSlice []int
	for _, v := range player.Friends {
		s, _ := strconv.Atoi(v.SteamID)
		friendsSlice = append(friendsSlice, s)
	}

	// Get friends for template
	friends, err := datastore.GetPlayersByIDs(friendsSlice)
	if err != nil {
		logger.Error(err)
	}

	// Add friends to rabbit
	if player.AddFriends {
		for _, v := range friendsSlice {
			queue.PlayerProducer(v)
		}
	}

	// Template
	template := playerTemplate{}
	template.Fill(r)
	template.Player = player
	template.Friends = friends

	returnTemplate(w, r, "player", template)
}

type playerTemplate struct {
	GlobalTemplate
	Player  datastore.Player
	Friends []datastore.Player
}

func playerIDHandler(w http.ResponseWriter, r *http.Request) {

	post := r.PostFormValue("id")

	// todo, check DB before doing api call

	id, err := steam.GetID(post)
	if err != nil {
		logger.Info(err.Error() + ": " + post)
		returnErrorTemplate(w, r, 404, "Can't find user: "+post)
		return
	}

	http.Redirect(w, r, "/players/"+id, 302)
}
