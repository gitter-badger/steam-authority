package steam

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/kr/pretty"
)

/**
https://partner.steamgames.com/doc/webapi/ISteamApps#GetAppList
https://partner.steamgames.com/doc/webapi/ISteamApps#GetCheatingReports
https://partner.steamgames.com/doc/webapi/ISteamApps#GetPlayersBanned

*/

func GetAppList() (apps []stGetAppList, err error) {

	bytes, err := get("ISteamApps/GetAppList/v2/", url.Values{})
	if err != nil {
		return apps, err
	}

	// Unmarshal JSON
	info := stGetAppList{}
	if err := json.Unmarshal(bytes, &info); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return apps, err
	}

	return apps, nil
}

type stGetAppList struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
