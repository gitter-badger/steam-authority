package mysql

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/streadway/amqp"
)

type Package struct {
	ID          int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"` //
	CreatedAt   *time.Time `gorm:"not null;column:created_at"`                    //
	UpdatedAt   *time.Time `gorm:"not null;column:updated_at"`                    //
	Name        string     `gorm:"not null;column:name"`                          //
	BillingType int8       `gorm:"not null;column:billing_type"`                  //
	LicenseType int8       `gorm:"not null;column:license_type"`                  //
	Status      int8       `gorm:"not null;column:status"`                        //
	Apps        string     `gorm:"not null;column:apps"`                          // JSON
	ChangeID    int        `gorm:"not null;column:change_id"`                     //
}

func (pack Package) GetPath() string {
	return "/packages/" + strconv.Itoa(int(pack.ID)) + "/" + slug.Make(pack.Name)
}

func (pack Package) GetName() (name string) {

	if pack.Name == "" {
		pack.Name = "Package " + strconv.Itoa(pack.ID)
	}

	return pack.Name
}

func (pack Package) GetBillingType() (string) {

	switch pack.BillingType {
	case 11:
		return "Repurchaseable"
	default:
		return "Unknown"
	}
}

func (pack Package) GetLicenseType() (string) {

	switch pack.LicenseType {
	case 0:
		return "No License"
	default:
		return "Unknown"
	}
}

func (pack Package) GetStatus() (string) {

	switch pack.LicenseType {
	case 0:
		return "Available"
	default:
		return "Unknown"
	}
}

func (pack Package) GetApps() (apps []int, err error) {

	bytes := []byte(pack.Apps)
	if err := json.Unmarshal(bytes, apps); err != nil {
		return apps, err
	}

	return apps, nil
}

func GetPackage(id int) (pack Package, err error) {

	db, err := getDB()
	if err != nil {
		return pack, err
	}

	db.First(&pack, id)
	if db.Error != nil {
		return pack, db.Error
	}

	return pack, nil
}

func GetPackages(ids []int, columns []string) (packages []Package, err error) {

	if len(ids) < 1 {
		return packages, nil
	}

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	db.Where("id IN (?)", ids).Find(&packages)
	if db.Error != nil {
		return packages, db.Error
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
		return packages, db.Error
	}

	return packages, nil
}

func GetPackagesAppIsIn(appID int) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	db = db.Where("JSON_CONTAINS(apps, '[\"?\"]')", appID).Limit(96).Order("id DESC").Find(&packages)
	if db.Error != nil {
		return packages, db.Error
	}

	return packages, nil
}

func ConsumePackage(msg amqp.Delivery) (err error) {

	id := string(msg.Body)
	idx, _ := strconv.Atoi(id)

	db, err := getDB()
	if err != nil {
		return err
	}

	pack := new(Package)

	db.Where(Package{ID: idx}).FirstOrInit(pack)

	pack.fill()

	db.Save(&pack)
	if db.Error != nil {
		return db.Error
	}

	return err
}

// GORM callback
func (pack *Package) fill() (err error) {

	// Get app details
	err = pack.fillFromPICS()
	if err != nil {
		return err
	}

	return nil
}

func (pack *Package) fillFromPICS() (err error) {

	// Call PICS
	resp, err := steam.GetPICSInfo([]int{}, []int{pack.ID})
	if err != nil {
		return err
	}

	var pics steam.JsPackage
	if val, ok := resp.Packages[strconv.Itoa(pack.ID)]; ok {
		pics = val
	} else {
		return errors.New("no package key in json")
	}

	// Apps
	appsString, err := json.Marshal(pics.AppIDs)
	if err != nil {
		return err
	}

	pack.ID = pics.PackageID
	pack.Apps = string(appsString)
	pack.BillingType = pics.BillingType
	pack.LicenseType = pics.LicenseType
	pack.Status = pics.Status

	return nil
}
