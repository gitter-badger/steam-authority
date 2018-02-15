package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/queue"
	"github.com/steam-authority/steam-authority/steam"
)

func adminReRankHandler(w http.ResponseWriter, r *http.Request) {

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

func adminUpdateAllAppsHandler(w http.ResponseWriter, r *http.Request) {

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
