package datastore

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/websockets"
	"github.com/streadway/amqp"
	"google.golang.org/api/iterator"
)

type Change struct {
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`
	ChangeID  int       `datastore:"change_id"`
	Apps      []int     `datastore:"apps"`
	Packages  []int     `datastore:"packages"`
}

func (change Change) GetKey() (key *datastore.Key) {
	return datastore.NameKey(CHANGE, strconv.Itoa(change.ChangeID), nil)
}

func (change *Change) Tidy() *Change {

	change.UpdatedAt = time.Now()
	if change.CreatedAt.IsZero() {
		change.CreatedAt = time.Now()
	}

	return change
}

func GetLatestChanges(limit int) (changes []Change, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return changes, err
	}

	q := datastore.NewQuery(CHANGE).Order("-change_id").Limit(limit)
	it := client.Run(context, q)

	for {
		var change Change
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

func GetChange(id string) (change *Change, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return change, err
	}

	key := datastore.NameKey(CHANGE, id, nil)

	change = &Change{}
	err = client.Get(context, key, change)
	if err != nil {
		logger.Error(err)
	}

	return change, nil
}

func AddChanges(changes []*Change) (err error) {

	changesLen := len(changes)
	if changesLen == 0 {
		return nil
	}

	client, context, err := getDSClient()
	if err != nil {
		return err
	}

	keys := make([]*datastore.Key, 0)

	for _, v := range changes {
		keys = append(keys, v.GetKey())
	}

	fmt.Println("Saving " + strconv.Itoa(changesLen) + " changes")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}

func ConsumeChange(msg amqp.Delivery) (err error) {

	var change Change
	if err := json.Unmarshal(msg.Body, &change); err != nil {
		return err
	}

	logger.Info("Reading change " + strconv.Itoa(change.ChangeID) + " from rabbit")

	// Save to DS
	err = AddChanges([]*Change{&change})
	if err != nil {
		return err
	}

	// Send websocket
	payload := changeWebsocketPayload{
		ID:        change.ChangeID,
		CreatedAt: change.CreatedAt.Unix(),
		Apps:      change.Apps,
		Packages:  change.Packages,
	}
	websockets.Send(websockets.CHANGES, payload)

	return nil
}

type changeWebsocketPayload struct {
	ID        int   `json:"id"`
	CreatedAt int64 `json:"created_at"`
	Apps      []int `json:"apps"`
	Packages  []int `json:"packages"`
}
