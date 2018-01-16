package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

func checkForChanges() {

	// Get the latest change to start from
	client, context := getDSClient()
	q := datastore.NewQuery("Change").Order("-change_id").Limit(1)
	it := client.Run(context, q)

	var latestChange dsChange
	_, err := it.Next(&latestChange)
	if err != nil {
		logger.Error(err)
	}

	// Grab the JSON from node
	response, err := http.Get("http://localhost:8086/changes/" + strconv.Itoa(latestChange.ChangeID))
	if err != nil {
		logger.Error(err)
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
	}

	// Unmarshal JSON
	jsChange := JsChange{}
	if err := json.Unmarshal(contents, &jsChange); err != nil {
		logger.Error(err)
	}

	// Make a list of changes to add
	dsChanges := make(map[int]*dsChange, 0)

	for k, v := range jsChange.Apps {

		_, ok := dsChanges[v]
		if !ok {
			dsChanges[v] = &dsChange{ChangeID: v}
		}

		dsChanges[v].Apps = append(dsChanges[v].Apps, k)
	}

	for k, v := range jsChange.Packages {

		_, ok := dsChanges[v]
		if !ok {
			dsChanges[v] = &dsChange{ChangeID: v}
		}

		dsChanges[v].Packages = append(dsChanges[v].Packages, k)
	}

	// Stop if there are no apps/packages
	if len(dsChanges) < 1 {
		return
	}

	var ChangeIDs []int
	for k := range dsChanges {
		ChangeIDs = append(ChangeIDs, k)
	}

	// Datastore can only bulk insert 500, so grab the oldest 500
	sort.Ints(ChangeIDs)
	count := int(math.Min(float64(len(ChangeIDs)), 500))
	ChangeIDs = ChangeIDs[:count]

	dsKeys := make([]*datastore.Key, 0)
	dsChangesSlice := make([]*dsChange, 0)

	for _, v := range ChangeIDs {
		dsKeys = append(dsKeys, datastore.NameKey("Change", strconv.Itoa(v), nil))
		dsChangesSlice = append(dsChangesSlice, dsChanges[v])
	}

	// Bulk add changes
	fmt.Println("Saving " + strconv.Itoa(count) + " change(s)")

	if _, err := client.PutMulti(context, dsKeys, dsChangesSlice); err != nil {
		logger.Error(err)
	}

	// Get apps/packages IDs
	apps := []string{}
	packages := []string{}

	for _, v := range dsChanges {
		for _, vv := range v.Apps {
			apps = append(apps, vv)
		}
		for _, vv := range v.Packages {
			packages = append(packages, vv)
		}
	}

	// Grab the JSON from node
	response, err = http.Get("http://localhost:8086/info?apps=" + strings.Join(apps, ",") + "&packages=" + strings.Join(packages, ",") + "&prettyprint=0")
	if err != nil {
		logger.Error(err)
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
	}

	// Unmarshal JSON
	info := JsInfo{}
	if err := json.Unmarshal(contents, &info); err != nil {
		logger.Error(err)
	}

	// Build up rows to bulk add apps
	dsApps := make([]*dsApp, 0)
	dsAppKeys := make([]*datastore.Key, 0)

	for _, v := range info.Apps {
		dsApps = append(dsApps, createDsAppFromJsApp(v))
		dsAppKeys = append(dsAppKeys, datastore.NameKey("App", v.AppID, nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len(dsApps)) + " app(s)")
	if _, err := client.PutMulti(context, dsAppKeys, dsApps); err != nil {
		logger.Error(err)
	}

	// Build up rows to bulk add packages
	dsPackages := make([]*dsPackage, 0)
	dsPackageKeys := make([]*datastore.Key, 0)

	for _, v := range info.Packages {
		dsPackages = append(dsPackages, createDsPackageFromJsPackage(v))
		dsPackageKeys = append(dsPackageKeys, datastore.NameKey("Package", strconv.Itoa(v.PackageID), nil))
	}

	fmt.Println("Saving " + strconv.Itoa(len(dsPackages)) + " package(s)")
	if _, err := client.PutMulti(context, dsPackageKeys, dsPackages); err != nil {
		logger.Error(err)
	}
}

// JsChange ...
type JsChange struct {
	Success            int8           `json:"success"`
	LatestChangeNumber int            `json:"current_changenumber"`
	Apps               map[string]int `json:"apps"`
	Packages           map[string]int `json:"packages"`
}

// JsInfo ...
type JsInfo struct {
	Success         int8              `json:"success"`
	Apps            map[int]JsApp     `json:"apps"`
	Packages        map[int]JsPackage `json:"packages"`
	UnknownApps     []int             `json:"unknown_apps"`
	UnknownPackages []int             `json:"unknown_packages"`
}

// JsApp ...
type JsApp struct {
	AppID              string                  `json:"appid"`
	PublicOnly         string                  `json:"public_only"`
	Common             JsAppCommon             `json:"common"`
	Extended           JsAppExtended           `json:"extended"`
	Config             JsAppConfig             `json:"config"`
	Depots             JsAppDepots             `json:"depots"`
	UFS                JsAppUFS                `json:"ufs"`
	SystemRequirements JsAppSystemRequirements `json:"sysreqs"`
	ChangeNumber       int                     `json:"change_number"`
}

// JsAppCommon ...
type JsAppCommon struct {
	ClientICNS            string                     `json:"clienticns"`
	ClientIcon            string                     `json:"clienticon"`
	ClientTGA             string                     `json:"clienttga"`
	CommunityHubVisible   string                     `json:"community_hub_visible"`
	CommunityVisibleStats string                     `json:"community_visible_stats"`
	ControllerSupport     string                     `json:"controller_support"`
	EULAs                 map[string]JsAppCommonEULA `json:"eulas"`
	Exfgls                string                     `json:"exfgls"`
	GameID                string                     `json:"gameid"`
	HeaderImage           map[string]string          `json:"header_image"`
	Icon                  string                     `json:"icon"`
	Languages             map[string]string          `json:"languages"`
	LinuxClientIcon       string                     `json:"linuxclienticon"`
	Logo                  string                     `json:"logo"`
	LogoSmall             string                     `json:"logo_small"`
	MetacriticFullurl     string                     `json:"metacritic_fullurl"`
	MetacriticName        string                     `json:"metacritic_name"`
	MetacriticScore       string                     `json:"metacritic_score"`
	MetacriticURL         string                     `json:"metacritic_url"`
	Name                  string                     `json:"name"`
	OSList                string                     `json:"oslist"`
	OSArch                string                     `json:"osarch"`
	ReleaseState          string                     `json:"releasestate"`
	SmallCapsule          map[string]string          `json:"small_capsule"`
	StoreTags             map[string]string          `json:"store_tags"`
	Type                  string                     `json:"type"`
	WorkshopVisible       string                     `json:"workshop_visible"`
}

// JsAppCommonEULA ...
type JsAppCommonEULA struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// JsAppExtended ...
type JsAppExtended struct {
	Developer           string `json:"developer"`
	DeveloperURL        string `json:"developer_url"`
	GameDir             string `json:"gamedir"`
	GameManualURL       string `json:"gamemanualurl"`
	Homepage            string `json:"homepage"`
	Icon                string `json:"icon"`
	Icon2               string `json:"icon2"`
	IsFreeApp           string `json:"isfreeapp"`
	Languages           string `json:"languages"`
	LoadAllBeforeLaunch string `json:"loadallbeforelaunch"`
	MinClientVersion    string `json:"minclientversion"`
	NoServers           string `json:"noservers"`
	PrimaryCache        string `json:"primarycache"`
	PrimaryCacheLinux   string `json:"primarycache_linux"`
	RequireSSSE         string `json:"requiressse"`
	ServerBrowserName   string `json:"serverbrowsername"`
	SourceGame          string `json:"sourcegame"`
	State               string `json:"state"`
	VacMacModuleCache   string `json:"vacmacmodulecache"`
	VacModuleCache      string `json:"vacmodulecache"`
	VacModuleFilename   string `json:"vacmodulefilename"`
	ValidosList         string `json:"validoslist"`
	Publisher           string `json:"publisher"`
	ListofDLC           string `json:"listofdlc"`
}

// JsAppConfig ...
type JsAppConfig struct {
}

// JsAppDepots ...
type JsAppDepots struct {
}

// JsAppUFS ...
type JsAppUFS struct {
	Quota       string `json:"quota"`
	MaxNumFiles string `json:"maxnumfiles"`
	HideCloudUI string `json:"hidecloudui"`
}

// JsAppSystemRequirements ...
type JsAppSystemRequirements struct {
}

// JsPackage ...
type JsPackage struct {
	PackageID   int               `json:"packageid"`
	BillingType int8              `json:"billingtype"`
	LicenseType int8              `json:"licensetype"`
	Status      int8              `json:"status"`
	Extended    JsPackageExtended `json:"extended"`
	AppIDs      []int             `json:"appids"`
	DepotIDs    []int             `json:"depotids"`
	AppItems    []int             `json:"appitems"` // todo, no data to test with
}

// JsPackageExtended ...
type JsPackageExtended struct {
	AlwaysCountAsOwned                int8   `json:"alwayscountasowned"`
	DevComp                           int8   `json:"devcomp"`
	ReleaseStateOverride              string `json:"releasestateoverride"`
	AllowCrossRegionTradingAndGifting string `json:"allowcrossregiontradingandgifting"`
}
