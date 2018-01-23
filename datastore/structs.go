package datastore

// DsPlayer kind
type DsPlayer struct {
	ID64            int
	ValintyURL      string `datastore:"vality_url"`
	Avatar          string `datastore:"avatar"`
	RealName        string `datastore:"real_name"`
	CountryCode     string `datastore:"country_code"`
	StateCode       string `datastore:"status_code"`
	LastUpdated     int64  `datastore:"last_updated"`
	Level           int    `datastore:"level"`
	LevelRank       int    `datastore:"level_rank"`
	Games           int    `datastore:"games"`
	GamesRank       int    `datastore:"games_rank"`
	Badges          int    `datastore:"badges"`
	BadgesRank      int    `datastore:"badges_rank"`
	PlayTime        int    `datastore:"play_time"`
	PlayTimeRank    int    `datastore:"play_time_rank"`
	TimeCreated     int    `datastore:"time_created"`
	TimeCreatedRank int    `datastore:"time_created_rank"`
}

// DsChange kind
type DsChange struct {
	ChangeID int   `datastore:"change_id"`
	Apps     []int `datastore:"apps"`
	Packages []int `datastore:"packages"`
}

// DsApp kind
type DsApp struct {
	AppID             int      `datastore:"app_id"`
	Name              string   `datastore:"name"`
	Type              string   `datastore:"type"`
	ReleaseState      string   `datastore:"releasestate"`
	OSList            []string `datastore:"oslist"`
	MetacriticScore   int8     `datastore:"metacritic_score"`
	MetacriticFullURL string   `datastore:"metacritic_fullurl"`
	StoreTags         []int    `datastore:"store_tags"`
	Developer         string   `datastore:"developer"`
	Publisher         string   `datastore:"publisher"`
	Homepage          string   `datastore:"homepage"`
	ChangeNumber      int      `datastore:"change_number"`
	Logo              string   `datastore:"logo"`
	Icon              string   `datastore:"icon"`
}

// DsPackage kind
type DsPackage struct {
	PackageID   int   `datastore:"package_id"`
	BillingType int8  `datastore:"billingtype"`
	LicenseType int8  `datastore:"licensetype"`
	Status      int8  `datastore:"status"`
	Apps        []int `datastore:"apps"`
	ChangeID    int   `datastore:"change_id"`
}
