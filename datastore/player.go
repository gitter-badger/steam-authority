package datastore

import (
	"errors"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetPlayer(id64 string) (player DsPlayer, err error) {

	client, context := getDSClient()

	key := datastore.NameKey(PLAYER, id64, nil)

	err = client.Get(context, key, &player)
	if err != nil {
		logger.Error(err)
	}

	if player.LastUpdated < (time.Now().Unix() - int64(10)) { //todo, make a day?
		return player, errors.New("expired")
	}

	return player, err
}

func GetPlayers(order string, limit int) (players []DsPlayer, err error) {

	client, context := getDSClient()

	q := datastore.NewQuery(PLAYER).Order(order).Limit(limit)
	it := client.Run(context, q)

	for {
		var dsPlayer DsPlayer
		_, err := it.Next(&dsPlayer)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		players = append(players, dsPlayer)
	}

	return players, err
}

func SavePlayer(data DsPlayer) {

	key := datastore.NameKey(
		PLAYER,
		strconv.Itoa(data.ID64),
		nil,
	)

	saveKind(key, &data)
}
