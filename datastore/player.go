package datastore

// DsPlayer kind
type DsPlayer struct {
	ValintyURL int `datastore:"vality_url"`
	Name       string
	Country    string
}
