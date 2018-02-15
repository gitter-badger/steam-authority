package pics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/steam-authority/steam-authority/queue"
	"github.com/Jleagle/go-helpers/logger"
)

const (
	changesLimit = 500
)

// Run triggers the PICS updater to run forever
func Run() {

	for {
		fmt.Println("Checking for changes")

		jsChange, err := GetLatestChanges()
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

		//_, err = saveChangesFromJSON(jsChange)
		//if err != nil {
		//	logger.Error(err)
		//}

		time.Sleep(10 * time.Second)
	}
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
