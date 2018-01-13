package main

// JSON ...
type JSON struct {
	Success         int8              `json:"success"`
	Apps            map[int]JsApp     `json:"apps"`
	Packages        map[int]JsPackage `json:"packages"`
	UnknownApps     []int             `json:"unknown_apps"`
	UnknownPackages []int             `json:"unknown_packages"`
}

// JsApp ...
type JsApp struct {
	AppID              string                     `json:"appid"`
	Common             JsAppCommon                `json:"common"`
	Extended           JsAppExtended              `json:"extended"`
	Config             JsAppConfig                `json:"config"`
	Depots             JsAppDepots                `json:"depots"`
	UFS                JsAppUFS                   `json:"ufs"`
	SystemRequirements JsAppUFSSystemRequirements `json:"sysreqs"`
	ChangeNumber       int                        `json:"change_number"`
}

// JsAppCommon ...
type JsAppCommon struct {
	Icon                  string            `json:"icon"`
	Logo                  string            `json:"logo"`
	LogoSmall             string            `json:"logo_small"`
	MetacriticURL         string            `json:"metacritic_url"`
	Name                  string            `json:"name"`
	ClientIcon            string            `json:"clienticon"`
	ClientTga             string            `json:"clienttga"`
	Languages             map[string]string `json:"languages"`
	ClientICNS            string            `json:"clienticns"`
	LinuxClientIcon       string            `json:"linuxclienticon"`
	OSList                string            `json:"oslist"`
	Type                  string            `json:"type"`
	MetacriticName        string            `json:"metacritic_name"`
	ControllerSupport     string            `json:"controller_support"`
	SmallCapsule          map[string]string `json:"small_capsule"`
	HeaderImage           map[string]string `json:"header_image"`
	MetacriticScore       string            `json:"metacritic_score"`
	MetacriticFullurl     string            `json:"metacritic_fullurl"`
	CommunityVisibleStats string            `json:"community_visible_stats"`
	WorkshopVisible       string            `json:"workshop_visible"`
	CommunityHubVisible   string            `json:"community_hub_visible"`
	GameID                string            `json:"gameid"`
	Exfgls                string            `json:"exfgls"`
	StoreTags             map[string]string `json:"store_tags"`
}

// JsAppExtended ...
type JsAppExtended struct {
	Developer           string `json:"developer"`
	DeveloperURL        string `json:"developer_url"`
	GameDir             string `json:"gamedir"`
	GameManualURL       string `json:"gamemanualurl"`
	Homepage            string `json:"homepage"`
	Icon                string `json:"icon"`
	Icon2               string `json:"icon2"`
	IsFreeApp           string `json:"isfreeapp"`
	Languages           string `json:"languages"`
	LoadAllBeforeLaunch string `json:"loadallbeforelaunch"`
	MinClientVersion    string `json:"minclientversion"`
	NoServers           string `json:"noservers"`
	PrimaryCache        string `json:"primarycache"`
	PrimaryCacheLinux   string `json:"primarycache_linux"`
	RequireSSSE         string `json:"requiressse"`
	ServerBrowserName   string `json:"serverbrowsername"`
	SourceGame          string `json:"sourcegame"`
	State               string `json:"state"`
	VacMacModuleCache   string `json:"vacmacmodulecache"`
	VacModuleCache      string `json:"vacmodulecache"`
	VacModuleFilename   string `json:"vacmodulefilename"`
	ValidosList         string `json:"validoslist"`
	Publisher           string `json:"publisher"`
	ListofDLC           string `json:"listofdlc"`
}

// JsAppConfig ...
type JsAppConfig struct {
}

// JsAppDepots ...
type JsAppDepots struct {
}

// JsAppUFS ...
type JsAppUFS struct {
	Quota       string `json:"quota"`
	MaxNumFiles string `json:"maxnumfiles"`
	HideCloudUI string `json:"hidecloudui"`
}

// JsAppUFSSystemRequirements ...
type JsAppUFSSystemRequirements struct {
}

// JsPackage ...
type JsPackage struct {
	PackageID   int               `json:"packageid"`
	BillingType int8              `json:"billingtype"`
	LicenseType int8              `json:"licensetype"`
	Status      int8              `json:"status"`
	Extended    JsPackageExtended `json:"extended"`
	AppIDs      []int             `json:"appids"`
	DepotIDs    []int             `json:"depotids"`
	//AppItems    []int             `json:"appitems"` // todo, no data to test with
}

// JsPackageExtended ...
type JsPackageExtended struct {
	Alwayscountasowned   int8   `json:"alwayscountasowned"`
	Devcomp              int8   `json:"devcomp"`
	Releasestateoverride string `json:"releasestateoverride"`
}
