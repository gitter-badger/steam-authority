package pics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/websockets"
)

const (
	changesLimit = 500
)

func Run() {

	for {
		fmt.Println("Checking for changes")

		changes, err := datastore.GetLatestChanges(1)
		if err != nil {
			logger.Error(err)
		}

		jsChange, err := getChangesJson(changes[0])
		if err != nil {
			logger.Error(err)
		}

		_, err = saveChangesFromJson(jsChange)
		if err != nil {
			logger.Error(err)
		}

		time.Sleep(10 * time.Second)
	}
}

func getChangesJson(latestChange datastore.DsChange) (jsChange JsChange, err error) {

	// Grab the JSON from node
	response, err := http.Get("http://localhost:8086/changes/" + strconv.Itoa(latestChange.ChangeID))
	if err != nil {
		return jsChange, err
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsChange, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(contents, &jsChange); err != nil {
		return jsChange, err
	}

	return jsChange, err
}

func getInfoJson(change *datastore.DsChange) (jsInfo JsInfo, err error) {

	apps := []string{}
	packages := []string{}

	for _, vv := range change.Apps {
		apps = append(apps, vv)
	}
	for _, vv := range change.Packages {
		packages = append(packages, vv)
	}

	// Grab the JSON from node
	response, err := http.Get("http://localhost:8086/info?apps=" + strings.Join(apps, ",") + "&packages=" + strings.Join(packages, ",") + "&prettyprint=0")
	if err != nil {
		logger.Error(err)
		return jsInfo, err
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
	}

	// Unmarshal JSON
	info := JsInfo{}
	if err := json.Unmarshal(contents, &info); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(contents))
		}
		logger.Error(err)
	}

	return info, nil
}

func saveChangesFromJson(jsChange JsChange) (changes []*datastore.DsChange, err error) {

	// Make a list of changes to add
	dsChanges := make(map[int]*datastore.DsChange, 0)

	for k, v := range jsChange.Apps {
		_, ok := dsChanges[v]
		if !ok {
			dsChanges[v] = &datastore.DsChange{ChangeID: v}
		}
		dsChanges[v].Apps = append(dsChanges[v].Apps, k)
	}

	for k, v := range jsChange.Packages {
		_, ok := dsChanges[v]
		if !ok {
			dsChanges[v] = &datastore.DsChange{ChangeID: v}
		}
		dsChanges[v].Packages = append(dsChanges[v].Packages, k)
	}

	// Stop if there are no apps/packages
	if len(dsChanges) < 1 {
		return
	}

	// Make a slice from map
	var ChangeIDs []int
	for k := range dsChanges {
		ChangeIDs = append(ChangeIDs, k)
	}

	// Datastore can only bulk insert 500, grab the oldest
	sort.Ints(ChangeIDs)
	count := int(math.Min(float64(len(ChangeIDs)), changesLimit))
	ChangeIDs = ChangeIDs[:count]

	dsChangesSlice := make([]*datastore.DsChange, 0)

	for _, v := range ChangeIDs {
		dsChangesSlice = append(dsChangesSlice, dsChanges[v])
		websockets.Send(dsChanges[v])
	}

	// Bulk add changes
	err = datastore.BulkAddChanges(dsChangesSlice)
	if err != nil {

	}

	// Get apps/packages IDs
	for _, v := range dsChangesSlice {

		info, err := getInfoJson(v)
		if err != nil {
			continue
		}

		// Build up rows to bulk add apps
		dsApps := make([]*datastore.DsApp, 0)

		for _, v := range info.Apps {
			dsApps = append(dsApps, createDsAppFromJsApp(v))
		}

		err = datastore.BulkAddApps(dsApps)
		if err != nil {
			logger.Error(err)
		}

		// Build up rows to bulk add packages
		dsPackages := make([]*datastore.DsPackage, 0)

		for _, v := range info.Packages {
			dsPackages = append(dsPackages, createDsPackageFromJsPackage(v))
		}

		err = datastore.BulkAddPackages(dsPackages)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil, err
}

func createDsAppFromJsApp(js JsApp) *datastore.DsApp {

	// Convert map of tags to slice
	jsTags := js.Common.StoreTags
	tags := make([]int, 0, len(jsTags))
	for _, value := range jsTags {
		valueInt, _ := strconv.Atoi(value)
		tags = append(tags, valueInt)
	}

	// String to int
	appIDInt, _ := strconv.Atoi(js.AppID)
	metacriticScoreInt, _ := strconv.Atoi(js.Common.MetacriticScore)

	//
	dsApp := datastore.DsApp{}
	dsApp.AppID = appIDInt
	dsApp.Name = js.Common.Name
	dsApp.Type = js.Common.Type
	dsApp.ReleaseState = js.Common.ReleaseState
	dsApp.OSList = strings.Split(js.Common.OSList, ",")
	dsApp.MetacriticScore = int8(metacriticScoreInt)
	dsApp.MetacriticFullURL = js.Common.MetacriticURL
	dsApp.StoreTags = tags
	dsApp.Developer = js.Extended.Developer
	dsApp.Publisher = js.Extended.Publisher
	dsApp.Homepage = js.Extended.Homepage
	dsApp.ChangeNumber = js.ChangeNumber

	return &dsApp
}

func createDsPackageFromJsPackage(js JsPackage) *datastore.DsPackage {

	dsPackage := datastore.DsPackage{}
	dsPackage.PackageID = js.PackageID
	dsPackage.Apps = js.AppIDs
	dsPackage.BillingType = js.BillingType
	dsPackage.LicenseType = js.LicenseType
	dsPackage.Status = js.Status

	return &dsPackage
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
	PackageID   int  `json:"packageid"`
	BillingType int8 `json:"billingtype"`
	LicenseType int8 `json:"licensetype"`
	Status      int8 `json:"status"`
	// Extended    JsPackageExtended `json:"extended"` // Sometimes shows as empty array, breaking unmarshal
	AppIDs   []int `json:"appids"`
	DepotIDs []int `json:"depotids"`
	AppItems []int `json:"appitems"` // todo, no data to test with
}

// JsPackageExtended ...
type JsPackageExtended struct {
	AlwaysCountsAsOwned               int8   `json:"alwayscountsasowned"`
	DevComp                           int8   `json:"devcomp"`
	ReleaseStateOverride              string `json:"releasestateoverride"`
	AllowCrossRegionTradingAndGifting string `json:"allowcrossregiontradingandgifting"`
}
