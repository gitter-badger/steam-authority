package datastore

import (
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/streadway/amqp"
	"google.golang.org/api/iterator"
)

type Player struct {
	CreatedAt      time.Time                   `datastore:"created_at"`
	UpdatedAt      time.Time                   `datastore:"updated_at"`
	FriendsAddedAt time.Time                   `datastore:"friends_added_at"`
	PlayerID       int                         `datastore:"player_id"`
	ValintyURL     string                      `datastore:"vality_url,noindex"`
	Avatar         string                      `datastore:"avatar,noindex"`
	PersonaName    string                      `datastore:"persona_name,noindex"`
	RealName       string                      `datastore:"real_name,noindex"`
	CountryCode    string                      `datastore:"country_code"`
	StateCode      string                      `datastore:"status_code"`
	Level          int                         `datastore:"level"`
	Games          int                         `datastore:"games"`
	Badges         int                         `datastore:"badges"`
	PlayTime       int                         `datastore:"play_time"`
	TimeCreated    int                         `datastore:"time_created"` // In Steam's DB
	Friends        []steam.GetFriendListFriend `datastore:"friends,noindex"`
	AddFriends     bool                        `datastore:"-"` // Internal
}

func (p Player) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PLAYER, strconv.Itoa(p.PlayerID), nil)
}

func (p Player) GetPath() string {
	return "/players/" + strconv.Itoa(p.PlayerID) + "/" + slug.Make(p.PersonaName)
}

func GetPlayer(id int, produceFriends bool) (ret Player, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return ret, err
	}

	key := datastore.NameKey(PLAYER, strconv.Itoa(id), nil)

	player := &Player{}
	player.PlayerID = id
	player.AddFriends = produceFriends

	err = client.Get(context, key, player)
	if err != nil {

		// todo, clear latest players cache
		// Not in DB, go get it!
		if err.Error() == "datastore: no such entity" {

			player.Fill()
			_, err = SaveKind(player.GetKey(), player)
		}
		return *player, err
	}

	if player.UpdatedAt.Unix() < (time.Now().Unix() - int64(86400)) {
		player.Fill()
		//player.AddFriends = false
		_, err = SaveKind(player.GetKey(), player)
		return *player, err
	}

	// todo, add friendsScannedAt field????????????????????????
	//player.AddFriends = false
	return *player, nil
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

func ConsumePlayer(msg amqp.Delivery) (err error) {

	id := string(msg.Body)
	idx, _ := strconv.Atoi(id)

	_, err = GetPlayer(idx, false)

	return err
}

func (p *Player) Fill() (player *Player, err error) {

	// todo, get player bans, groups

	//Get summary
	summary, err := steam.GetPlayerSummaries([]int{p.PlayerID})
	if err != nil {
		return p, err
	} else if len(summary.Response.Players) > 0 {
		p.Avatar = summary.Response.Players[0].AvatarFull
		p.ValintyURL = path.Base(summary.Response.Players[0].ProfileURL)
		p.RealName = summary.Response.Players[0].RealName
		p.CountryCode = summary.Response.Players[0].LOCCountryCode
		p.StateCode = summary.Response.Players[0].LOCStateCode
		p.PersonaName = summary.Response.Players[0].PersonaName
	}

	//Get friends
	friends, err := steam.GetFriendList(p.PlayerID)
	if err != nil {
		logger.Error(err)
	} else {
		p.Friends = friends
	}

	// Dates
	p.UpdatedAt = time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	return p, nil
}
