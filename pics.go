package main

// JSON ...
type JSON struct {
	success         int8
	apps            map[int]JsApp
	packages        map[int]JsPackage
	UnknownApps     []int `json:"unknown_apps"`
	UnknownPackages []int `json:"unknown_packages"`
}

// JsApp ...
type JsApp struct {
	appID        int
	changeNumber int
}

// JsAppCommon ...
type JsAppCommon struct {
	icon          string
	logo          string
	LogoSmall     string `json:"logo_small"`
	metacriticURL string
	name          string
	ClientIcon    string `json:"clienticon"`
	ClientTga     string `json:"clienttga"`
	languages     map[string]int8
	// todo
}

// JsAppExtended ...
type JsAppExtended struct {
	developer           string
	DeveloperURL        string `json:"developer_url"`
	GameDir             string `json:"gamedir"`
	GameManualURL       string `json:"gamemanualurl"`
	homepage            string
	icon                string
	icon2               string
	IsFreeApp           int8 `json:"isfreeapp"`
	languages           string
	LoadAllBeforeLaunch int8   `json:"loadallbeforelaunch"`
	MinClientVersion    string `json:"minclientversion"`
	NoServers           int8   `json:"noservers"`
	PrimaryCache        int    `json:"primarycache"`
	PrimaryCacheLinux   int    `json:"primarycache_linux"`
	RequireSSSE         int8   `json:"requiressse"`
	ServerBrowserName   string `json:"serverbrowsername"`
	sourceGame          int8
	state               string
	VacMacModuleCache   int    `json:"vacmacmodulecache"`
	VacModuleCache      int    `json:"vacmodulecache"`
	VacModuleFilename   string `json:"vacmodulefilename"`
	ValidosList         string `json:"validoslist"`
	publisher           string
	ListofDLC           string `json:"listofdlc"`
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
	//AppItems    []int             `json:"appitems"` // todo
}

// JsPackageExtended ...
type JsPackageExtended struct {
	Alwayscountasowned   int8   `json:"alwayscountasowned"`
	Devcomp              int8   `json:"devcomp"`
	Releasestateoverride string `json:"releasestateoverride"`
}
