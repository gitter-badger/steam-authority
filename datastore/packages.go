package datastore

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
	"google.golang.org/api/iterator"
)

func GetPackagesAppIsIn(appID int) (packages []DsPackage, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return packages, err
	}

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

	client, context, err := getDSClient()
	if err != nil {
		return nil, err
	}

	var keysReal []*datastore.Key
	for _, v := range keys {
		keysReal = append(keysReal, datastore.NameKey(PACKAGE, strconv.Itoa(v), nil))
	}

	packages = make([]DsPackage, len(keys), len(keys))

	err = client.GetMulti(context, keysReal, packages)
	if err != nil {
		return packages, err
	}

	return packages, nil
}

func GetLatestUpdatedPackages(limit int) (packages []DsPackage, err error) {

	client, context, err := getDSClient()
	if err != nil {
		logger.Error(err)
		return packages, err
	}

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

func GetPackage(id string) (packagex DsPackage, err error) {

	client, context, err := getDSClient()
	if err != nil {
		return packagex, err
	}

	key := datastore.NameKey(PACKAGE, id, nil)

	err = client.Get(context, key, &packagex)
	if err != nil {
		return packagex, err
	}

	return packagex, nil
}

func BulkAddPackages(changes []*DsPackage) (err error) {

	packagesLen := len(changes)
	if packagesLen == 0 {
		return nil
	}

	client, context, err := getDSClient()
	if err != nil {
		return err
	}

	keys := make([]*datastore.Key, 0, packagesLen)

	for _, v := range changes {
		keys = append(keys, v.GetKey())
	}

	fmt.Println("Saving " + strconv.Itoa(packagesLen) + " packages")

	_, err = client.PutMulti(context, keys, changes)
	if err != nil {
		return err
	}

	return nil
}
