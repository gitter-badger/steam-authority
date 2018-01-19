package datastore

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

const (
	// CHANGE datastore kind
	CHANGE = "Change"

	// APP atastore kind
	APP = "App"

	// PACKAGE atastore kind
	PACKAGE = "Package"
)

func GetLatestChanges(limit int) (changes []DsChange, err error) {

	client, context := getDSClient()

	q := datastore.NewQuery(CHANGE).Order("-change_id").Limit(limit)
	it := client.Run(context, q)

	for {
		var change DsChange
		_, err := it.Next(&change)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		changes = append(changes, change)
	}

	return changes, err
}

func GetLatestUpdatedApps(limit int) (apps []DsApp, err error) {

	client, context := getDSClient()

	q := datastore.NewQuery(APP).Order("-change_number").Limit(limit)
	it := client.Run(context, q)

	for {
		var app DsApp
		_, err := it.Next(&app)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		apps = append(apps, app)
	}

	return apps, err
}

func GetChange(id string) (change *DsChange, err error) {

	client, context := getDSClient()

	key := datastore.NameKey(CHANGE, id, nil)

	err = client.Get(context, key, change)
	if err != nil {
		logger.Error(err)
	}

	return change, err
}

func GetApp(id string) (app DsApp, err error) {

	client, context := getDSClient()

	key := datastore.NameKey(APP, id, nil)

	err = client.Get(context, key, app)
	if err != nil {
		logger.Error(err)
	}

	return app, err
}

func BulkAddChanges(changes []*DsChange) (err error) {

	len := len(changes)
	if len == 0 {
		return nil
	}

	client, context := getDSClient()
	keys := make([]*datastore.Key, 0, len)

	for _, v := range changes {
		keys = append(keys, datastore.NameKey(CHANGE, strconv.Itoa(v.ChangeID), nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len) + " changes")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}

func BulkAddApps(changes []*DsApp) (err error) {

	len := len(changes)
	if len == 0 {
		return nil
	}

	client, context := getDSClient()
	keys := make([]*datastore.Key, 0, len)

	for _, v := range changes {
		keys = append(keys, datastore.NameKey(APP, strconv.Itoa(v.AppID), nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len) + " apps")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}

func BulkAddPackages(changes []*DsPackage) (err error) {

	len := len(changes)
	if len == 0 {
		return nil
	}

	client, context := getDSClient()
	keys := make([]*datastore.Key, 0, len)

	for _, v := range changes {
		keys = append(keys, datastore.NameKey(PACKAGE, strconv.Itoa(v.PackageID), nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len) + " packages")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}

func savePackage(data DsPackage) {

	packageIDString := strconv.Itoa(data.PackageID)

	key := datastore.NameKey(PACKAGE, packageIDString, nil)

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
