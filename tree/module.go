package tree

var CacheArtist []Artist
var CacheRelation []Relation
var CacheDate []Date
var CacheLocation []Location

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Date         string   `json:"concertDates"`
	Location     string   `json:"locations"`
	Relations    string   `json:"relations"`
}

type RelationIndex struct {
	Index []Relation `json:"index"`
}

type LocationIndex struct {
	Index []Location `json:"index"`
}

type DatesIndex struct {
	Index []Date `json:"index"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Date      []Date
}

type UserProfile struct {
	Artist []Artist
}

type ProfileAPI struct {
	Artist   string `json:"artists"`
	Location string `json:"locations"`
	Dates    string `json:"dates"`
	Relation string `json:"relation"`
}

type LocationDetail struct {
	Name string  `json:"display_name"`
	Lat  float64 `json:"lat,string"`
	Lon  float64 `json:"lon,string"`
}
