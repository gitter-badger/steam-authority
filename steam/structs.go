package steam

type StGetAppList struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
