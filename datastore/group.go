package datastore

import (
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type Group struct {
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`
	GroupID   int       `datastore:"group_id"`
	Name      int       `datastore:"name"`
}

func (g Group) GetKey() (key *datastore.Key) {
	return datastore.NameKey(PLAYER, strconv.Itoa(g.GroupID), nil)
}

func GetGroupsByIDs(ids []int) (groups []Group, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return groups, err
	}

	var keys []*datastore.Key
	for _, v := range ids {
		key := datastore.NameKey(PLAYER, strconv.Itoa(v), nil)
		keys = append(keys, key)
	}

	groups = make([]Group, len(keys))
	err = client.GetMulti(context, keys, groups)
	if err != nil && !strings.Contains(err.Error(), "no such entity") {
		return groups, err
	}

	return groups, nil
}
