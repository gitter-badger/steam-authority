package datastore

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
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

func getDSClient() (client *datastore.Client, ctx context.Context, err error) {

	ctx = context.Background()
	client, err = datastore.NewClient(ctx, os.Getenv("STEAM_GOOGLE_PROJECT"))
	if err != nil {
		return client, ctx, err
	}

	return client, ctx, nil
}

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

// DsPlayer has everything visible on a player page
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
	TimeUpdated int64     `datastore:"time_updated"`
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

	player.Avatar = summary.Response.Players[0].AvatarFull
	player.ValintyURL = path.Base(summary.Response.Players[0].ProfileURL)
	player.RealName = summary.Response.Players[0].RealName
	player.CountryCode = summary.Response.Players[0].LOCCountryCode
	player.StateCode = summary.Response.Players[0].LOCStateCode
	player.PersonaName = summary.Response.Players[0].PersonaName

	return player
}

// DsRank has only the things that need to be visible on the frontend ranks page
type DsRank struct {
	CreatedAt       time.Time `datastore:"created_at"`
	UpdatedAt       time.Time `datastore:"updated_at"`
	ID64            int       `datastore:"id64"`
	ValintyURL      string    `datastore:"vality_url"`
	Avatar          string    `datastore:"avatar"`
	PersonaName     string    `datastore:"persona_name"`
	CountryCode     string    `datastore:"country_code"`
	Level           int       `datastore:"level"`
	LevelRank       int       `datastore:"level_rank"`
	Games           int       `datastore:"games"`
	GamesRank       int       `datastore:"games_rank"`
	Badges          int       `datastore:"badges"`
	BadgesRank      int       `datastore:"badges_rank"`
	PlayTime        int       `datastore:"play_time"`
	PlayTimeRank    int       `datastore:"play_time_rank"`
	TimeCreated     int       `datastore:"time_created"`
	TimeCreatedRank int       `datastore:"time_created_rank"`
	Friends         int       `datastore:"friends"`
	FriendsRank     int       `datastore:"friends_rank"`

	Rank int `datastore:"-"` // Just for the frontend
}

func (rank DsRank) GetKey() (key *datastore.Key) {
	return datastore.NameKey(RANK, strconv.Itoa(rank.ID64), nil)
}

func (rank *DsRank) Tidy() *DsRank {

	rank.UpdatedAt = time.Now()
	if rank.CreatedAt.IsZero() {
		rank.CreatedAt = time.Now()
	}

	return rank
}

func (rank *DsRank) FillFromPlayer(player DsPlayer) *DsRank {

	rank.ID64 = player.ID64
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

// DsChange kind
type DsChange struct {
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`
	ChangeID  int       `datastore:"change_id"`
	Apps      []int     `datastore:"apps"`
	Packages  []int     `datastore:"packages"`
}

func (change DsChange) GetKey() (key *datastore.Key) {
	return datastore.NameKey(CHANGE, strconv.Itoa(change.ChangeID), nil)
}

func (change *DsChange) Tidy() *DsChange {

	change.UpdatedAt = time.Now()
	if change.CreatedAt.IsZero() {
		change.CreatedAt = time.Now()
	}

	return change
}

// DsApp kind
type DsApp struct {
	AppID     int       `datastore:"app_id"`
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`

	// In API call
	Name              string   `datastore:"name"`
	Type              string   `datastore:"type"`
	Free              bool     `datastore:"is_free"`
	DLC               []int    `datastore:"dlc"`
	ShortDescription  string   `datastore:"short_description"`
	HeaderImage       string   `datastore:"header_image"`
	Developers        []string `datastore:"developer"`
	Publishers        []string `datastore:"publisher"`
	Packages          []int    `datastore:"packages"`
	MetacriticScore   int8     `datastore:"metacritic_score"`
	MetacriticFullURL string   `datastore:"metacritic_fullurl"`
	Categories        []int8   `datastore:"categories"`
	Genres            []int8   `datastore:"genres"`
	Screenshots       string   `datastore:"screenshots"`
	Movies            string   `datastore:"movies"`
	Achievements      string   `datastore:"achievements"`
	Background        string   `datastore:"background"`
	OSList            []string `datastore:"oslist"`

	// These are not in the api call, only in pics
	ReleaseState string `datastore:"releasestate"`
	StoreTags    []int  `datastore:"store_tags"`
	Homepage     string `datastore:"homepage"`
	ChangeNumber int    `datastore:"change_number"`
	Logo         string `datastore:"logo"`
	Icon         string `datastore:"icon"`

	// Struct versions of stores JSON
	ScreenshotStruct   []steam.AppDetailsScreenshot `datastore:"-"`
	AchievementsStruct steam.AppDetailsAchievements `datastore:"-"`
}

func (app DsApp) GetKey() (key *datastore.Key) {
	return datastore.NameKey(APP, strconv.Itoa(app.AppID), nil)
}

func (app *DsApp) Tidy() *DsApp {

	app.UpdatedAt = time.Now()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}

	app.Type = strings.ToLower(app.Type)
	app.ReleaseState = strings.ToLower(app.ReleaseState)

	if app.Name == "" {
		app.Name = "App " + strconv.Itoa(app.AppID)
	}

	return app
}

func (app *DsApp) GetPath() string {
	return "/apps/" + strconv.Itoa(app.AppID) + "/" + slug.Make(app.Name)
}

func (app *DsApp) FillFromAppDetails(appDetails steam.AppDetailsBody) *DsApp {

	// Screenshots
	screenshotsString, err := json.Marshal(appDetails.Data.Screenshots)
	if err != nil {
		logger.Error(err)
	}

	// Movies
	moviesString, err := json.Marshal(appDetails.Data.Movies)
	if err != nil {
		logger.Error(err)
	}

	// Achievements
	achievementsString, err := json.Marshal(appDetails.Data.Achievements)
	if err != nil {
		logger.Error(err)
	}

	// Categories
	var categories []int8
	for _, v := range appDetails.Data.Categories {
		categories = append(categories, v.ID)
	}

	// Genres
	var genres []int8
	for _, v := range appDetails.Data.Genres {
		genre, _ := strconv.ParseInt(v.ID, 10, 8)
		genres = append(genres, int8(genre))
	}

	// Platforms
	var platforms []string
	if appDetails.Data.Platforms.Linux {
		platforms = append(platforms, "linux")
	}
	if appDetails.Data.Platforms.Windows {
		platforms = append(platforms, "windows")
	}
	if appDetails.Data.Platforms.Windows {
		platforms = append(platforms, "macos")
	}

	//
	app.Name = appDetails.Data.Name
	app.Type = appDetails.Data.Type
	app.Free = appDetails.Data.IsFree
	app.DLC = appDetails.Data.DLC
	app.ShortDescription = appDetails.Data.ShortDescription
	app.HeaderImage = appDetails.Data.HeaderImage
	app.Developers = appDetails.Data.Developers
	app.Publishers = appDetails.Data.Publishers
	app.Packages = appDetails.Data.Packages
	app.MetacriticScore = appDetails.Data.Metacritic.Score
	app.MetacriticFullURL = appDetails.Data.Metacritic.URL
	app.Categories = categories
	app.Genres = genres
	app.Screenshots = string(screenshotsString)
	app.Movies = string(moviesString)
	app.Achievements = string(achievementsString)
	app.Background = appDetails.Data.Background
	app.OSList = platforms

	return app
}

func (app *DsApp) FillFromJSON() *DsApp {

	var bytes []byte

	// Screenshots
	bytes = []byte(app.Screenshots)
	if err := json.Unmarshal(bytes, &app.ScreenshotStruct); err != nil {
		logger.Error(err)
	}

	// Achievements
	bytes = []byte(app.Achievements)
	if err := json.Unmarshal(bytes, &app.AchievementsStruct); err != nil {
		logger.Error(err)
	}

	return app
}

// DsPackage kind
type DsPackage struct {
	CreatedAt   time.Time `datastore:"created_at"`
	UpdatedAt   time.Time `datastore:"updated_at"`
	PackageID   int       `datastore:"package_id"`
	BillingType int8      `datastore:"billingtype"`
	LicenseType int8      `datastore:"licensetype"`
	Status      int8      `datastore:"status"`
	Apps        []int     `datastore:"apps"`
	ChangeID    int       `datastore:"change_id"`
}

func (packagex DsPackage) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PACKAGE, strconv.Itoa(packagex.PackageID), nil)
}

func (packagex *DsPackage) Tidy() *DsPackage {

	packagex.UpdatedAt = time.Now()
	if packagex.CreatedAt.IsZero() {
		packagex.CreatedAt = time.Now()
	}

	return packagex
}
