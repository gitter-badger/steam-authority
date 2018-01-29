package datastore

import (
	"fmt"
	"net/url"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetMultiAppsByKey(keys []int) (apps []DsApp, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return apps, err
	}

	var keysReal []*datastore.Key
	for _, v := range keys {
		keysReal = append(keysReal, datastore.NameKey(APP, strconv.Itoa(v), nil))
	}

	apps = make([]DsApp, len(keys), len(keys))
	err = client.GetMulti(context, keysReal, apps)
	if err != nil {
		return apps, err
	}

	return apps, nil
}

func GetApp(id string) (app DsApp, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return app, err
	}

	key := datastore.NameKey(APP, id, nil)
	err = client.Get(context, key, &app)
	if err != nil {
		return app, err
	}

	return app, nil
}

func SearchApps(query url.Values, limit int) (apps []DsApp, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return apps, err
	}

	q := datastore.NewQuery(APP).Limit(limit)

	addFilter(q, query, "os", "oslist")
	addFilter(q, query, "tag", "store_tags")

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

func addFilter(q *datastore.Query, query url.Values, formName string, dbName string) *datastore.Query {

	formValue := query.Get(formName)
	if formValue != "" {
		q = q.Filter(dbName+" =", formValue)
	}

	return q
}

func BulkAddApps(changes []*DsApp) (err error) {

	appsLen := len(changes)
	if appsLen == 0 {
		return nil
	}

	client, context, err := getDSClient()
	if err != nil {
		return err
	}

	keys := make([]*datastore.Key, 0, appsLen)

	for _, v := range changes {
		keys = append(keys, v.GetKey())
	}

	fmt.Println("Saving " + strconv.Itoa(appsLen) + " apps")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}

func CountApps() (count int, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return count, err
	}

	q := datastore.NewQuery(APP)
	count, err = client.Count(ctx, q)
	if err != nil {
		return count, err
	}

	return count, nil
}
