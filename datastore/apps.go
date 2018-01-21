package datastore

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetMultiAppsByKey(keys []int) (apps []DsApp, err error) {

	client, context := getDSClient()

	keysReal := []*datastore.Key{}
	for _, v := range keys {
		keysReal = append(keysReal, datastore.NameKey(APP, strconv.Itoa(v), nil))
	}

	apps = make([]DsApp, len(keys), len(keys))
	err = client.GetMulti(context, keysReal, apps)

	return apps, err
}

func GetApp(id string) (app DsApp, err error) {

	client, context := getDSClient()

	key := datastore.NameKey(APP, id, nil)
	err = client.Get(context, key, &app)

	return app, err
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
