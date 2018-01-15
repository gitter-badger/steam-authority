package main

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

func saveApp(data dsApp) {

	key := datastore.NameKey(
		"App",
		data.AppID,
		nil,
	)

	saveKind(key, &data)
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
}

type dsPackage struct {
	PackageID string  `datastore:"package_id"`
	Apps      []dsApp `datastore:"apps"`
}
