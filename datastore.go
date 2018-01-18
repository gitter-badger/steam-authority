package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

func createDsAppFromJsApp(js JsApp) *dsApp {

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
	dsApp := dsApp{}
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

func createDsPackageFromJsPackage(js JsPackage) *dsPackage {

	dsPackage := dsPackage{}
	dsPackage.PackageID = js.PackageID
	dsPackage.Apps = js.AppIDs
	dsPackage.BillingType = js.BillingType
	dsPackage.LicenseType = js.LicenseType
	dsPackage.Status = js.Status

	return &dsPackage
}

func savePackage(data dsPackage) {

	packageIDString := strconv.Itoa(data.PackageID)

	key := datastore.NameKey("Package", packageIDString, nil)

	saveKind(key, &data)
}

func saveKind(key *datastore.Key, data interface{}) (newKey *datastore.Key) {

	client, context := getDSClient()
	newKey, err := client.Put(context, key, data)
	if err != nil {
		logger.Error(err)
	}

	return newKey
}

func getDSClient() (*datastore.Client, context.Context) {

	context := context.Background()
	client, err := datastore.NewClient(
		context,
		os.Getenv("STEAM_GOOGLE_PROJECT"),
	)
	if err != nil {
		logger.Error(err)
	}

	return client, context
}

type dsChange struct {
	ChangeID int      `datastore:"change_id"`
	Apps     []string `datastore:"apps"`
	Packages []string `datastore:"packages"`
}

type dsApp struct {
	AppID             int      `datastore:"app_id"`
	Name              string   `datastore:"name"`
	Type              string   `datastore:"type"`
	ReleaseState      string   `datastore:"releasestate"`
	OSList            []string `datastore:"oslist"`
	MetacriticScore   int8     `datastore:"metacritic_score"`
	MetacriticFullURL string   `datastore:"metacritic_fullurl"`
	StoreTags         []int    `datastore:"store_tags"`
	Developer         string   `datastore:"developer"`
	Publisher         string   `datastore:"publisher"`
	Homepage          string   `datastore:"homepage"`
	ChangeNumber      int      `datastore:"change_number"`
}

type dsPackage struct {
	PackageID   int   `datastore:"package_id"`
	BillingType int8  `datastore:"billingtype"`
	LicenseType int8  `datastore:"licensetype"`
	Status      int8  `datastore:"status"`
	Apps        []int `datastore:"apps"`
}
