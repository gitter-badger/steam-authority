package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/kr/pretty"
)

/**
https://partner.steamgames.com/doc/webapi/ISteamApps#GetCheatingReports
https://partner.steamgames.com/doc/webapi/ISteamApps#GetPlayersBanned
*/

// todo, list of apps that wont unmarshal
// 2130, 2720, 4720, 4580, 4580
func GetAppDetails(id string) (app AppDetailsBody, err error) {

	options := url.Values{}
	options.Set("appids", id)

	bytes, err := get("", options)
	if err != nil {
		return app, err
	}

	// Check for no app
	if string(bytes) == "null" {
		return app, errors.New("invalid app id")
	}

	// Fix values that can change type, causing unmarshal errors
	b := strings.Replace(string(bytes), "\"pc_requirements\":[],", "\"pc_requirements\":null,", 1)
	b = strings.Replace(b, "\"mac_requirements\":[],", "\"mac_requirements\":null,", 1)
	b = strings.Replace(b, "\"linux_requirements\":[],", "\"linux_requirements\":null,", 1)
	bytes = []byte(b)

	// Unmarshal JSON
	resp := make(map[string]AppDetailsBody)
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return app, err
	}

	if resp[id].Success == false {
		return app, errors.New("no app with id")
	}

	return resp[id], nil
}

type AppDetailsBody struct {
	Success bool `json:"success"`
	Data struct {
		Type                string `json:"type"`
		Name                string `json:"name"`
		SteamAppID          int    `json:"steam_appid"`
		RequiredAge         int    `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		DLC                 []int  `json:"dlc"`
		ControllerSupport   string `json:"controller_support"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		Fullgame struct {
			AppID string `json:"appid"`
			Name  string `json:"name"`
		} `json:"fullgame"`
		SupportedLanguages string `json:"supported_languages"`
		Reviews            string `json:"reviews"`
		HeaderImage        string `json:"header_image"`
		Website            string `json:"website"`
		PcRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"`
		MacRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"mac_requirements"`
		LinuxRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"linux_requirements"`
		LegalNotice string   `json:"legal_notice"`
		Developers  []string `json:"developers"`
		Publishers  []string `json:"publishers"`
		Demos []struct {
			Appid       int    `json:"appid"`
			Description string `json:"description"`
		} `json:"demos"`
		PriceOverview struct {
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
		} `json:"price_overview"`
		Packages []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             int    `json:"display_type"`
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs []struct {
				Packageid                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        string `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int8   `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories  []AppDetailsCategory   `json:"categories"`
		Genres      []AppDetailsGenre      `json:"genres"`
		Screenshots []AppDetailsScreenshot `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements AppDetailsAchievements `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background string `json:"background"`
	} `json:"data"`
}

type AppDetailsScreenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type AppDetailsAchievements struct {
	Total int `json:"total"`
	Highlighted []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"highlighted"`
}

type AppDetailsGenre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type AppDetailsCategory struct {
	ID          int8   `json:"id"`
	Description string `json:"description"`
}

func GetAppList() (apps []GetAppListApp, err error) {

	bytes, err := get("ISteamApps/GetAppList/v2/", url.Values{})
	if err != nil {
		return apps, err
	}

	// Unmarshal JSON
	resp := GetAppListBody{}
	if err := json.Unmarshal(bytes, &resp); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			pretty.Print(string(bytes))
		}
		return apps, err
	}

	return resp.AppList.Apps, nil
}

type GetAppListBody struct {
	AppList GetAppListAppList `json:"applist"`
}

type GetAppListAppList struct {
	Apps []GetAppListApp `json:"apps"`
}

type GetAppListApp struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
