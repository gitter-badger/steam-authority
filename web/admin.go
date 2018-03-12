package web

import (
	"encoding/json"
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
	case "queues":
		r.ParseForm()
		go adminQueues(w, r, r.PostForm)
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
		bytes, _ := json.Marshal(queue.AppMessage{
			AppID:    v.AppID,
			ChangeID: 0,
		})

		queue.Produce(queue.AppQueue, bytes)
	}

	logger.Info(strconv.Itoa(len(apps)) + " apps added to rabbit")
}

func adminDeploy(w http.ResponseWriter, r *http.Request) {

}

func adminDonations(w http.ResponseWriter, r *http.Request) {

	donations, err := datastore.GetDonations(0, 0)
	if err != nil {
		logger.Error(err)
		return
	}

	// ma[player]total
	counts := make(map[int]int)

	for _, v := range donations {

		if _, ok := counts[v.PlayerID]; ok {
			counts[v.PlayerID] = counts[v.PlayerID] + v.AmountUSD
		} else {
			counts[v.PlayerID] = v.AmountUSD
		}
	}

	for k, v := range counts {
		player, err := datastore.GetPlayer(k)
		if err != nil {
			logger.Error(err)
			continue
		}

		player.Donated = v
		_, err = datastore.SaveKind(player.GetKey(), player)
	}

	logger.Info("Updated " + strconv.Itoa(len(counts)) + " player donation counts")
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

func adminQueues(w http.ResponseWriter, r *http.Request, form url.Values) {

	if val := form.Get("change-id"); val != "" {

		logger.Info("Change: " + val)
		appID, _ := strconv.Atoi(val)
		bytes, _ := json.Marshal(queue.AppMessage{
			AppID: appID,
		})
		queue.Produce(queue.AppQueue, bytes)
	}

	if val := form.Get("player-id"); val != "" {

		logger.Info("Player: " + val)
		playerID, _ := strconv.Atoi(val)
		bytes, _ := json.Marshal(queue.PlayerMessage{
			PlayerID: playerID,
		})
		queue.Produce(queue.PlayerQueue, bytes)
	}

	if val := form.Get("app-id"); val != "" {

		logger.Info("App: " + val)
		appID, _ := strconv.Atoi(val)
		bytes, _ := json.Marshal(queue.AppMessage{
			AppID: appID,
		})
		queue.Produce(queue.AppQueue, bytes)
	}

	if val := form.Get("package-id"); val != "" {

		logger.Info("Package: " + val)
		packageID, _ := strconv.Atoi(val)
		bytes, _ := json.Marshal(queue.PackageMessage{
			PackageID: packageID,
		})
		queue.Produce(queue.PackageQueue, bytes)
	}

}

// todo, handle tags that no longer have any games.
func adminTags(w http.ResponseWriter, r *http.Request) {

	filter := url.Values{}
	filter.Set("json_depth", "3")

	apps, err := mysql.SearchApps(filter, 0, "")
	if err != nil {
		logger.Error(err)
	}

	// map[player]struct
	counts := make(map[int]*adminTag)

	for _, app := range apps {
		tags, err := app.GetTags()
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, tag := range tags {
			//logger.Info(genre.Description)

			if _, ok := counts[tag]; ok {
				counts[tag].count++
				counts[tag].totalPrice = counts[tag].totalPrice + app.PriceFinal
				counts[tag].totalDiscount = counts[tag].totalDiscount + app.PriceDiscount
			} else {
				counts[tag] = &adminTag{
					count:         1,
					totalPrice:    app.PriceFinal,
					totalDiscount: app.PriceDiscount,
				}
			}
		}
	}

	for k, v := range counts {
		err := mysql.SaveOrUpdateTag(k, mysql.Tag{
			Apps:         v.count,
			MeanPrice:    v.GetMeanPrice(),
			MeanDiscount: v.GetMeanDiscount(),
		})
		if err != nil {
			logger.Error(err)
		}
	}

	logger.Info("Tags updated")
}

type adminTag struct {
	count         int
	totalPrice    int
	totalDiscount int
}

func (t adminTag) GetMeanPrice() float64 {
	return float64(t.totalPrice) / float64(t.count)
}

func (t adminTag) GetMeanDiscount() float64 {
	return float64(t.totalDiscount) / float64(t.count)
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
