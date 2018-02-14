package mysql

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gosimple/slug"
	"strings"
)

type Package struct {
	ID          int       `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt   time.Time `gorm:"not null;column:created_at"`
	UpdatedAt   time.Time `gorm:"not null;column:updated_at"`
	Name        string    `gorm:"not null;column:name"`
	BillingType int8      `gorm:"not null;column:billing_type"`
	LicenseType int8      `gorm:"not null;column:license_type"`
	Status      int8      `gorm:"not null;column:status"`
	Apps        string    `gorm:"not null;column:apps"` // JSON
	ChangeID    int       `gorm:"not null;column:change_id"`
}

func (pack Package) GetPath() string {
	return "/packages/" + strconv.Itoa(int(pack.ID)) + "/" + slug.Make(pack.Name)
}

func GetPackage(id int) (pack Package, err error) {

	db, err := getDB()
	if err != nil {
		return pack, err
	}

	db.First(&pack, id)
	if db.Error != nil {
		return pack, err
	}

	if pack.UpdatedAt.Unix() < time.Now().AddDate(0, 0, -1).Unix() {

	}

	// Don't bother checking steam to see if it exists, we should know about all packs.

	return pack, nil
}

func GetPackages(ids []int) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	db.Where("id IN (?)", ids).Find(&packages)
	if db.Error != nil {
		return packages, err
	}

	return packages, nil
}

func (pack Package) GetApps() (apps []int, err error) {

	bytes := []byte(pack.Apps)
	if err := json.Unmarshal(bytes, apps); err != nil {
		return apps, err
	}

	return apps, nil
}

func (pack *Package) Tidy() *Package {

	pack.UpdatedAt = time.Now()
	if pack.CreatedAt.IsZero() {
		pack.CreatedAt = time.Now()
	}

	return pack
}

func GetPackagesAppIsIn(appID int) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	db = db.Where("JSON_CONTAINS(apps, '[\"?\"]')", appID).Limit(96).Order("id DESC").Find(&packages)
	if db.Error != nil {
		return packages, err
	}

	return packages, nil
}

func GetLatestPackages() (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	db.Limit(20).Order("created_at DESC").Find(&packages)
	if db.Error != nil {
		return packages, err
	}

	return packages, nil
}

func (pack *Package) FillFromPICS() (err error) {

	dsPackage := datastore.DsPackage{}
	dsPackage.PackageID = js.PackageID
	dsPackage.Apps = js.AppIDs
	dsPackage.BillingType = js.BillingType
	dsPackage.LicenseType = js.LicenseType
	dsPackage.Status = js.Status

	return &dsPackage

}