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
		apps = append(apps, strconv.Itoa(vv))
	}
	for _, vv := range change.Packages {
		packages = append(packages, strconv.Itoa(vv))
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

		int, _ := strconv.Atoi(k)
		dsChanges[v].Apps = append(dsChanges[v].Apps, int)
	}

	for k, v := range jsChange.Packages {
		_, ok := dsChanges[v]
		if !ok {
			dsChanges[v] = &datastore.DsChange{ChangeID: v}
		}

		int, _ := strconv.Atoi(k)
		dsChanges[v].Packages = append(dsChanges[v].Packages, int)
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
		websockets.Send(websockets.CHANGES, dsChanges[v])
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
