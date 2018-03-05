package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/kr/pretty"
)

func GetFriendList(id int) (friends []GetFriendListFriend, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))
	options.Set("relationship", "friend")

	bytes, err := get("ISteamUser/GetFriendList/v1/", options)
	if err != nil {
		return friends, err
	}

	if strings.Contains(string(bytes), "Internal Server Error") {
		return friends, errors.New("no such user")
	}

	// Unmarshal JSON
	var resp *GetFriendListBody
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return friends, err
	}

	return resp.Friendslist.Friends, nil
}

type GetFriendListBody struct {
	Friendslist struct {
		Friends []GetFriendListFriend `json:"friends"`
	} `json:"friendslist"`
}

type GetFriendListFriend struct {
	SteamID      string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendSince  int    `json:"friend_since"`
}

func ResolveVanityURL(id string) (resp ResolveVanityURLBody, err error) {

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
			logger.Error(err)
		}
		return resp, err
	}

	if resp.Response.Success != 1 {
		return resp, errors.New("no user found")
	}

	return resp, nil
}

type ResolveVanityURLBody struct {
	Response ResolveVanityURLResponse
}

type ResolveVanityURLResponse struct {
	SteamID string `json:"steamid"`
	Success int8   `json:"success"`
	Message string `json:"message"`
}

// todo, only return the needed response
func GetPlayerSummaries(ids []int) (resp PlayerSummariesBody, err error) {

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
			logger.Error(err)
		}
		return resp, err
	}

	return resp, nil
}

type PlayerSummariesBody struct {
	Response PlayerSummariesResponse
}

type PlayerSummariesResponse struct {
	Players []PlayerSummariesPlayer
}

type PlayerSummariesPlayer struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	LastLogOff               int64  `json:"lastlogoff"`
	CommentPermission        int    `json:"commentpermission"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            string `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LOCCountryCode           string `json:"loccountrycode"`
	LOCStateCode             string `json:"locstatecode"`
}

func GetPlayerBans(id int) (bans GetPlayerBanResponse, err error) {

	options := url.Values{}
	options.Set("steamids", strconv.Itoa(id))

	bytes, err := get("ISteamUser/GetPlayerBans/v1", options)
	if err != nil {
		return bans, err
	}

	// Unmarshal JSON
	var resp GetPlayerBansResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return bans, err
	}

	if len(resp.Players) == 0 {
		return bans, nil
	}

	return resp.Players[0], nil
}

type GetPlayerBansResponse struct {
	Players []GetPlayerBanResponse `json:"players"`
}

type GetPlayerBanResponse struct {
	SteamID          string `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

func GetUserGroupList(id int) (groups []int, err error) {

	options := url.Values{}
	options.Set("steamid", strconv.Itoa(id))

	bytes, err := get("ISteamUser/GetUserGroupList/v1", options)
	if err != nil {
		return bans, err
	}

	// Unmarshal JSON
	var resp GetPlayerBansResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return bans, err
	}

	return resp.Response, nil
}

type AutoGenerated struct {
	Response struct {
		Success bool `json:"success"`
		Groups  []struct {
			GID string `json:"gid"`
		} `json:"groups"`
	} `json:"response"`
}
