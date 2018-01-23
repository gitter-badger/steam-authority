package steam

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
)

func GetAppList() (apps []StGetAppList) {

	bytes, err := get("ISteamApps/GetAppList/v2/", url.Values{})
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

type StGetAppList struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
