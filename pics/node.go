package pics

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
	"github.com/steam-authority/steam-authority/datastore"
)

var latestChangeSaved int

func GetLatestChanges() (jsChange JsChange, err error) {

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

func GetInfo(apps []int, packages []int) (jsInfo JsInfo, err error) {

	var stringApps []string
	var stringPackages []string

	for _, vv := range apps {
		stringApps = append(stringApps, strconv.Itoa(vv))
	}
	for _, vv := range packages {
		stringPackages = append(stringPackages, strconv.Itoa(vv))
	}

	// Grab the JSON from node
	url := "http://localhost:8086/info?apps=" + strings.Join(stringApps, ",") + "&packages=" + strings.Join(stringPackages, ",") + "&prettyprint=0"
	logger.Info("PICS: " + url)
	response, err := http.Get(url)
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
