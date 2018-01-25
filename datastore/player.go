package datastore

import (
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetPlayer(id64 string) (player DsPlayer, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return player, err
	}

	key := datastore.NameKey(PLAYER, id64, nil)

	err = client.Get(context, key, &player)
	if err != nil {
		logger.Error(err)
	}

	// Error if data is older than a day
	if player.TimeUpdated < (time.Now().Unix() - int64(86400)) {
		return player, errors.New("expired")
	}

	return player, nil
}

func GetPlayers(order string, limit int) (players []DsPlayer, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return players, err
	}

	q := datastore.NewQuery(PLAYER).Order(order).Limit(limit)
	it := client.Run(ctx, q)

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

func CountPlayers() (count int, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return count, err
	}

	q := datastore.NewQuery(PLAYER)
	count, err = client.Count(ctx, q)
	if err != nil {
		return count, err
	}

	return count, nil
}
