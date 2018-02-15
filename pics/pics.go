package pics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/queue"
)

const (
	changesLimit = 500
	checkSeconds = 10
)

var latestChangeSaved int

// Run triggers the PICS updater to run forever
func Run() {

	for {
		fmt.Println("Checking for changes")

		jsChange, err := getLatestChanges()
		if err != nil {
			logger.Error(err)
		}

		for k, v := range jsChange.Apps {
			appID, _ := strconv.Atoi(k)
			queue.AppProducer(appID, v)
		}

		for k, v := range jsChange.Packages {
			packageID, _ := strconv.Atoi(k)
			queue.PackageProducer(packageID, v)
		}

		// todo, make unique list of change IDs and call queue.ChangeProducer() on them

		time.Sleep(checkSeconds * time.Second)
	}
}

func getLatestChanges() (jsChange JsChange, err error) {

	// Get the last change
	if latestChangeSaved == 0 {

		changes, err := datastore.GetLatestChanges(1)
		if err != nil {
			logger.Error(err)
		}

		if len(changes) > 0 {
			latestChangeSaved = changes[0].ChangeID
		} else {
			latestChangeSaved = 4059093
		}
	}

	// Grab the JSON from node
	url := "http://localhost:8086/changes/" + strconv.Itoa(latestChangeSaved)
	logger.Info("PICS: " + url)
	response, err := http.Get(url)
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

	latestChangeSaved = jsChange.LatestChangeNumber

	return jsChange, nil
}

type JsChange struct {
	Success            int8           `json:"success"`
	LatestChangeNumber int            `json:"current_changenumber"`
	Apps               map[string]int `json:"apps"`
	Packages           map[string]int `json:"packages"`
}

//func saveChangesFromJSON(jsChange JsChange) (changes []*datastore.Change, err error) {
//
//	// Make a list of changes to add
//	dsChanges := make(map[int]*datastore.Change, 0)
//
//	for k, v := range jsChange.Apps {
//		_, ok := dsChanges[v]
//		if !ok {
//			dsChanges[v] = &datastore.Change{ChangeID: v}
//		}
//
//		intx, _ := strconv.Atoi(k)
//		dsChanges[v].Apps = append(dsChanges[v].Apps, intx)
//	}
//
//	for k, v := range jsChange.Packages {
//		_, ok := dsChanges[v]
//		if !ok {
//			dsChanges[v] = &datastore.Change{ChangeID: v}
//		}
//
//		intx, _ := strconv.Atoi(k)
//		dsChanges[v].Packages = append(dsChanges[v].Packages, intx)
//	}
//
//	// Stop if there are no apps/packages
//	if len(dsChanges) < 1 {
//		return
//	}
//
//	// Make a slice from map
//	var ChangeIDs []int
//	for k := range dsChanges {
//		ChangeIDs = append(ChangeIDs, k)
//	}
//
//	// Datastore can only bulk insert 500, grab the oldest
//	sort.Ints(ChangeIDs)
//	count := int(math.Min(float64(len(ChangeIDs)), changesLimit))
//	ChangeIDs = ChangeIDs[:count]
//
//	dsChangesSlice := make([]*datastore.Change, 0)
//
//	for _, v := range ChangeIDs {
//		dsChangesSlice = append(dsChangesSlice, dsChanges[v])
//	}
//
//	// Bulk add changes
//	err = datastore.BulkAddChanges(dsChangesSlice)
//	if err != nil {
//
//	}
//
//	// Get apps/packages IDs
//	for _, v := range dsChangesSlice {
//
//		info, err := getInfoJSON(v)
//		if err != nil {
//			continue
//		}
//
//		// Build up rows to bulk add apps
//		dsApps := make([]*datastore.DsApp, 0)
//
//		for _, v := range info.Apps {
//			dsApps = append(dsApps, createDsAppFromJsApp(v))
//		}
//
//		websockets.Send(websockets.CHANGES, dsApps)
//
//		err = datastore.BulkAddApps(dsApps)
//		if err != nil {
//			logger.Error(err)
//		}
//
//		// Build up rows to bulk add packages
//		dsPackages := make([]*datastore.DsPackage, 0)
//
//		for _, vv := range info.Packages {
//			dsPackage := createDsPackageFromJsPackage(vv)
//			dsPackage.ChangeID = v.ChangeID
//
//			dsPackages = append(dsPackages, dsPackage)
//		}
//
//		websockets.Send(websockets.CHANGES, dsPackages)
//
//		err = datastore.BulkAddPackages(dsPackages)
//		if err != nil {
//			logger.Error(err)
//		}
//	}
//
//	return nil, err
//}
