package datastore

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

type AppPrice struct {
	CreatedAt time.Time `datastore:"created_at"`
	AppID     int       `datastore:"app_id"`
	Price     int       `datastore:"price"`
	Discount  int       `datastore:"discount"`
	Currency  string    `datastore:"currency"`
}

func (price AppPrice) GetKey() (key *datastore.Key) {
	return datastore.IncompleteKey(PRICE, nil)
}

func CreatePrice(appID int, price int, discount int) (err error) {

	p := new(AppPrice)
	p.CreatedAt = time.Now()
	p.AppID = appID
	p.Price = price
	p.Discount = discount
	p.Currency = "usd"

	_, err = SaveKind(p.GetKey(), p)

	return err
}

func GetPrices(appID int) (prices []AppPrice, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return prices, err
	}

	q := datastore.NewQuery(PRICE).Order("created_at").Limit(500)
	q = q.Filter("app_id =", appID)

	it := client.Run(ctx, q)

	for {
		var price AppPrice
		_, err := it.Next(&price)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		prices = append(prices, price)
	}

	return prices, err
}
