package datastore

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

type Donation struct {
	CreatedAt time.Time `datastore:"created_at"`
	PlayerID  int       `datastore:"player_id"`
	Amount    int       `datastore:"amount"`
	AmountUSD int       `datastore:"amount_usd"`
	Currency  string    `datastore:"currency"`
}

func (d Donation) GetKey() (key *datastore.Key) {
	return datastore.IncompleteKey(DONATION, nil)
}

func GetDonations(playerID int, limit int) (donations []Donation, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return donations, err
	}

	q := datastore.NewQuery(DONATION).Order("-created_at")

	if limit != 0 {
		q = q.Limit(limit)
	}

	if playerID != 0 {
		q = q.Filter("player_id =", playerID)
	}

	it := client.Run(ctx, q)

	for {
		var donation Donation
		_, err := it.Next(&donation)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		donations = append(donations, donation)
	}

	return donations, err
}
