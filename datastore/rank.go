package datastore

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

var ranksLimit = 500

func GetRanksBy(order string) (ranks []DsRank, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return ranks, err
	}

	q := datastore.NewQuery(RANK).Order(order).Limit(ranksLimit)
	it := client.Run(context, q)

	for {
		var dsRank DsRank
		_, err := it.Next(&dsRank)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		ranks = append(ranks, dsRank)
	}

	return ranks, err
}

func BulkSaveRanks(ranks []*DsRank) (err error) {

	RanksLen := len(ranks)
	if RanksLen == 0 {
		return nil
	}

	client, context, err := getDSClient()
	if err != nil {
		return err
	}

	keys := make([]*datastore.Key, 0, RanksLen)
	for _, v := range ranks {
		keys = append(keys, v.GetKey())
	}

	fmt.Println("Saving " + strconv.Itoa(RanksLen) + " ranks")
	_, err = client.PutMulti(context, keys, ranks)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func GetRankKeys() (keysMap map[string]*datastore.Key, err error) {

	keysMap = make(map[string]*datastore.Key)

	client, ctx, err := getDSClient()
	if err != nil {
		return keysMap, err
	}

	q := datastore.NewQuery(RANK).KeysOnly().Limit(1000)
	keys, err := client.GetAll(ctx, q, nil)
	if err != nil {
		return keysMap, err
	}

	for _, v := range keys {
		keysMap[v.Name] = v
	}

	return keysMap, nil
}

func CountRankedPlayers() (count int, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return count, err
	}

	q := datastore.NewQuery(RANK)
	count, err = client.Count(ctx, q)
	if err != nil {
		return count, err
	}

	return count, nil
}

func BulkDeleteRanks(keys map[string]*datastore.Key) (err error) {

	// Make map a slice
	var keysToDelete []*datastore.Key
	for _, v := range keys {
		keysToDelete = append(keysToDelete, v)
	}

	client, ctx, err := getDSClient()
	if err != nil {
		return err
	}

	err = client.DeleteMulti(ctx, keysToDelete)
	if err != nil {
		return err
	}

	return nil
}
