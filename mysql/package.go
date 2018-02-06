package mysql

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
)

type Package struct {
	ID          uint      `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"name"`
	BillingType int8      `db:"billing_type"`
	LicenseType int8      `db:"license_type"`
	Status      int8      `db:"status"`
	Apps        string    `db:"apps"` // JSON
	ChangeID    int       `db:"change_id"`
}

func (packagex Package) GetPath() string {
	return "/packages/" + strconv.Itoa(int(packagex.ID)) + "/" + slug.Make(packagex.Name)
}

func GetPackage(id uint) (packagex Package, err error) {

	db, err := getDB()
	if err != nil {
		return packagex, err
	}

	err = db.Get(packagex, "SELECT * FROM apps WHERE id = $1;", id)
	if err != nil {
		return packagex, err
	}

	return packagex, nil
}

func GetPackages(ids []uint) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	// Build query
	query, args, err := squirrel.Select("*").From("packages").Where(squirrel.Eq{"id": ids}).ToSql()

	// Query
	err = db.Select(packages, query, args...)
	if err != nil {
		return packages, err
	}

	return packages, nil
}

func (packagex Package) GetApps() (apps []uint, err error) {

	bytes := []byte(packagex.Apps)
	if err := json.Unmarshal(bytes, apps); err != nil {
		return apps, err
	}

	return apps, nil
}

func (packagex *Package) Tidy() *Package {

	packagex.UpdatedAt = time.Now()
	if packagex.CreatedAt.IsZero() {
		packagex.CreatedAt = time.Now()
	}

	return packagex
}

func GetPackagesAppIsIn(appID int) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	err = db.Select(&packages, "SELECT * FROM packages WHERE JSON_CONTAINS(apps, '[\""+strconv.Itoa(appID)+"\"]')")
	if err != nil {
		return packages, err
	}

	return packages, nil
}

func GetLatestPackages() (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	err = db.Select(packages, "SELECT * FROM packages ORDER BY created_at LIMIT 20")
	if err != nil {
		return packages, err
	}

	return packages, nil
}

//func BulkAddPackages(changes []*Package) (err error) {
//
//	packagesLen := len(changes)
//	if packagesLen == 0 {
//		return nil
//	}
//
//	client, context, err := getDSClient()
//	if err != nil {
//		return err
//	}
//
//	keys := make([]*datastore.Key, 0, packagesLen)
//
//	for _, v := range changes {
//		keys = append(keys, v.GetKey())
//	}
//
//	fmt.Println("Saving " + strconv.Itoa(packagesLen) + " packages")
//
//	_, err = client.PutMulti(context, keys, changes)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
