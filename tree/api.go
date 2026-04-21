package tree

import (
	"encoding/json"
	"net/http"
)

const baseApi = "https://groupietrackers.herokuapp.com/api"

func fetch(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetArtist() ([]Artist, error) {
	var artists []Artist
	err := fetch(baseApi+"/artists", &artists)
	return artists, err
}

func GetRelation() ([]Relation, error) {
	var relations RelationIndex
	err := fetch(baseApi+"/relation", &relations)
	return relations.Index, err
}

func GetDate() ([]Date, error) {
	var dates DatesIndex
	err := fetch(baseApi+"/dates", &dates)
	return dates.Index, err
}

func GetLocation() ([]Location, error) {
	var locations LocationIndex
	err := fetch(baseApi+"/locations", &locations)
	return locations.Index, err
}
