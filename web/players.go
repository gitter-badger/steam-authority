package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	slugify "github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/queue"
	"github.com/steam-authority/steam-authority/steam"
)

func PlayersHandler(w http.ResponseWriter, r *http.Request) {

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

	// Count players
	playersCount, err := datastore.CountPlayers()
	if err != nil {
		logger.Error(err)
	}

	// Count ranks
	ranksCount, err := datastore.CountRankedPlayers()
	if err != nil {
		logger.Error(err)
	}

	template := playersTemplate{}
	template.Fill(r)
	template.Ranks = ranks
	template.PlayersCount = playersCount
	template.RanksCount = ranksCount

	returnTemplate(w, r, "players", template)
	return
}

type playersTemplate struct {
	GlobalTemplate
	Ranks        []datastore.Rank
	PlayersCount int
	RanksCount   int
}

func PlayerHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	slug := chi.URLParam(r, "slug")

	idx, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	//queue.PlayerProducer(76561197995497914)

	player, err := datastore.GetPlayer(idx)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	err = player.UpdateIfNeeded()
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	// Queue friends
	if player.ShouldUpdateFriends() {

		for _, v := range player.Friends {
			vv, _ := strconv.Atoi(v.SteamID)
			queue.PlayerProducer(vv)
		}

		player.FriendsAddedAt = time.Now()
		player.Save()
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

	// Template
	template := playerTemplate{}
	template.Fill(r)
	template.Player = player
	template.Friends = friends

	returnTemplate(w, r, "player", template)
}

type playerTemplate struct {
	GlobalTemplate
	Player  *datastore.Player
	Friends []*datastore.Player
}

func PlayerIDHandler(w http.ResponseWriter, r *http.Request) {

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
