package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/kr/pretty"
)

/**
https://partner.steamgames.com/doc/webapi/ISteamUser#GetUserGroupList
https://partner.steamgames.com/doc/webapi/ISteamUser#GetPlayerBans
https://partner.steamgames.com/doc/webapi/ISteamUser#GetFriendList
https://partner.steamgames.com/doc/webapi/ISteamUser#GetAppPriceInfo
*/

func ResolveVanityURL(id string) (resp resolveVanityURLBody, err error) {

	options := url.Values{}
	options.Set("vanityurl", id)
	options.Set("url_type", "1")

	bytes, err := get("ISteamUser/ResolveVanityURL/v1/", options)
	if err != nil {
		return resp, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return resp, err
	}

	if resp.Response.Success != 1 {
		return resp, errors.New("No user found")
	}

	return resp, nil
}

type resolveVanityURLBody struct {
	Response struct {
		SteamID string `json:"steamid"`
		Success int8   `json:"success"`
		Message string `json:"message"`
	}
}

func GetPlayerSummaries(ids []int) (resp getPlayerSummariesBody, err error) {

	if len(ids) > 100 {
		return resp, errors.New("100 ids max")
	}

	var idsString []string
	for _, v := range ids {
		idsString = append(idsString, strconv.Itoa(v))
	}

	options := url.Values{}
	options.Set("steamids", strings.Join(idsString, ","))

	bytes, err := get("ISteamUser/GetPlayerSummaries/v2/", options)
	if err != nil {
		return resp, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return resp, err
	}

	return resp, nil
}

type getPlayerSummariesBody struct {
	Response struct {
		Players []struct {
			SteamID                  string `json:"steamid"`
			CommunityVisibilityState int8   `json:"communityvisibilitystate"`
			ProfileState             int8   `json:"profilestate"`
			PersonaName              string `json:"personaname"`
			LastLogOff               int64  `json:"lastlogoff"`
			CommentPermission        int8   `json:"commentpermission"`
			ProfileURL               string `json:"profileurl"`
			Avatar                   string `json:"avatar"`
			AvatarMedium             string `json:"avatarmedium"`
			AvatarFull               string `json:"avatarfull"`
			PersonaState             int8   `json:"personastate"`
			RealName                 string `json:"realname"`
			PrimaryClanID            string `json:"primaryclanid"`
			TimeCreated              int64  `json:"timecreated"`
			PersonaStateFlags        int8   `json:"personastateflags"`
			LOCCountryCode           string `json:"loccountrycode"`
			LOCStateCode             string `json:"locstatecode"`
		}
	}
}
