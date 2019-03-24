package main

// run: alfred_workflow_data=workflow alfred_workflow_cache=/tmp/alfred alfred_workflow_bundleid=mk_testing go run alfred-deezer.go track mötley

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	resty "gopkg.in/resty.v1"
)

// TrackResult represents a result of a Deezer track query
type TrackResult struct {
	// Inspired by https://medium.com/@IndianGuru/consuming-json-apis-with-go-d711efc1dcf9
	Data []struct {
		ID             int    `json:"id"`
		Readable       bool   `json:"readable"`
		Title          string `json:"title"`
		TitleShort     string `json:"title_short"`
		TitleVersion   string `json:"title_version"`
		Link           string `json:"link"`
		Duration       int    `json:"duration"`
		Rank           int    `json:"rank"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Preview        string `json:"preview"`
		Artist         struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Album struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Cover       string `json:"cover"`
			CoverSmall  string `json:"cover_small"`
			CoverMedium string `json:"cover_medium"`
			CoverBig    string `json:"cover_big"`
			CoverXl     string `json:"cover_xl"`
			Tracklist   string `json:"tracklist"`
			Type        string `json:"type"`
		} `json:"album"`
		Type string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

// AlbumResult represents a result of a Deezer album query
type AlbumResult struct {
	Data []struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		Link           string `json:"link"`
		Cover          string `json:"cover"`
		CoverSmall     string `json:"cover_small"`
		CoverMedium    string `json:"cover_medium"`
		CoverBig       string `json:"cover_big"`
		CoverXl        string `json:"cover_xl"`
		GenreID        int    `json:"genre_id"`
		NbTracks       int    `json:"nb_tracks"`
		RecordType     string `json:"record_type"`
		Tracklist      string `json:"tracklist"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Artist         struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Type string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

// aw.Workflow is the main API
var wf *aw.Workflow

func init() {
	// Create a new *Workflow using default configuration
	// (workflow settings are read from the environment variables
	// set by Alfred)
	wf = aw.New()
}

func main() {
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying
	wf.Run(run)
}

func queryDeezer(query string, resource string) string {
	// https://api.deezer.com/search/track?q="$query"&limit=1&order=RANKING_DESC"
	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"q":     query,
			"limit": "10",
			"order": "RANKING_DESC",
		}).
		Get("https://api.deezer.com/search/" + resource)

	if err == nil {
	}

	return resp.String()
}

func run() {
	contentType := os.Args[1]
	title := os.Args[2]

	switch contentType {
	case "track":
		runTracks(title)
	case "album":
		runAlbum(title)
	case "artist":
		runArtist(title)
	}

}

func getLocalURL(url string) string {
	return strings.Replace(url, "https://", "deezer://", -1)
}

func runAlbum(title string) {
	response := queryDeezer(title, "album")

	var albums AlbumResult
	if err := json.NewDecoder(strings.NewReader(response)).Decode(&albums); err != nil {
		// log.Println(err)
	}

	for _, album := range albums.Data {
		var icon aw.Icon
		icon.Value = album.CoverSmall

		id := strconv.Itoa(album.ID)
		url := "https://www.deezer.com/en/album/" + id

		wf.NewItem(album.Artist.Name + " - " + album.Title).
			// Subtitle(album.Album.Title).
			Valid(true).
			// Icon(&icon).
			Arg(url).
			Quicklook(url).
			UID("album" + id).
			NewModifier("cmd").
			Subtitle("Open in Deezer App").
			Arg(getLocalURL(url))
	}

	// And send the results to Alfred
	wf.SendFeedback()
}
func runArtist(title string) {
	// TODO implement me
}

func runTracks(title string) {

	response := queryDeezer(title, "track")

	var tracks TrackResult
	if err := json.NewDecoder(strings.NewReader(response)).Decode(&tracks); err != nil {
		// log.Println(err)
	}

	for _, track := range tracks.Data {
		var icon aw.Icon
		icon.Value = track.Album.CoverSmall

		id := strconv.Itoa(track.ID)
		url := "https://www.deezer.com/en/track/" + id

		wf.NewItem(track.Artist.Name + " -  " + track.Title).
			Subtitle(track.Album.Title).
			Valid(true).
			// Icon(&icon).
			Arg(url).
			Quicklook(url).
			UID("track" + id).
			NewModifier("cmd").
			Subtitle("Open in Deezer App").
			Arg(getLocalURL(url))
	}

	// And send the results to Alfred
	wf.SendFeedback()
}
