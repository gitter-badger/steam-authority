package datastore

import (
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

type DsChange struct {
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`
	ChangeID  int       `datastore:"change_id"`
	Apps      []int     `datastore:"apps"`
	Packages  []int     `datastore:"packages"`
}

func (change DsChange) GetKey() (key *datastore.Key) {
	return datastore.NameKey(CHANGE, strconv.Itoa(change.ChangeID), nil)
}

func (change *DsChange) Tidy() *DsChange {

	change.UpdatedAt = time.Now()
	if change.CreatedAt.IsZero() {
		change.CreatedAt = time.Now()
	}

	return change
}


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
