package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	resty "gopkg.in/resty.v1"
)

type Track struct {
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

func queryDeezer(query string) Track {
	// https://api.deezer.com/search/track?q="$query"&limit=1&order=RANKING_DESC"
	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"q":     query,
			"limit": "10",
			"order": "RANKING_DESC",
		}).
		Get("https://api.deezer.com/search/track")

	if err == nil {
	}

	var track Track
	if err := json.NewDecoder(strings.NewReader(resp.String())).Decode(&track); err != nil {
		// log.Println(err)
	}

	return track
}

func run() {
	title := os.Args[1]
	tracks := queryDeezer(title)

	for _, track := range tracks.Data {
		var icon aw.Icon
		icon.Value = track.Album.CoverSmall

		id := strconv.Itoa(track.ID)
		url := "https://www.deezer.com/en/track/" + id

		wf.NewItem(track.Artist.Name + " -  " + track.Title).
			Subtitle(track.Album.Title).
			Valid(true).
			Icon(&icon).
			Arg(id).
			Quicklook(url).
			UID("track" + id)
	}

	// And send the results to Alfred
	wf.SendFeedback()
}
