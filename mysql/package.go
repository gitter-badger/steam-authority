package mysql

import (
	"encoding/json"
	"errors"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/streadway/amqp"
)

type Package struct {
	ID              int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"` //
	CreatedAt       *time.Time `gorm:"not null;column:created_at"`                    //
	UpdatedAt       *time.Time `gorm:"not null;column:updated_at"`                    //
	Name            string     `gorm:"not null;column:name"`                          //
	ImagePage       string     `gorm:"not null;column:image_page"`                    //
	ImageHeader     string     `gorm:"not null;column:image_header"`                  //
	ImageLogo       string     `gorm:"not null;column:image_logo"`                    //
	BillingType     int8       `gorm:"not null;column:billing_type"`                  //
	LicenseType     int8       `gorm:"not null;column:license_type"`                  //
	Status          int8       `gorm:"not null;column:status"`                        //
	Apps            string     `gorm:"not null;column:apps;default:'[]'"`             // JSON
	ChangeID        int        `gorm:"not null;column:change_id"`                     //
	Extended        string     `gorm:"not null;column:extended;default:'{}'"`         // JSON
	PurchaseText    string     `gorm:"not null;column:purchase_text"`                 //
	PriceInitial    int        `gorm:"not null;column:price_initial"`                 //
	PriceFinal      int        `gorm:"not null;column:price_final"`                   //
	PriceDiscount   int        `gorm:"not null;column:price_discount"`                //
	PriceIndividual int        `gorm:"not null;column:price_individual"`              //
	Controller      string     `gorm:"not null;column:controller;default:'{}'"`       // JSON
	ComingSoon      bool       `gorm:"not null;column:coming_soon"`                   //
	ReleaseDate     *time.Time `gorm:"not null;column:controller"`                    //
	Platforms       string     `gorm:"not null;column:platforms;default:'[]'"`        // JSON
}

func getDefaultPackageJSON() Package {
	return Package{
		Apps:     "[]",
		Extended: "{}",
	}
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
	case 0:
		return "No Cost"
	case 1:
		return "Store"
	case 2:
		return "Bill Monthly"
	case 3:
		return "CD Key"
	case 4:
		return "Guest Pass"
	case 5:
		return "Hardware Promo"
	case 6:
		return "Gift"
	case 7:
		return "Free Weekend"
	case 8:
		return "OEM Ticket"
	case 9:
		return "Recurring Option"
	case 10:
		return "Store or CD Key"
	case 11:
		return "Repurchaseable"
	case 12:
		return "Free on Demand"
	case 13:
		return "Rental"
	case 14:
		return "Commercial License"
	case 15:
		return "Free Commercial License"
	default:
		return "Unknown"
	}
}

func (pack Package) GetLicenseType() (string) {

	switch pack.LicenseType {
	case 0:
		return "No License"
	case 1:
		return "Single Purchase"
	case 2:
		return "Single Purchase (Limited Use)"
	case 3:
		return "Recurring Charge"
	case 6:
		return "Recurring"
	case 7:
		return "Limited Use Delayed Activation"
	default:
		return "Unknown"
	}
}

func (pack Package) GetStatus() (string) {

	switch pack.Status {
	case 0:
		return "Available"
	case 2:
		return "Unavailable"
	default:
		return "Unknown"
	}
}

func (pack Package) GetApps() (apps []int, err error) {

	bytes := []byte(pack.Apps)
	if err := json.Unmarshal(bytes, &apps); err != nil {
		return apps, err
	}

	return apps, nil
}

func (pack Package) GetExtended() (extended map[string]interface{}, err error) {

	extended = make(map[string]interface{})

	bytes := []byte(pack.Extended)
	if err := json.Unmarshal(bytes, &extended); err != nil {
		return extended, err
	}

	return extended, nil
}

func (pack Package) GetPlatforms() (platforms []string, err error) {

	bytes := []byte(pack.Platforms)
	if err := json.Unmarshal(bytes, &platforms); err != nil {
		return platforms, err
	}

	return platforms, nil
}

func (pack Package) GetPlatformImages() (ret template.HTML, err error) {

	platforms, err := pack.GetPlatforms()
	if err != nil {
		return ret, err
	}

	for _, v := range platforms {
		if v == "macos" {
			ret = ret + `<i class="fab fa-apple"></i>`
		} else if v == "windows" {
			ret = ret + `<i class="fab fa-windows"></i>`
		} else if v == "linux" {
			ret = ret + `<i class="fab fa-linux"></i>`
		}
	}

	return ret, nil
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

	if pack.ID == 0 {
		return pack, errors.New("no id")
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

func GetLatestPackages(limit int) (packages []Package, err error) {

	db, err := getDB()
	if err != nil {
		return packages, err
	}

	db.Limit(limit).Order("created_at DESC").Find(&packages)
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

	db = db.Where("JSON_CONTAINS(apps, '[\"" + strconv.Itoa(appID) + "\"]')").Limit(96).Order("id DESC").Find(&packages)

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

	db.Attrs(getDefaultPackageJSON()).FirstOrCreate(pack, Package{ID: idx})

	pack.fill()

	db.Save(pack)
	if db.Error != nil {
		return db.Error
	}

	return err
}

// GORM callback
func (pack *Package) fill() (err error) {

	// Get app details
	err = pack.fillFromAPI()
	if err != nil {
		return err
	}

	// Get app details from PICS
	err = pack.fillFromPICS()
	if err != nil {
		return err
	}

	// Default JSON values
	if pack.Apps == "" || pack.Apps == "null" {
		pack.Apps = "[]"
	}

	if pack.Extended == "" || pack.Extended == "null" {
		pack.Extended = "{}"
	}

	return nil
}

func (pack *Package) fillFromAPI() (err error) {

	// Get data
	response, err := steam.GetPackageDetailsFromStore(pack.ID)
	if err != nil {

		// Not all packages can be found
		if err.Error() == "no package with id in steam" || strings.HasPrefix(err.Error(), "invalid package id:") {
			return nil
		}

		return err
	}

	// Controller
	controllerString, err := json.Marshal(response.Data.Controller)
	if err != nil {
		return err
	}

	// Platforms
	var platforms []string
	if response.Data.Platforms.Linux {
		platforms = append(platforms, "linux")
	}
	if response.Data.Platforms.Windows {
		platforms = append(platforms, "windows")
	}
	if response.Data.Platforms.Windows {
		platforms = append(platforms, "macos")
	}

	platformsString, err := json.Marshal(platforms)
	if err != nil {
		return err
	}

	// Release date
	var releaseDate = time.Time{}
	if response.Data.ReleaseDate.Date != "" {
		releaseDate, err = time.Parse("2 Jan, 2006", response.Data.ReleaseDate.Date)
		if err != nil {
			return err
		}
	}

	//
	pack.Name = response.Data.Name
	pack.ImageHeader = response.Data.HeaderImage
	pack.ImageLogo = response.Data.SmallLogo
	pack.ImageHeader = response.Data.HeaderImage
	//pack.Apps = string(appsString) // Can get from PICS
	pack.PriceInitial = response.Data.Price.Initial
	pack.PriceFinal = response.Data.Price.Final
	pack.PriceDiscount = response.Data.Price.DiscountPercent
	pack.PriceIndividual = response.Data.Price.Individual
	pack.Platforms = string(platformsString)
	pack.Controller = string(controllerString)
	pack.ReleaseDate = &releaseDate
	pack.ComingSoon = response.Data.ReleaseDate.ComingSoon

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

	// Extended
	extended, err := json.Marshal(pics.Extended)
	if err != nil {
		return err
	}

	pack.ID = pics.PackageID
	pack.Apps = string(appsString)
	pack.BillingType = pics.BillingType
	pack.LicenseType = pics.LicenseType
	pack.Status = pics.Status
	pack.Extended = string(extended)

	return nil
}

// todo, make these nice, put into the GetExtended func?
var PackageKeys = map[string]string{
	"allowcrossregiontradingandgifting":     "allowcrossregiontradingandgifting",
	"allowpurchasefromretrictedcountries":   "allowpurchasefromretrictedcountries",
	"allowpurchaseinrestrictedcountries":    "allowpurchaseinrestrictedcountries",
	"allowpurchaserestrictedcountries":      "allowpurchaserestrictedcountries",
	"allowrunincountries":                   "allowrunincountries",
	"alwayscountasowned":                    "alwayscountasowned",
	"alwayscountsasowned":                   "Always Counts As Owned",
	"alwayscountsasunowned":                 "alwayscountsasunowned",
	"appid":                                 "appid",
	"appidownedrequired":                    "appidownedrequired",
	"billingagreementtype":                  "billingagreementtype",
	"blah":                                  "blah",
	"canbegrantedfromexternal":              "canbegrantedfromexternal",
	"cantownapptopurchase":                  "cantownapptopurchase",
	"complimentarypackagegrant":             "complimentarypackagegrant",
	"complimentarypackagegrants":            "complimentarypackagegrants",
	"curatorconnect":                        "curatorconnect",
	"devcomp":                               "devcomp",
	"dontallowrunincountries":               "dontallowrunincountries",
	"dontgrantifappidowned":                 "dontgrantifappidowned",
	"enforceintraeeaactivationrestrictions": "enforceintraeeaactivationrestrictions",
	"excludefromsharing":                    "excludefromsharing",
	"exfgls":                                "exfgls",
	"expirytime":                            "expirytime",
	"extended":                              "extended",
	"fakechange":                            "fakechange",
	"foo":                                   "foo",
	"freeondemand":                          "freeondemand",
	"freeweekend":                           "freeweekend",
	"giftsaredeletable":                     "giftsaredeletable",
	"giftsaremarketable":                    "giftsaremarketable",
	"giftsaretradable":                      "giftsaretradable",
	"grantexpirationdays":                   "grantexpirationdays",
	"grantguestpasspackage":                 "grantguestpasspackage",
	"grantpassescount":                      "grantpassescount",
	"hardwarepromotype":                     "hardwarepromotype",
	"ignorepurchasedateforrefunds":          "ignorepurchasedateforrefunds",
	"initialperiod":                         "initialperiod",
	"initialtimeunit":                       "initialtimeunit",
	"iploginrestriction":                    "iploginrestriction",
	"languages":                             "languages",
	"launcheula":                            "launcheula",
	"legacygamekeyappid":                    "legacygamekeyappid",
	"lowviolenceinrestrictedcountries":      "lowviolenceinrestrictedcountries",
	"martinotest":                           "martinotest",
	"mustownapptopurchase":                  "mustownapptopurchase",
	"onactivateguestpassmsg":                "onactivateguestpassmsg",
	"onexpiredmsg":                          "onexpiredmsg",
	"ongrantguestpassmsg":                   "ongrantguestpassmsg",
	"onlyallowincountries":                  "onlyallowincountries",
	"onlyallowrestrictedcountries":          "onlyallowrestrictedcountries",
	"onlyallowrunincountries":               "onlyallowrunincountries",
	"onpurchasegrantguestpasspackage":       "onpurchasegrantguestpasspackage",
	"onpurchasegrantguestpasspackage0":      "onpurchasegrantguestpasspackage0",
	"onpurchasegrantguestpasspackage1":      "onpurchasegrantguestpasspackage1",
	"onpurchasegrantguestpasspackage2":      "onpurchasegrantguestpasspackage2",
	"onpurchasegrantguestpasspackage3":      "onpurchasegrantguestpasspackage3",
	"onpurchasegrantguestpasspackage4":      "onpurchasegrantguestpasspackage4",
	"onpurchasegrantguestpasspackage5":      "onpurchasegrantguestpasspackage5",
	"onpurchasegrantguestpasspackage6":      "onpurchasegrantguestpasspackage6",
	"onpurchasegrantguestpasspackage7":      "onpurchasegrantguestpasspackage7",
	"onpurchasegrantguestpasspackage8":      "onpurchasegrantguestpasspackage8",
	"onpurchasegrantguestpasspackage9":      "onpurchasegrantguestpasspackage9",
	"onpurchasegrantguestpasspackage10":     "onpurchasegrantguestpasspackage10",
	"onpurchasegrantguestpasspackage11":     "onpurchasegrantguestpasspackage11",
	"onpurchasegrantguestpasspackage12":     "onpurchasegrantguestpasspackage12",
	"onpurchasegrantguestpasspackage13":     "onpurchasegrantguestpasspackage13",
	"onpurchasegrantguestpasspackage14":     "onpurchasegrantguestpasspackage14",
	"onpurchasegrantguestpasspackage15":     "onpurchasegrantguestpasspackage15",
	"onpurchasegrantguestpasspackage16":     "onpurchasegrantguestpasspackage16",
	"onpurchasegrantguestpasspackage17":     "onpurchasegrantguestpasspackage17",
	"onpurchasegrantguestpasspackage18":     "onpurchasegrantguestpasspackage18",
	"onpurchasegrantguestpasspackage19":     "onpurchasegrantguestpasspackage19",
	"onpurchasegrantguestpasspackage20":     "onpurchasegrantguestpasspackage20",
	"onpurchasegrantguestpasspackage21":     "onpurchasegrantguestpasspackage21",
	"onpurchasegrantguestpasspackage22":     "onpurchasegrantguestpasspackage22",
	"onquitguestpassmsg":                    "onquitguestpassmsg",
	"overridetaxtype":                       "overridetaxtype",
	"permitrunincountries":                  "permitrunincountries",
	"prohibitrunincountries":                "prohibitrunincountries",
	"purchaserestrictedcountries":           "purchaserestrictedcountries",
	"purchaseretrictedcountries":            "purchaseretrictedcountries",
	"recurringoptions":                      "recurringoptions",
	"recurringpackageoption":                "recurringpackageoption",
	"releaseoverride":                       "releaseoverride",
	"releasestatecountries":                 "releasestatecountries",
	"releasestateoverride":                  "releasestateoverride",
	"releasestateoverridecountries":         "releasestateoverridecountries",
	"relesestateoverride":                   "relesestateoverride",
	"renewalperiod":                         "renewalperiod",
	"renewaltimeunit":                       "renewaltimeunit",
	"requiredps3apploginforpurchase":        "requiredps3apploginforpurchase",
	"requirespreapproval":                   "requirespreapproval",
	"restrictedcountries":                   "restrictedcountries",
	"runrestrictedcountries":                "runrestrictedcountries",
	"shippableitem":                         "shippableitem",
	"skipownsallappsinpackagecheck":         "skipownsallappsinpackagecheck",
	"starttime":                             "starttime",
	"state":                                 "state",
	"test":                                  "test",
	"testchange":                            "testchange",
	"trading_card_drops":                    "trading_card_drops",
	"violencerestrictedcountries":           "violencerestrictedcountries",
	"violencerestrictedterritorycodes":      "violencerestrictedterritorycodes",
	"virtualitemreward":                     "virtualitemreward",
}
