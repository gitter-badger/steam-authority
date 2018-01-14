package main

import (
	"context"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

func saveChange(data dsChange) {

	key := datastore.NameKey(
		"Change",
		strconv.Itoa(data.ChangeID)+"-"+strconv.Itoa(data.LatestChangeID),
		nil,
	)

	saveKind(key, &data)
}

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

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("STEAM_GOOGLE_PROJECT"))
	if err != nil {
		logger.Error(err)
	}

	return client, ctx
}

type dsChange struct {
	ChangeID       int      `datastore:"change_id"`
	LatestChangeID int      `datastore:"latest_change_id"`
	Apps           []string `datastore:"apps"`
	Packages       []string `datastore:"packages"`
}

type dsApp struct {
	AppID string `datastore:"app_id"`
	Test  string `datastore:"test"`
}

type dsPackage struct {
	PackageID string  `datastore:"package_id"`
	Apps      []dsApp `datastore:"apps"`
}
