package datastore

import (
	"errors"
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
	"google.golang.org/api/iterator"
)

type Player struct {
	CreatedAt   time.Time `datastore:"created_at"`
	UpdatedAt   time.Time `datastore:"updated_at"`
	PlayerID    int       `datastore:"player_id"`
	ValintyURL  string    `datastore:"vality_url"`
	Avatar      string    `datastore:"avatar"`
	PersonaName string    `datastore:"persona_name"`
	RealName    string    `datastore:"real_name"`
	CountryCode string    `datastore:"country_code"`
	StateCode   string    `datastore:"status_code"`
	Level       int       `datastore:"level"`
	Games       int       `datastore:"games"`
	Badges      int       `datastore:"badges"`
	PlayTime    int       `datastore:"play_time"`
	TimeCreated int       `datastore:"time_created"` // In Steam's DB
	Friends     []int     `datastore:"friends"`
}

func (player Player) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PLAYER, strconv.Itoa(player.PlayerID), nil)
}

func (player Player) GetPath() string {
	return "/players/" + strconv.Itoa(player.PlayerID) + "/" + slug.Make(player.PersonaName)
}

func (player *Player) Tidy() *Player {

	player.UpdatedAt = time.Now()
	if player.CreatedAt.IsZero() {
		player.CreatedAt = time.Now()
	}

	return player
}

func (player *Player) FillFromSummary(summary steam.PlayerSummariesBody) *Player {

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

func (player *Player) FillFromFriends(summary []steam.GetFriendListFriend) *Player {

	var friends []int
	for _, v := range summary {
		i, _ := strconv.Atoi(v.Steamid)
		friends = append(friends, i)
	}

	player.Friends = friends
	return player
}

// todo, Only return 1 player not slice
func GetPlayer(id64 string) (player Player, err error) {

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
	if player.UpdatedAt.Unix() < (time.Now().Unix() - int64(86400)) {
		return player, errors.New("expired")
	}

	return player, nil
}

func GetPlayers(order string, limit int) (players []Player, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return players, err
	}

	q := datastore.NewQuery(PLAYER).Order(order).Limit(limit)
	it := client.Run(ctx, q)

	for {
		var dsPlayer Player
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

func GetPlayersByIDs(ids []int) (friends []Player, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return friends, err
	}

	var keys []*datastore.Key
	for _, v := range ids {
		key := datastore.NameKey(PLAYER, strconv.Itoa(v), nil)
		keys = append(keys, key)
	}

	friends = make([]Player, len(keys))
	err = client.GetMulti(context, keys, friends)
	if err != nil && !strings.Contains(err.Error(), "no such entity") {
		return friends, err
	}

	return friends, nil
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
