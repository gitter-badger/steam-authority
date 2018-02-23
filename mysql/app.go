package mysql

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/gosimple/slug"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/streadway/amqp"
)

type App struct {
	ID                int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"` //
	CreatedAt         *time.Time `gorm:"not null;column:created_at"`                    //
	UpdatedAt         *time.Time `gorm:"not null;column:updated_at"`                    //
	Name              string     `gorm:"not null;column:name"`                          //
	Type              string     `gorm:"not null;column:type"`                          //
	IsFree            bool       `gorm:"not null;column:is_free;type:tinyint(1)"`       //
	DLC               string     `gorm:"not null;column:dlc;default:'[]'"`              // JSON
	ShortDescription  string     `gorm:"not null;column:description_short"`             //
	HeaderImage       string     `gorm:"not null;column:image_header"`                  //
	Developers        string     `gorm:"not null;column:developer;default:'[]'"`        // JSON
	Publishers        string     `gorm:"not null;column:publisher;default:'[]'"`        // JSON
	Packages          string     `gorm:"not null;column:packages;default:'[]'"`         // JSON
	MetacriticScore   int8       `gorm:"not null;column:metacritic_score"`              //
	MetacriticFullURL string     `gorm:"not null;column:metacritic_url"`                //
	Categories        string     `gorm:"not null;column:categories;default:'[]'"`       // JSON
	Genres            string     `gorm:"not null;column:genres;default:'[]'"`           // JSON
	Screenshots       string     `gorm:"not null;column:screenshots;default:'[]'"`      // JSON
	Movies            string     `gorm:"not null;column:movies;default:'[]'"`           // JSON
	Achievements      string     `gorm:"not null;column:achievements;default:'[]'"`     // JSON
	Background        string     `gorm:"not null;column:background"`                    //
	Platforms         string     `gorm:"not null;column:platforms;default:'[]'"`        // JSON
	GameID            int        `gorm:"not null;column:game_id"`                       //
	GameName          string     `gorm:"not null;column:game_name"`                     //
	ReleaseState      string     `gorm:"not null;column:release_state"`                 // PICS
	StoreTags         string     `gorm:"not null;column:tags;default:'[]';type:json"`   // PICS JSON
	Homepage          string     `gorm:"not null;column:homepage"`                      // PICS
	ChangeNumber      int        `gorm:"not null;column:change_number"`                 // PICS
	Logo              string     `gorm:"not null;column:logo"`                          // PICS
	Icon              string     `gorm:"not null;column:icon"`                          // PICS
	ClientIcon        string     `gorm:"not null;column:client_icon"`                   // PICS
}

func (app App) GetPath() (ret string) {
	ret = "/apps/" + strconv.Itoa(int(app.ID))

	if app.Name != "" {
		ret = ret + "/" + slug.Make(app.Name)
	}

	return ret
}

// advertising, demo, dlc, episode, game, mod, movie, series, tool, video
func (app App) GetType() (ret string) {

	switch app.Type {
	case "dlc":
		return "DLC"
	default:
		return strings.Title(app.Type)
	}
}

func (app App) GetCommunityLink() (string) {
	return "https://steamcommunity.com/app/" + strconv.Itoa(app.ID) + "/?utm_source=SteamAuthority&utm_medium=SteamAuthority&utm_campaign=SteamAuthority"
}

func (app App) GetStoreLink() (string) {
	return "https://store.steampowered.com/app/" + strconv.Itoa(app.ID) + "/?utm_source=SteamAuthority&utm_medium=SteamAuthority&utm_campaign=SteamAuthority"
}

func (app App) GetInstallLink() (string) {
	return "steam://install/" + strconv.Itoa(app.ID)
}

// Used in frontend
func (app App) GetScreenshots() (screenshots []steam.AppDetailsScreenshot, err error) {

	bytes := []byte(app.Screenshots)
	if err := json.Unmarshal(bytes, &screenshots); err != nil {
		return screenshots, err
	}

	return screenshots, nil
}

func (app App) GetAchievements() (achievements steam.AppDetailsAchievements, err error) {

	bytes := []byte(app.Achievements)
	if err := json.Unmarshal(bytes, &achievements); err != nil {
		return achievements, err
	}

	return achievements, nil
}

func (app App) GetPlatforms() (platforms []string, err error) {

	bytes := []byte(app.Platforms)
	if err := json.Unmarshal(bytes, &platforms); err != nil {
		return platforms, err
	}

	return platforms, nil
}

func (app App) GetDLC() (dlcs []int, err error) {

	bytes := []byte(app.DLC)
	if err := json.Unmarshal(bytes, &dlcs); err != nil {
		return dlcs, err
	}

	return dlcs, nil
}

func (app App) GetPackages() (packages []int, err error) {

	bytes := []byte(app.Packages)
	if err := json.Unmarshal(bytes, &packages); err != nil {
		return packages, err
	}

	return packages, nil
}

func (app App) GetGenres() (genres []steam.AppDetailsGenre, err error) {

	bytes := []byte(app.Genres)
	if err := json.Unmarshal(bytes, &genres); err != nil {
		return genres, err
	}

	return genres, nil
}

func (app App) GetCategories() (categories []string, err error) {

	bytes := []byte(app.Categories)
	if err := json.Unmarshal(bytes, &categories); err != nil {
		return categories, err
	}

	return categories, nil
}

func (app App) GetName() (name string) {

	if app.Name == "" {
		app.Name = "App " + strconv.Itoa(int(app.ID))
	}

	return app.Name
}

func GetApp(id int) (app App, err error) {

	// Connect
	db, err := getDB()
	if err != nil {
		return app, err
	}

	// Get app
	db.FirstOrInit(&app, App{ID: id})
	if db.Error != nil {
		return app, db.Error
	}

	// If new or expired app
	if app.UpdatedAt.Unix() < time.Now().AddDate(0, 0, -1).Unix() {

		err = app.Save()
		if err != nil {
			return app, err
		}
	}

	// Don't bother checking steam to see if it exists, we should know about all apps.

	return app, nil
}

func GetApps(ids []int) (apps []App, err error) {

	if len(ids) == 0 {
		return apps, nil
	}

	db, err := getDB()
	if err != nil {
		return apps, err
	}

	db.Where("id IN (?)", ids).Find(&apps)
	if db.GetErrors() != nil {
		return apps, err
	}

	return apps, nil
}

func SearchApps(query url.Values, limit int, sort string) (apps []App, err error) {

	db, err := getDB()
	if err != nil {
		return apps, err
	}

	if sort == "" {
		sort = "id DESC" // todo, order by popularity?
	}

	db = db.Limit(limit)
	db = db.Order(sort)

	// Free
	if _, ok := query["is_free"]; ok {
		db = db.Where("is_free = ?", query.Get("is_free"))
	}

	// Platforms
	if _, ok := query["platforms"]; ok {
		db = db.Where("JSON_CONTAINS(platforms, [\"?\"])", query.Get("platforms"))

	}

	// Tag
	if _, ok := query["tags"]; ok {
		db = db.Where("JSON_CONTAINS(tags, ?)", "[\""+query.Get("tags")+"\"]")
	}

	// Query
	db = db.Find(&apps)
	if db.Error != nil {
		return apps, err
	}

	return apps, err
}

func CountApps() (count int, err error) {

	db, err := getDB()
	if err != nil {
		return count, err
	}

	db.Model(&App{}).Count(&count)
	if db.Error != nil {
		return count, err
	}

	return count, nil
}

func NewApp(id int) (app App) {

	app.ID = id
	return app
}

func (app *App) Save() (err error) {

	// Save
	db, err := getDB()
	if err != nil {
		return err
	}

	db.Save(&app)
	if db.Error != nil {
		return err
	}

	return nil
}

// GORM callback
// If any callback returns an error, GORM will stop future operations and rollback current transaction.
func (app *App) BeforeSave() {

	var err error

	// Get app details
	err = app.FillFromAppDetails()
	if err != nil {
		if err.Error() != "no app with id" {
			logger.Error(err)
		}
	}

	err = app.FillFromPICS()
	if err != nil {
		logger.Error(err)
	}

	// Tidy
	app.Type = strings.ToLower(app.Type)
	app.ReleaseState = strings.ToLower(app.ReleaseState)

	if app.StoreTags == "" || app.StoreTags == "null" {
		app.StoreTags = "[]"
	}

	if app.Developers == "" || app.Developers == "null" {
		app.Developers = "[]"
	}

	if app.Publishers == "" || app.Publishers == "null" {
		app.Publishers = "[]"
	}

	if app.Categories == "" || app.Categories == "null" {
		app.Categories = "[]"
	}

	if app.Genres == "" || app.Genres == "null" {
		app.Genres = "[]"
	}

	if app.Screenshots == "" || app.Screenshots == "null" {
		app.Screenshots = "[]"
	}

	if app.Movies == "" || app.Movies == "null" {
		app.Movies = "[]"
	}

	if app.Achievements == "" || app.Achievements == "null" {
		app.Achievements = "{}"
	}

	if app.Platforms == "" || app.Platforms == "null" {
		app.Platforms = "[]"
	}

	if app.DLC == "" || app.DLC == "null" {
		app.DLC = "[]"
	}

	if app.Packages == "" || app.Packages == "null" {
		app.Packages = "[]"
	}
}

func ConsumeApp(msg amqp.Delivery) (err error) {

	id := string(msg.Body)
	idx, _ := strconv.Atoi(id)

	app := NewApp(idx)
	err = app.Save()

	return err
}

func (app *App) FillFromPICS() (err error) {

	// Call PICS
	resp, err := steam.GetPICSInfo([]int{app.ID}, []int{})
	if err != nil {
		return err
	}

	var js steam.JsApp
	if val, ok := resp.Apps[app.ID]; ok {
		js = val
	} else {
		return errors.New("no app key in json")
	}

	// Tags
	tags, err := json.Marshal(js.Common.StoreTags)
	if err != nil {
		return err
	}

	// String to int
	var metacriticScoreInt int
	if js.Common.MetacriticScore == "" {
		metacriticScoreInt = 0
	} else {
		metacriticScoreInt, err = strconv.Atoi(js.Common.MetacriticScore)
		if err != nil {
			return err
		}
	}

	//
	app.Name = js.Common.Name
	app.Type = js.Common.Type
	app.ReleaseState = js.Common.ReleaseState
	//app.Platforms = strings.Split(js.Common.OSList, ",") // Can get from API
	app.MetacriticScore = int8(metacriticScoreInt)
	app.MetacriticFullURL = js.Common.MetacriticURL
	app.StoreTags = string(tags)
	app.Developers = js.Extended.Developer
	app.Publishers = js.Extended.Publisher
	app.Homepage = js.Extended.Homepage
	app.ChangeNumber = js.ChangeNumber
	app.Logo = js.Common.Logo
	app.Icon = js.Common.Icon
	app.ClientIcon = js.Common.ClientIcon

	return nil
}

func (app *App) FillFromAppDetails() (err error) {

	// Get data
	appDetails, err := steam.GetAppDetails(strconv.Itoa(app.ID))
	if err != nil {
		return err
	}

	// Screenshots
	screenshotsString, err := json.Marshal(appDetails.Data.Screenshots)
	if err != nil {
		return err
	}

	// Movies
	moviesString, err := json.Marshal(appDetails.Data.Movies)
	if err != nil {
		return err
	}

	// Achievements
	achievementsString, err := json.Marshal(appDetails.Data.Achievements)
	if err != nil {
		return err
	}

	// DLC
	dlcString, err := json.Marshal(appDetails.Data.DLC)
	if err != nil {
		return err
	}

	// Developers
	developersString, err := json.Marshal(appDetails.Data.Developers)
	if err != nil {
		return err
	}

	// Publishers
	publishersString, err := json.Marshal(appDetails.Data.Publishers)
	if err != nil {
		return err
	}

	// Packages
	packagesString, err := json.Marshal(appDetails.Data.Packages)
	if err != nil {
		return err
	}

	// Categories
	var categories []int8
	for _, v := range appDetails.Data.Categories {
		categories = append(categories, v.ID)
	}

	categoriesString, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	// Genres
	var genres []int8
	for _, v := range appDetails.Data.Genres {
		//genre, _ := strconv.ParseInt(v.ID, 10, 8)
		genres = append(genres, v.ID)
	}

	genresString, err := json.Marshal(genres)
	if err != nil {
		return err
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

	platformsString, err := json.Marshal(platforms)
	if err != nil {
		return err
	}

	//
	app.Name = appDetails.Data.Name
	app.Type = appDetails.Data.Type
	app.IsFree = appDetails.Data.IsFree
	app.DLC = string(dlcString)
	app.ShortDescription = appDetails.Data.ShortDescription
	app.HeaderImage = appDetails.Data.HeaderImage
	app.Developers = string(developersString)
	app.Publishers = string(publishersString)
	app.Packages = string(packagesString)
	app.MetacriticScore = appDetails.Data.Metacritic.Score
	app.MetacriticFullURL = appDetails.Data.Metacritic.URL
	app.Categories = string(categoriesString)
	app.Genres = string(genresString)
	app.Screenshots = string(screenshotsString)
	app.Movies = string(moviesString)
	app.Achievements = string(achievementsString)
	app.Background = appDetails.Data.Background
	app.Platforms = string(platformsString)
	app.GameID = appDetails.Data.Fullgame.AppID
	app.GameName = appDetails.Data.Fullgame.Name

	return nil
}
