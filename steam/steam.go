package steam

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func get(path string, query url.Values, useKey ...bool) (bytes []byte, err error) {

	if path != "" {
		query.Add("format", "json")

		if !(len(useKey) > 0 && !useKey[0]) {
			query.Add("key", os.Getenv("STEAM_API_KEY"))
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
