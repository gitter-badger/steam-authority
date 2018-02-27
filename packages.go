package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/mysql"
)

func packagesHandler(w http.ResponseWriter, r *http.Request) {

	packages, err := mysql.GetLatestPackages(100)
	if err != nil {
		logger.Error(err)
	}

	template := packagesTemplate{}
	template.Fill(r)
	template.Packages = packages

	returnTemplate(w, r, "packages", template)
}

type packagesTemplate struct {
	GlobalTemplate
	Packages []mysql.Package
}

func packageHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		returnErrorTemplate(w, r, 404, err.Error())
		return
	}

	pack, err := mysql.GetPackage(id)
	if err != nil {

		if err.Error() == "no id" {
			returnErrorTemplate(w, r, 404, "We can't find this package in our database, there may not be one with this ID.")
			return
		}

		logger.Error(err)
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	appIDs, err := pack.GetApps()
	if err != nil {
		logger.Error(err)
	}

	apps, err := mysql.GetApps(appIDs, []string{"id", "icon", "type", "platforms", "dlc"})
	if err != nil {
		logger.Error(err)
	}
	// Make banners
	banners := make(map[string][]string)
	var primary []string

	//if pack.GetExtended() == "prerelease" {
	//	primary = append(primary, "This package is intended for developers and publishers only.")
	//}

	if len(primary) > 0 {
		banners["primary"] = primary
	}

	// Template
	template := packageTemplate{}
	template.Fill(r)
	template.Package = pack
	template.Apps = apps
	template.Keys = packageKeys

	returnTemplate(w, r, "package", template)
}

type packageTemplate struct {
	GlobalTemplate
	Package mysql.Package
	Apps    []mysql.App
	Keys    map[string]string
	Banners map[string][]string
}

// todo, make these nice, put into the GetExtended func?
var packageKeys = map[string]string{
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
