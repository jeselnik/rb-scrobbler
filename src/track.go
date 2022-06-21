package main

import (
	"fmt"

	"github.com/shkh/lastfm-go/lastfm"
)

const (
	SUCCESS_STR          = "[OK] %s - %s\n"
	FAIL_STR             = "[FAIL] %s - %s\n"
	SUCCESS_STR_COLOURED = "\u001b[32;1m[OK]\u001b[0m %s - %s\n"
	FAIL_STR_COLOURED    = "\u001b[31;1m[FAIL]\u001b[0m %s - %s\n"
)

type Track struct {
	artist    string
	album     string
	title     string
	timestamp string
}

type Tracks []Track

func (tracks Tracks) scrobble(api *lastfm.Api, colours *bool) (
	success, fail uint) {

	var (
		successString string
		failString    string
	)

	if *colours {
		successString = SUCCESS_STR_COLOURED
		failString = FAIL_STR_COLOURED
	} else {
		successString = SUCCESS_STR
		failString = FAIL_STR
	}

	for _, track := range tracks {
		p := lastfm.P{"artist": track.artist, "album": track.album,
			"track": track.title, "timestamp": track.timestamp}

		res, err := api.Track.Scrobble(p)
		if err != nil || res.Ignored != "0" {
			fmt.Printf(failString, track.artist, track.title)
			fail++
		} else {
			fmt.Printf(successString, track.artist, track.title)
			success++
		}
	}

	return
}
