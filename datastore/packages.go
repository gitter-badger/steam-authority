package datastore

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetPackagesAppIsIn(appID int) (packages []DsPackage, err error) {

	client, context := getDSClient()

	q := datastore.NewQuery(PACKAGE).Filter("apps =", appID)
	it := client.Run(context, q)

	for {
		var packagex DsPackage
		_, err := it.Next(&packagex)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
			break
		}

		packages = append(packages, packagex)
	}

	return packages, err
}

func GetMultiPackagesByKey(keys []int) (packages []DsPackage, err error) {

	client, context := getDSClient()

	keysReal := []*datastore.Key{}
	for _, v := range keys {
		keysReal = append(keysReal, datastore.NameKey(PACKAGE, strconv.Itoa(v), nil))
	}

	packages = make([]DsPackage, len(keys), len(keys))

	err = client.GetMulti(context, keysReal, packages)
	if err != nil {
		logger.Error(err)
	}

	return packages, err
}

func GetLatestUpdatedPackages(limit int) (packages []DsPackage, err error) {

	client, context := getDSClient()

	q := datastore.NewQuery(PACKAGE).Order("-change_id").Limit(limit)
	it := client.Run(context, q)

	for {
		var dsPackage DsPackage
		_, err := it.Next(&dsPackage)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		packages = append(packages, dsPackage)
	}

	return packages, err
}

func GetPackage(id string) (app DsPackage, err error) {

	client, context := getDSClient()

	key := datastore.NameKey(PACKAGE, id, nil)

	err = client.Get(context, key, &app)
	if err != nil {
		logger.Error(err)
	}

	return app, err
}

func savePackage(data DsPackage) {

	packageIDString := strconv.Itoa(data.PackageID)

	key := datastore.NameKey(PACKAGE, packageIDString, nil)

	saveKind(key, &data)
}

func BulkAddPackages(changes []*DsPackage) (err error) {

	len := len(changes)
	if len == 0 {
		return nil
	}

	client, context := getDSClient()
	keys := make([]*datastore.Key, 0, len)

	for _, v := range changes {
		keys = append(keys, datastore.NameKey(PACKAGE, strconv.Itoa(v.PackageID), nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len) + " packages")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}
