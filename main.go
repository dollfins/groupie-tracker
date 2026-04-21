package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tree/tree"
)

func loadData() error {
	var err error

	tree.CacheArtist, err = tree.GetArtist()
	if err != nil {
		return err
	}
	tree.CacheRelation, err = tree.GetRelation()
	if err != nil {
		return err
	}
	tree.CacheLocation, err = tree.GetLocation()
	if err != nil {
		return err
	}
	tree.CacheDate, err = tree.GetDate()
	if err != nil {
		return err
	}
	return nil
}

func autoRefresh(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			artists, err1 := tree.GetArtist()
			relation, err2 := tree.GetRelation()
			date, err3 := tree.GetDate()
			location, err4 := tree.GetLocation()

			log.Println("✅ Cache Data Successfully Updated")

			if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
				tree.CacheArtist = artists
				tree.CacheRelation = relation
				tree.CacheLocation = location
				tree.CacheDate = date
			}
		}
	}()
}
func main() {
	loadData()
	autoRefresh(10 * time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", tree.HomeHandler)
	mux.HandleFunc("/details", tree.DetailHandler)
	mux.HandleFunc("/filter", tree.FilterHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server Running on http://localhost:8080/")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
