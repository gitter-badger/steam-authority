package mysql

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
)

type App struct {
	ID                int    `db:"id"`
	CreatedAt         int    `db:"created_at"`
	UpdatedAt         int    `db:"updated_at"`
	Name              string `db:"name"`
	Type              string `db:"type"`
	IsFree            bool   `db:"is_free"`
	DLC               string `db:"dlc"` // []int
	ShortDescription  string `db:"description_short"`
	HeaderImage       string `db:"image_header"`
	Developers        string `db:"developer"` // []string
	Publishers        string `db:"publisher"` // []string
	Packages          string `db:"packages"`  // []int
	MetacriticScore   int8   `db:"metacritic_score"`
	MetacriticFullURL string `db:"metacritic_url"`
	Categories        string `db:"categories"`  // []int8
	Genres            string `db:"genres"`      // []int8
	Screenshots       string `db:"screenshots"` // []
	Movies            string `db:"movies"`
	Achievements      string `db:"achievements"`
	Background        string `db:"background"`
	Platforms         string `db:"platforms"` // []string

	// These are not in the api call, only in pics
	ReleaseState string `dbx:"releasestate"`
	StoreTags    []int  `dbx:"store_tags"`
	Homepage     string `dbx:"homepage"`
	ChangeNumber int    `dbx:"change_number"`
	Logo         string `dbx:"logo"`
	Icon         string `dbx:"icon"`
}

func (app App) GetPath() (ret string) {
	ret = "/apps/" + strconv.Itoa(int(app.ID))

	if app.Name != "" {
		ret = ret + "/" + slug.Make(app.Name)
	}

	return ret
}

func GetApp(id uint) (app App, err error) {

	db, err := getDB()
	if err != nil {
		return app, err
	}

	err = db.Get(&app, "SELECT * FROM apps WHERE id = ?", id)
	if err != nil {
		return app, err
	}

	app.tidyJSON()

	return app, nil
}

func GetApps(ids []uint) (apps []App, err error) {

	db, err := getDB()
	if err != nil {
		return apps, err
	}

	// Build query
	query, args, err := squirrel.Select("*").From("apps").Where(squirrel.Eq{"id": ids}).ToSql()

	// Query
	err = db.Select(apps, query, args...)
	if err != nil {
		return apps, err
	}

	return apps, nil
}

func (app App) GetScreenshots() (screenshots []steam.AppDetailsScreenshot, err error) {

	fmt.Println("xx")
	fmt.Println(app.Screenshots)
	fmt.Println("xx")

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

func (app App) GetName() (name string) {

	if app.Name == "" {
		app.Name = "App " + strconv.Itoa(int(app.ID))
	}

	return app.Name
}

func (app *App) tidyJSON() {

	if app.DLC == "" {
		app.DLC = "[]"
	}

	if app.Packages == "" {
		app.Packages = "[]"
	}

	if app.Categories == "" {
		app.Categories = "[]"
	}

	if app.Genres == "" {
		app.Genres = "[]"
	}

	if app.Platforms == "" {
		app.Platforms = "[]"
	}

	if app.Screenshots == "" {
		app.Screenshots = "[]"
	}

	if app.Achievements == "" {
		app.Achievements = "{}"
	}
}

// todo, use args in JSON function queries when its fixed in sqlx
func SearchApps(query url.Values) (apps []App, err error) {

	searchQuery := squirrel.Select("*").From("apps").Limit(96).OrderBy("id DESC") // todo, order by popularity

	// Platforms
	if _, ok := query["platforms"]; ok {
		searchQuery = searchQuery.Where("JSON_CONTAINS(platforms, '[\"" + query.Get("platforms") + "\"]')")
	}

	// Tag
	if _, ok := query["tags"]; ok {
		searchQuery = searchQuery.Where("JSON_CONTAINS(tags, '[\"" + query.Get("tags") + "\"]')")
	}

	// Query
	db, err := getDB()
	if err != nil {
		return apps, err
	}

	sql, args, err := searchQuery.ToSql()

	err = db.Select(&apps, sql, args...)
	if err != nil {
		return apps, err
	}

	return apps, err
}

func (app App) Save(id int) (err error) {

	// Tidy
	app.ID = id

	now := int(time.Now().Unix())
	app.UpdatedAt = now
	if app.CreatedAt == 0 {
		app.CreatedAt = now
	}

	app.tidyJSON()

	app.Type = strings.ToLower(app.Type)
	app.ReleaseState = strings.ToLower(app.ReleaseState)

	// Get values from struct
	var fields []string
	var values []interface{}

	v := reflect.ValueOf(app)
	t := reflect.TypeOf(app)

	for i := 0; i < v.NumField(); i++ {

		tag := t.Field(i).Tag.Get("db")
		if tag != "" {
			fields = append(fields, tag)
			values = append(values, v.Field(i).Interface())
		}
	}

	// Make SQL query
	sqlString, args, err := squirrel.Insert("apps").Columns(fields...).Values(values...).ToSql()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Save
	db, err := getDB()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sqlString)
	fmt.Println(args)

	_, err = db.Query(sqlString, args...)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) FillFromAppDetails() (ret *App, err error) {

	// Get data
	appDetails, err := steam.GetAppDetails(strconv.Itoa(int(app.ID)))
	if err != nil {
		return ret, err
	}

	// Screenshots
	screenshotsString, err := json.Marshal(appDetails.Data.Screenshots)
	if err != nil {
		return ret, err
	}

	// Movies
	moviesString, err := json.Marshal(appDetails.Data.Movies)
	if err != nil {
		return ret, err
	}

	// Achievements
	achievementsString, err := json.Marshal(appDetails.Data.Achievements)
	if err != nil {
		return ret, err
	}

	// DLC
	dlcString, err := json.Marshal(appDetails.Data.DLC)
	if err != nil {
		return ret, err
	}

	// Developers
	developersString, err := json.Marshal(appDetails.Data.Developers)
	if err != nil {
		return ret, err
	}

	// Publishers
	publishersString, err := json.Marshal(appDetails.Data.Publishers)
	if err != nil {
		return ret, err
	}

	// Packages
	packagesString, err := json.Marshal(appDetails.Data.Packages)
	if err != nil {
		return ret, err
	}

	// Categories
	var categories []int8
	for _, v := range appDetails.Data.Categories {
		categories = append(categories, v.ID)
	}

	categoriesString, err := json.Marshal(categories)
	if err != nil {
		return ret, err
	}

	// Genres
	var genres []int8
	for _, v := range appDetails.Data.Genres {
		genre, _ := strconv.ParseInt(v.ID, 10, 8)
		genres = append(genres, int8(genre))
	}

	genresString, err := json.Marshal(genres)
	if err != nil {
		return ret, err
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
		return ret, err
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

	return ret, nil
}
