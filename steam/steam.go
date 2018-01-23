package steam

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
)

var apiKey = os.Getenv("STEAM_API_KEY")

func GetAppList() (apps []StGetAppList) {

	bytes, err := get("ISteamApps/GetAppList/v2", url.Values{})
	if err != nil {
		logger.Error(err)
	}

	// Unmarshal JSON
	info := StGetAppList{}
	if err := json.Unmarshal(bytes, &info); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		logger.Error(err)
	}

	return apps
}

func get(path string, query url.Values, useKey ...bool) (bytes []byte, err error) {

	if path != "" {
		query.Add("format", "json")

		if !(len(useKey) > 0 && !useKey[0]) {
			query.Add("key", apiKey)
		}
		path = "http://api.steampowered.com/" + path
	} else {
		path = "http://store.steampowered.com/api/appdetails"
	}

	// Grab the JSON from node
	response, err := http.Get(path + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, err
}
