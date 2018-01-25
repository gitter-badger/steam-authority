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
https://partner.steamgames.com/doc/webapi/ISteamUser#ResolveVanityURL
https://partner.steamgames.com/doc/webapi/ISteamUser#GetUserGroupList
https://partner.steamgames.com/doc/webapi/ISteamUser#GetPlayerBans
https://partner.steamgames.com/doc/webapi/ISteamUser#GetFriendList
https://partner.steamgames.com/doc/webapi/ISteamUser#GetAppPriceInfo
*/

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
	resp = getPlayerSummariesBody{}
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return resp, err
	}

	return resp, nil
}

type getPlayerSummariesBody struct {
	Response getPlayerSummariesResponse
}

type getPlayerSummariesResponse struct {
	Players []getPlayerSummaries
}

type getPlayerSummaries struct {
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
