package main

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

func saveApp(jsApp JsApp) {

	jsTags := jsApp.Common.StoreTags
	tags := make([]string, 0, len(jsTags))
	for _, value := range jsTags {
		tags = append(tags, value)
	}

	dsApp := dsApp{}
	dsApp.AppID = jsApp.AppID
	dsApp.Name = jsApp.Common.Name
	dsApp.Type = jsApp.Common.Type
	dsApp.ReleaseState = jsApp.Common.ReleaseState
	dsApp.OSList = strings.Split(jsApp.Common.OSList, ",")
	dsApp.MetacriticScore = jsApp.Common.MetacriticScore
	dsApp.MetacriticFullURL = jsApp.Common.MetacriticURL
	dsApp.StoreTags = tags
	dsApp.Developer = jsApp.Extended.Developer
	dsApp.Publisher = jsApp.Extended.Publisher
	dsApp.Homepage = jsApp.Extended.Homepage
	dsApp.ChangeNumber = jsApp.ChangeNumber

	key := datastore.NameKey(
		"App",
		dsApp.AppID,
		nil,
	)

	saveKind(key, &dsApp)
}

func savePackage(data dsPackage) {

	key := datastore.NameKey(
		"Package",
		data.PackageID,
		nil,
	)

	saveKind(key, &data)
}

func saveKind(key *datastore.Key, data interface{}) (newKey *datastore.Key) {

	client, context := getDSClient()
	newKey, err := client.Put(context, key, data)
	if err != nil {
		logger.Error(err)
	}

	return newKey
}

func getDSClient() (*datastore.Client, context.Context) {

	context := context.Background()
	client, err := datastore.NewClient(
		context,
		os.Getenv("STEAM_GOOGLE_PROJECT"),
	)
	if err != nil {
		logger.Error(err)
	}

	return client, context
}

type dsChange struct {
	ChangeID int      `datastore:"change_id"`
	Apps     []string `datastore:"apps"`
	Packages []string `datastore:"packages"`
}

type dsApp struct {
	AppID             string   `datastore:"app_id"`
	Name              string   `datastore:"name"`
	Type              string   `datastore:"type"`
	ReleaseState      string   `datastore:"releasestate"`
	OSList            []string `datastore:"oslist"`
	MetacriticScore   string   `datastore:"metacritic_score"`
	MetacriticFullURL string   `datastore:"metacritic_fullurl"`
	StoreTags         []string `datastore:"store_tags"`
	Developer         string   `datastore:"developer"`
	Publisher         string   `datastore:"publisher"`
	Homepage          string   `datastore:"homepage"`
	ChangeNumber      int      `datastore:"change_number"`
}

type dsPackage struct {
	PackageID string  `datastore:"package_id"`
	Apps      []dsApp `datastore:"apps"`
}
