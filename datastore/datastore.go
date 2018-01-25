package datastore

import (
	"context"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
)

const (
	// CHANGE datastore kind
	CHANGE = "Change"

	// APP datastore kind
	APP = "App"

	// PACKAGE datastore kind
	PACKAGE = "Package"

	// PLAYER datastore kind
	PLAYER = "Player"

	// RANK datastore kind
	RANK = "Rank"
)

func SaveKind(key *datastore.Key, data interface{}) (newKey *datastore.Key, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return nil, err
	}

	newKey, err = client.Put(ctx, key, data)
	if err != nil {
		return newKey, err
	}

	return newKey, nil
}

func getDSClient() (client *datastore.Client, ctx context.Context, err error) {

	ctx = context.Background()
	client, err = datastore.NewClient(ctx, os.Getenv("STEAM_GOOGLE_PROJECT"))
	if err != nil {
		return client, ctx, err
	}

	return client, ctx, nil
}

// DsPlayer has everything visible on a player page
type DsPlayer struct {
	ID64        int
	ValintyURL  string `datastore:"vality_url"`
	Avatar      string `datastore:"avatar"`
	RealName    string `datastore:"real_name"`
	PersonaName string `datastore:"persona_name"`
	CountryCode string `datastore:"country_code"`
	StateCode   string `datastore:"status_code"`
	TimeUpdated int64  `datastore:"time_updated"`
	Level       int    `datastore:"level"`
	Games       int    `datastore:"games"`
	Badges      int    `datastore:"badges"`
	PlayTime    int    `datastore:"play_time"`
	TimeCreated int    `datastore:"time_created"`
	Friends     []int  `datastore:"friends"`

	Rank int `datastore:"-"`
}

func (p *DsPlayer) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PLAYER, strconv.Itoa(p.ID64), nil)
}

// DsRank has only the things that need to be visible on the frontend ranks page
type DsRank struct {
	ID64            int
	ValintyURL      string `datastore:"vality_url"`
	Avatar          string `datastore:"avatar"`
	PersonaName     string `datastore:"persona_name"`
	CountryCode     string `datastore:"country_code"`
	Level           int    `datastore:"level"`
	LevelRank       int    `datastore:"level_rank"`
	Games           int    `datastore:"games"`
	GamesRank       int    `datastore:"games_rank"`
	Badges          int    `datastore:"badges"`
	BadgesRank      int    `datastore:"badges_rank"`
	PlayTime        int    `datastore:"play_time"`
	PlayTimeRank    int    `datastore:"play_time_rank"`
	TimeCreated     int    `datastore:"time_created"`
	TimeCreatedRank int    `datastore:"time_created_rank"`
	Friends         int    `datastore:"friends"`
	FriendsRank     int    `datastore:"friends_rank"`

	Rank int `datastore:"-"` // Just for the frontend
}

func (rank *DsRank) GetKey() (key *datastore.Key) {
	return datastore.NameKey(RANK, strconv.Itoa(rank.ID64), nil)
}

// DsChange kind
type DsChange struct {
	ChangeID int   `datastore:"change_id"`
	Apps     []int `datastore:"apps"`
	Packages []int `datastore:"packages"`
}

func (change *DsChange) GetKey() (key *datastore.Key) {
	return datastore.NameKey(CHANGE, strconv.Itoa(change.ChangeID), nil)
}

// DsApp kind
type DsApp struct {
	AppID             int      `datastore:"app_id"`
	Name              string   `datastore:"name"`
	Type              string   `datastore:"type"`
	ReleaseState      string   `datastore:"releasestate"`
	OSList            []string `datastore:"oslist"`
	MetacriticScore   int8     `datastore:"metacritic_score"`
	MetacriticFullURL string   `datastore:"metacritic_fullurl"`
	StoreTags         []int    `datastore:"store_tags"`
	Developer         string   `datastore:"developer"`
	Publisher         string   `datastore:"publisher"`
	Homepage          string   `datastore:"homepage"`
	ChangeNumber      int      `datastore:"change_number"`
	Logo              string   `datastore:"logo"`
	Icon              string   `datastore:"icon"`
}

func (app *DsApp) GetKey() (key *datastore.Key) {
	return datastore.NameKey(APP, strconv.Itoa(app.AppID), nil)
}

// DsPackage kind
type DsPackage struct {
	PackageID   int   `datastore:"package_id"`
	BillingType int8  `datastore:"billingtype"`
	LicenseType int8  `datastore:"licensetype"`
	Status      int8  `datastore:"status"`
	Apps        []int `datastore:"apps"`
	ChangeID    int   `datastore:"change_id"`
}

func (packagex *DsPackage) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PACKAGE, strconv.Itoa(packagex.PackageID), nil)
}
