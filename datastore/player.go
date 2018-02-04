package datastore

import (
	"errors"
	"path"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
	"google.golang.org/api/iterator"
)

type DsPlayer struct {
	CreatedAt   time.Time `datastore:"created_at"`
	UpdatedAt   time.Time `datastore:"updated_at"`
	ID64        int       `datastore:"id64"`
	ValintyURL  string    `datastore:"vality_url"`
	Avatar      string    `datastore:"avatar"`
	RealName    string    `datastore:"real_name"`
	PersonaName string    `datastore:"persona_name"`
	CountryCode string    `datastore:"country_code"`
	StateCode   string    `datastore:"status_code"`
	TimeUpdated int64     `datastore:"time_updated"` // todo, Remove?
	Level       int       `datastore:"level"`
	Games       int       `datastore:"games"`
	Badges      int       `datastore:"badges"`
	PlayTime    int       `datastore:"play_time"`
	TimeCreated int       `datastore:"time_created"`
	Friends     []int     `datastore:"friends"`

	Rank int `datastore:"-"`
}

func (player DsPlayer) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PLAYER, strconv.Itoa(player.ID64), nil)
}

func (player DsPlayer) GetPath() string {
	return "/players/" + strconv.Itoa(player.ID64) + "/" + slug.Make(player.PersonaName)
}

func (player *DsPlayer) Tidy() *DsPlayer {

	player.UpdatedAt = time.Now()
	if player.CreatedAt.IsZero() {
		player.CreatedAt = time.Now()
	}

	return player
}

func (player *DsPlayer) FillFromSummary(summary steam.PlayerSummariesBody) *DsPlayer {

	if len(summary.Response.Players) > 0 {
		player.Avatar = summary.Response.Players[0].AvatarFull
		player.ValintyURL = path.Base(summary.Response.Players[0].ProfileURL)
		player.RealName = summary.Response.Players[0].RealName
		player.CountryCode = summary.Response.Players[0].LOCCountryCode
		player.StateCode = summary.Response.Players[0].LOCStateCode
		player.PersonaName = summary.Response.Players[0].PersonaName
	}

	return player
}

func (player *DsPlayer) FillFromFriends(summary []steam.GetFriendListFriend) *DsPlayer {

	var friends []int
	for _, v := range summary {
		i, _ := strconv.Atoi(v.Steamid)
		friends = append(friends, i)
	}

	player.Friends = friends
	return player
}

// todo, Only return 1 player not slice
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
