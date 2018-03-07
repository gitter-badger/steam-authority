package datastore

import (
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

type Price struct {
	CreatedAt time.Time `datastore:"created_at"`
	AppID     int       `datastore:"app_id"`
	Price     string    `datastore:"price"`
	Currency  string    `datastore:"currency"`
}

func (price Price) GetKey() (key *datastore.Key) {
	return datastore.IncompleteKey(PRICE, nil)
}

func CreateLogin(playerID int, r *http.Request) (err error) {

	login := new(Login)
	login.CreatedAt = time.Now()
	login.PlayerID = playerID
	login.UserAgent = r.Header.Get("User-Agent")
	login.IP = r.RemoteAddr

	_, err = SaveKind(login.GetKey(), login)

	return err
}

func GetLogins(playerID int, limit int) (logins []Login, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return logins, err
	}

	q := datastore.NewQuery(LOGIN).Order("-created_at").Limit(limit)
	q = q.Filter("player_id =", playerID)

	it := client.Run(ctx, q)

	for {
		var price Price
		_, err := it.Next(&login)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		logins = append(logins, login)
	}

	return logins, err
}
