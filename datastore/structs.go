package datastore

type DsChange struct {
	ChangeID int      `datastore:"change_id"`
	Apps     []string `datastore:"apps"`
	Packages []string `datastore:"packages"`
}

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
}

type DsPackage struct {
	PackageID   int   `datastore:"package_id"`
	BillingType int8  `datastore:"billingtype"`
	LicenseType int8  `datastore:"licensetype"`
	Status      int8  `datastore:"status"`
	Apps        []int `datastore:"apps"`
}
