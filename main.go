package main

import (
	"fmt"
	"html/template"
	"log"
	"ltvco/http_handler"
	"ltvco/repository"
	"net/http"
	"sync"
)

var instance *http_handler.HttpInstance

type Args struct {
	From   string
	Until  string
	Artist string
}

func getReleases(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query()["from"]
	to := r.URL.Query()["to"]
	artist := r.URL.Query()["artist"]

	cachedRecords, uncachedRecords, err := instance.GetRecords(from[0], to[0], artist[0])

	if err != nil {
		fmt.Fprintf(w, "<script>alert('There was an error: "+err.Error()+"')</script>")
	}

	html := fmt.Sprintf("<h1>RECORDS FOUND: %d | CACHED: %d | UNCACHED: %d</h1>", (len(cachedRecords) + len(uncachedRecords)), len(cachedRecords), len(uncachedRecords))
	html += "<table border='1'>"
	html += "<tr><td></td><td><b>Id</b></td><td><b>Released At</b></td><td><b>Duration</b></td><td><b>Artist</b></td><td><b>Name</b></td><td><b>Stats</b></td>"

	var wg sync.WaitGroup
	wg.Add(2)

	go func(html *string, cachedRecords *[]repository.SongResponse, wg *sync.WaitGroup) {
		defer wg.Done()

		for _, record := range *cachedRecords {
			stats := fmt.Sprintf("Last_played_at: %b</br>Times_played: %d</br>Global_rank: %d</br>", record.Stats.Last_played_at, record.Stats.Times_played, record.Stats.Global_rank)

			*html += fmt.Sprintf("<tr><td><b>Cached</b></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>", record.Song_id, record.Released_at, record.Duration, record.Artist, record.Name, stats)
		}
	}(&html, &cachedRecords, &wg)

	go func(html *string, cachedRecords *[]repository.SongResponse, wg *sync.WaitGroup) {
		defer wg.Done()

		for _, record := range uncachedRecords {
			stats := fmt.Sprintf("Last_played_at: %b</br>Times_played: %d</br>Global_rank: %d</br>", record.Stats.Last_played_at, record.Stats.Times_played, record.Stats.Global_rank)

			*html += fmt.Sprintf("<tr><td><b>Uncached</b></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>", record.Song_id, record.Released_at, record.Duration, record.Artist, record.Name, stats)
		}
	}(&html, &uncachedRecords, &wg)

	wg.Wait()

	html += "</br><button onclick='window.history.back()'><< Go Back</button></br>"
	html += "</html>"

	fmt.Fprintf(w, html)
}

func loadMain(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("ui/main.html")
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", loadMain)
	http.HandleFunc("/releases", getReleases)

	instance = http_handler.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server started at port 8080")

	select {}
}
