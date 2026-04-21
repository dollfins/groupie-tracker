package tree

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func convertDates(dateStr string) int64 {
	layout := "02-01-2006"
	date, _ := time.Parse(layout, dateStr)
	return date.Unix()

}
func geocode(address string) (LocationDetail, error) {
	apiURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
		url.QueryEscape(address),
	)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return LocationDetail{}, err
	}

	req.Header.Set("User-Agent", "MyGoMapApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return LocationDetail{}, err
	}

	if resp == nil {
		return LocationDetail{}, fmt.Errorf("empty response from server")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationDetail{}, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var results []LocationDetail
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return LocationDetail{}, err
	}

	if len(results) == 0 {
		return LocationDetail{}, fmt.Errorf("no results found for address: %s", address)
	}

	return results[0], nil
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	opt := r.URL.Query().Get("filters")

	artists := CacheArtist
	if search != "" {
		var filter []Artist
		for _, artist := range artists {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(search)) {
				filter = append(filter, artist)
			}
		}
		artists = filter
	}
	data := struct {
		Artist []Artist
		Opt    string
	}{
		Artist: artists,
		Opt:    opt,
	}
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		return
	}
	tmpl.Execute(w, data)
}
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filterOpt := r.URL.Query().Get("filters")

	artist := CacheArtist

	switch filterOpt {
	case "creationDate":
		for i := 0; i < len(artist)-1; i++ {
			for j := 0; j < len(artist)-i-1; j++ {
				if artist[j].CreationDate > artist[j+1].CreationDate {
					artist[j], artist[j+1] = artist[j+1], artist[j]
				}
			}
		}
	case "firstAlbum":
		for i := 0; i < len(artist)-1; i++ {
			for j := 0; j < len(artist)-i-1; j++ {
				if convertDates(artist[j].FirstAlbum) > convertDates(artist[j+1].FirstAlbum) {
					artist[j], artist[j+1] = artist[j+1], artist[j]
				}
			}
		}
	case "members":
		for i := 0; i < len(artist)-1; i++ {
			for j := 0; j < len(artist)-i-1; j++ {
				if len(artist[j].Members) > len(artist[j+1].Members) {
					artist[j], artist[j+1] = artist[j+1], artist[j]
				}
			}
		}
	}

	data := struct {
		Artist []Artist
		Opt    string
	}{
		Artist: artist,
		Opt:    filterOpt,
	}
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		return
	}
	tmpl.Execute(w, data)
}
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Number", 500)
	}

	artists := CacheArtist
	relations := CacheRelation
	dates := CacheDate
	locations := CacheLocation

	var artist Artist
	var relation Relation
	var location Location
	var date Date

	var addresses []string

	for _, a := range artists {
		if a.Id == idNum {
			artist = a
			break
		}
	}

	for _, r := range relations {
		if r.Id == idNum {
			relation = r
			break
		}
	}

	for _, l := range locations {
		if l.Id == idNum {
			location = l
			addresses = l.Locations
			break
		}
	}

	for _, d := range dates {
		if d.Id == idNum {
			date = d
			break
		}
	}

	var locationMap []LocationDetail

	for _, addr := range addresses {
		cleanAddr := strings.ReplaceAll(addr, "_", " ")
		cleanAddr = strings.ReplaceAll(cleanAddr, "-", ", ")

		loc, err := geocode(cleanAddr)
		if err != nil {
			fmt.Println("Geocode error:", err)
			continue
		}
		locationMap = append(locationMap, loc)
	}

	jsonData, err := json.Marshal(locationMap)
	if err != nil {
		http.Error(w, "JSON error", 500)
		return
	}

	output := struct {
		Artist   Artist
		Relation Relation
		Date     Date
		Location Location
		Data     template.JS
	}{
		Artist:   artist,
		Relation: relation,
		Date:     date,
		Location: location,
		Data:     template.JS(jsonData),
	}

	tmpl, err := template.ParseFiles("template/detail.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		return
	}

	tmpl.Execute(w, output)
}
