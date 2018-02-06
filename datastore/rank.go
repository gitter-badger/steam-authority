package datastore

import (
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

var ranksLimit = 500

type Rank struct {
	CreatedAt   time.Time `datastore:"created_at"`
	UpdatedAt   time.Time `datastore:"updated_at"`
	PlayerID    int       `datastore:"player_id"`
	ValintyURL  string    `datastore:"vality_url"`
	Avatar      string    `datastore:"avatar"`
	PersonaName string    `datastore:"persona_name"`
	CountryCode string    `datastore:"country_code"`

	// Ranks
	Level           int `datastore:"level"`
	LevelRank       int `datastore:"level_rank"`
	Games           int `datastore:"games"`
	GamesRank       int `datastore:"games_rank"`
	Badges          int `datastore:"badges"`
	BadgesRank      int `datastore:"badges_rank"`
	PlayTime        int `datastore:"play_time"`
	PlayTimeRank    int `datastore:"play_time_rank"`
	TimeCreated     int `datastore:"time_created"`
	TimeCreatedRank int `datastore:"time_created_rank"`
	Friends         int `datastore:"friends"`
	FriendsRank     int `datastore:"friends_rank"`

	Rank int `datastore:"-"` // Just for the frontend
}

func (rank Rank) GetKey() (key *datastore.Key) {
	return datastore.NameKey(RANK, strconv.Itoa(rank.PlayerID), nil)
}

func (rank *Rank) Tidy() *Rank {

	rank.UpdatedAt = time.Now()
	if rank.CreatedAt.IsZero() {
		rank.CreatedAt = time.Now()
	}

	return rank
}

func (rank *Rank) FillFromPlayer(player Player) *Rank {

	rank.PlayerID = player.PlayerID
	rank.ValintyURL = player.ValintyURL
	rank.Avatar = player.Avatar
	rank.PersonaName = player.PersonaName
	rank.CountryCode = player.CountryCode
	rank.Level = player.Level
	rank.Games = player.Games
	rank.Badges = player.Badges
	rank.PlayTime = player.PlayTime
	rank.TimeCreated = player.TimeCreated
	rank.Friends = len(player.Friends)

	return rank
}

func GetRanksBy(order string) (ranks []Rank, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return ranks, err
	}

	q := datastore.NewQuery(RANK).Order(order).Limit(ranksLimit)
	it := client.Run(context, q)

	for {
		var dsRank Rank
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

func BulkSaveRanks(ranks []*Rank) (err error) {

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
