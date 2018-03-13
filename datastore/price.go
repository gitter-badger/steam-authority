package datastore

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

type AppPrice struct {
	CreatedAt    time.Time `datastore:"created_at"`
	AppID        int       `datastore:"app_id"`
	PriceInitial int       `datastore:"price_initial"`
	PriceFinal   int       `datastore:"price_final"`
	Discount     int       `datastore:"discount"`
	Currency     string    `datastore:"currency"`
	Change       int       `datastore:"change"`
}

func (price AppPrice) GetKey() (key *datastore.Key) {
	return datastore.IncompleteKey(KindPriceApp, nil)
}

type PackagePrice struct {
	CreatedAt    time.Time `datastore:"created_at"`
	PackageID    int       `datastore:"package_id"`
	PriceInitial int       `datastore:"price_initial"`
	PriceFinal   int       `datastore:"price_final"`
	Discount     int       `datastore:"discount"`
	Currency     string    `datastore:"currency"`
	Change       int       `datastore:"change"`
}

func (pack PackagePrice) GetKey() (key *datastore.Key) {
	return datastore.IncompleteKey(KindPricePackage, nil)
}

func GetAppPrices(appID int) (prices []AppPrice, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return prices, err
	}

	q := datastore.NewQuery(KindPriceApp).Order("created_at").Limit(500)
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

func GetPackagePrices(packageID int) (prices []PackagePrice, err error) {

	client, ctx, err := getDSClient()
	if err != nil {
		return prices, err
	}

	q := datastore.NewQuery(KindPricePackage).Order("created_at").Limit(500)
	q = q.Filter("package_id =", packageID)

	it := client.Run(ctx, q)

	for {
		var price PackagePrice
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
