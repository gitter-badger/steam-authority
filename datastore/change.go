package datastore

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetLatestChanges(limit int) (changes []DsChange, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return changes, err
	}

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

func GetChange(id string) (change *DsChange, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return change, err
	}

	key := datastore.NameKey(CHANGE, id, nil)

	change = &DsChange{}
	err = client.Get(context, key, change)
	if err != nil {
		logger.Error(err)
	}

	return change, nil
}

func BulkAddChanges(changes []*DsChange) (err error) {

	changesLen := len(changes)
	if changesLen == 0 {
		return nil
	}

	client, context, err := getDSClient()
	if err != nil {
		return err
	}

	keys := make([]*datastore.Key, 0, changesLen)

	for _, v := range changes {
		keys = append(keys, v.GetKey(), nil)
	}

	fmt.Println("Saving " + strconv.Itoa(changesLen) + " changes")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}
