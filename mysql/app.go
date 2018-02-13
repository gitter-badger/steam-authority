package mysql

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type App struct {
	ID                int    `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt         int    `gorm:"not null;column:created_at"`
	UpdatedAt         int    `gorm:"not null;column:updated_at"`
	Name              string `gorm:"not null;column:name"`
	Type              string `gorm:"not null;column:type"`
	IsFree            bool   `gorm:"not null;column:is_free"`
	DLC               string `gorm:"not null;column:dlc"` // JSON
	ShortDescription  string `gorm:"not null;column:description_short"`
	HeaderImage       string `gorm:"not null;column:image_header"`
	Developers        string `gorm:"not null;column:developer"` // JSON
	Publishers        string `gorm:"not null;column:publisher"` // JSON
	Packages          string `gorm:"not null;column:packages"`  // JSON
	MetacriticScore   int8   `gorm:"not null;column:metacritic_score"`
	MetacriticFullURL string `gorm:"not null;column:metacritic_url"`
	Categories        string `gorm:"not null;column:categories"`  // JSON
	Genres            string `gorm:"not null;column:genres"`      // JSON
	Screenshots       string `gorm:"not null;column:screenshots"` // JSON
	Movies            string `gorm:"not null;column:movies"`
	Achievements      string `gorm:"not null;column:achievements"`
	Background        string `gorm:"not null;column:background"`
	Platforms         string `gorm:"not null;column:platforms"`    // JSON
	ReleaseState      string `gorm:"not null;column:releasestate"` // These down are not in the api call, only in pics
	StoreTags         []int  `gorm:"not null;column:store_tags"`
	Homepage          string `gorm:"not null;column:homepage"`
	ChangeNumber      int    `gorm:"not null;column:change_number"`
	Logo              string `gorm:"not null;column:logo"`
	Icon              string `gorm:"not null;column:icon"`
}

func (app App) GetPath() (ret string) {
	ret = "/apps/" + strconv.Itoa(int(app.ID))

	if app.Name != "" {
		ret = ret + "/" + slug.Make(app.Name)
	}

	return ret
}

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

//func GetApp(id int) (app App, err error) {
//
//	db, err := getDB()
//	if err != nil {
//		return app, err
//	}
//
//	err = db.Get(&app, "SELECT * FROM apps WHERE id = ?", id)
//	if err != nil {
//		return app, err
//	}
//
//
//
//	return app, nil
//}

func GetApp() (app App, err error) {

	db, err := getDB()
	if err != nil {
		return app, err
	}

	app = App{}

	db.First(&app, 10)

	if int64(app.UpdatedAt) < (time.Now().Unix() - int64(time.Hour*24)) {
	}

	// Don't bother checking steam to see if it exists, we should know about all apps.

	return app, nil

}

func GetApps(ids []int) (apps []App, err error) {

	db, err := getDB()
	if err != nil {
		return apps, err
	}

	db.Where("id IN (?)", ids).Find(&apps)
	if db.Error != nil {
		return apps, err
	}

	return apps, nil
}

func SearchApps(query url.Values) (apps []App, err error) {

	db, err := getDB()
	if err != nil {
		return apps, err
	}

	db = db.Limit(96).Order("id DESC") // todo, order by popularity?

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

func (app *App) Save() (err error) {

	// Tidy
	now := int(time.Now().Unix())
	app.UpdatedAt = now
	if app.CreatedAt == 0 {
		app.CreatedAt = now
	}

	if app.Developers == "" {
		app.Developers = "[]"
	}

	if app.Publishers == "" {
		app.Publishers = "[]"
	}

	if app.Categories == "" {
		app.Categories = "[]"
	}

	if app.Genres == "" {
		app.Genres = "[]"
	}

	if app.Screenshots == "" {
		app.Screenshots = "[]"
	}

	if app.Movies == "" {
		app.Movies = "[]"
	}

	if app.Achievements == "" {
		app.Achievements = "{}"
	}

	if app.Platforms == "" {
		app.Platforms = "[]"
	}

	if app.DLC == "" {
		app.DLC = "[]"
	}

	if app.Packages == "" {
		app.Packages = "[]"
	}

	app.Type = strings.ToLower(app.Type)
	app.ReleaseState = strings.ToLower(app.ReleaseState)

	// Get app details
	err = app.FillFromAppDetails()
	if err != nil {
		return err
	}

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

//func (app *App) Save() (err error) {
//
//
//
//	// Get values from struct
//	var fields []string
//	var values []interface{}
//
//	v := reflect.ValueOf(app).Elem()
//	t := reflect.TypeOf(app).Elem()
//
//	for i := 0; i < v.NumField(); i++ {
//
//		tag := t.Field(i).Tag.Get("db")
//		if tag != "" {
//			fields = append(fields, tag)
//			values = append(values, v.Field(i).Interface())
//		}
//	}
//
//	// Make SQL query
//	exists, err := app.ExistsInSQL()
//	if err != nil {
//		logger.Error(err)
//	}
//
//	var sqlString string
//	var args []interface{}
//
//	if exists {
//		sqlString, args, err = squirrel.Update("apps").ToSql()
//	} else {
//		sqlString, args, err = squirrel.Insert("apps").Columns(fields...).Values(values...).ToSql()
//	}
//
//	if err != nil {
//		return err
//	}
//
//	// Save
//	db, err := getDB()
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Query(sqlString, args...)
//	if err != nil {
//		fmt.Println(sqlString)
//		return err
//	}
//
//	return nil
//}

//func CreateApp(id int) (app App, err error) {
//
//	app.ID = id
//
//
//
//	// Save
//	err = app.Save()
//	if err != nil {
//		return app, err
//	}
//
//	return app, nil
//}

//func CreateOrUpdateApp(id int) (app App, err error) {
//
//	app.ID = 9
//	err = app.Save()
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println("xx")
//
//	//id = 9
//	//
//	//app, err = GetApp(id)
//	//if err != nil {
//	//	if err.Error() == "expired" {
//	//
//	//		app, err = GetApp(id)
//	//
//	//	} else if err.Error() == "sql: no rows in result set" {
//	//
//	//		app := App{}
//	//		app.ID = id
//	//		err = app.FillFromAppDetails()
//	//		if err != nil {
//	//			logger.Error(err)
//	//		}
//	//
//	//	}
//	//	logger.Error(err)
//	//}
//
//	return app, err
//
//}

//func (app *App) ExistsInSQL() (exists bool, err error) {
//
//	sqlString, args, err := squirrel.Select("id").Where("id = ?", app.ID).ToSql()
//	if err != nil {
//		return false, err
//	}
//
//	// Save
//	db, err := getDB()
//	if err != nil {
//		return false, err
//	}
//
//	_, err = db.Query(sqlString, args...)
//	if err != nil {
//		fmt.Println(sqlString)
//		return false, err
//	}
//
//	return true, nil
//}

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
		genre, _ := strconv.ParseInt(v.ID, 10, 8)
		genres = append(genres, int8(genre))
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

	return nil
}
