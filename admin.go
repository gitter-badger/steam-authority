package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/Masterminds/squirrel"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/steam-authority/steam-authority/steam"
)

func reRankHandler(w http.ResponseWriter, r *http.Request) {

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

func fillAppsHandler(w http.ResponseWriter, r *http.Request) {

	// Get apps
	apps, err := steam.GetAppList()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, v := range apps {
		// Build query
		builder := squirrel.Insert("apps")
		builder = builder.Columns("id", "created_at", "updated_at", "name", "type", "is_free", "packages", "dlc", "categories", "genres", "screenshots", "movies", "achievements", "platforms")
		sql, args, err := builder.Values(v.AppID, int(time.Now().Unix()), int(time.Now().Unix()), v.Name, "", "0", "[]", "[]", "[]", "[]", "[]", "[]", "[]", "{}").ToSql()

		// Query
		_, err = mysql.ExecQuery(sql, args)
		if err != nil {
			logger.Error(err)
			//return
		}
	}

	logger.Info("Complete")
}
