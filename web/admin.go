package web

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/steam-authority/steam-authority/queue"
	"github.com/steam-authority/steam-authority/steam"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {

	option := chi.URLParam(r, "option")

	switch option {
	case "apps": // Add all apps to queue
		go adminApps(w, r)
	case "deploy":
		go adminDeploy(w, r)
	case "donations":
		go adminDonations(w, r)
	case "genres":
		go adminGenres(w, r)
	case "ranks":
		go adminRanks(w, r)
	case "tags":
		go adminTags(w, r)
	}

	if option != "" {
		http.Redirect(w, r, "/admin?"+option, 302)
		return
	}

	// Template
	template := adminTemplate{}
	template.Fill(r)

	returnTemplate(w, r, "admin", template)
	return
}

type adminTemplate struct {
	GlobalTemplate
}

func adminApps(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := steam.GetAppList()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, v := range apps {
		queue.AppProducer(v.AppID, 0)
	}

	logger.Info(strconv.Itoa(len(apps)) + " apps added to rabbit")
}

func adminDeploy(w http.ResponseWriter, r *http.Request) {

}

func adminDonations(w http.ResponseWriter, r *http.Request) {

}

// todo, handle genres that no longer have any games.
func adminGenres(w http.ResponseWriter, r *http.Request) {

	filter := url.Values{}
	filter.Set("json_depth", "3")

	apps, err := mysql.SearchApps(filter, 0, "")
	if err != nil {
		logger.Error(err)
	}

	counts := make(map[int]*adminGenreCount)

	for _, app := range apps {
		genres, err := app.GetGenres()
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, genre := range genres {
			//logger.Info(genre.Description)

			if _, ok := counts[genre.ID]; ok {
				counts[genre.ID].Count++
			} else {
				counts[genre.ID] = &adminGenreCount{
					Count: 1,
					Genre: genre,
				}
			}
		}
	}

	for _, v := range counts {
		err := mysql.SaveOrUpdateGenre(v.Genre.ID, v.Genre.Description, v.Count)
		if err != nil {
			logger.Error(err)
		}
	}

	logger.Info("Genres updated")
}

type adminGenreCount struct {
	Count int
	Genre steam.AppDetailsGenre
}

// todo, handle tags that no longer have any games.
func adminTags(w http.ResponseWriter, r *http.Request) {

	filter := url.Values{}
	filter.Set("json_depth", "2")

	apps, err := mysql.SearchApps(filter, 0, "")
	if err != nil {
		logger.Error(err)
	}

	counts := make(map[int]int)

	for _, app := range apps {
		tags, err := app.GetTags()
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, tag := range tags {
			//logger.Info(genre.Description)

			if _, ok := counts[tag]; ok {
				counts[tag]++
			} else {
				counts[tag] = 1
			}
		}
	}

	for k, v := range counts {
		err := mysql.SaveOrUpdateTag(k, v)
		if err != nil {
			logger.Error(err)
		}
	}

	logger.Info("Tags updated")
}

func adminRanks(w http.ResponseWriter, r *http.Request) {

	var playersToRank = 500

	// Get keys, will delete any that are not removed from this map
	oldKeys, err := datastore.GetRankKeys()

	newRanks := make(map[int]*datastore.Rank)

	// Get players by level
	players, err := datastore.GetPlayers("-level", playersToRank)
	if err != nil {
		logger.Error(err)
		return
	}

	for k, v := range players {

		_, ok := newRanks[v.PlayerID]
		if !ok {

			rank := &datastore.Rank{}
			rank.FillFromPlayer(v)

			newRanks[v.PlayerID] = rank
		}
		newRanks[v.PlayerID].LevelRank = k + 1

		_, ok = oldKeys[strconv.Itoa(v.PlayerID)]
		if ok {
			delete(oldKeys, strconv.Itoa(v.PlayerID))
		}
	}

	// Convert new ranks to slice
	var ranks []*datastore.Rank
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
