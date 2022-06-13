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

func (tracks Tracks) scrobble(api *lastfm.Api, noColours *bool) (uint, uint) {
	var success uint = 0
	var fail uint = 0

	var successString string
	var failString string

	if *noColours {
		successString = SUCCESS_STR
		failString = FAIL_STR
	} else {
		successString = SUCCESS_STR_COLOURED
		failString = FAIL_STR_COLOURED
	}

	for _, track := range tracks {
		p := lastfm.P{"artist": track.artist, "album": track.album, "track": track.title, "timestamp": track.timestamp}

		_, err := api.Track.Scrobble(p)
		if err != nil {
			fmt.Printf(failString, track.artist, track.title)
			fail++
		} else {
			fmt.Printf(successString, track.artist, track.title)
			success++
		}
	}

	return success, fail
}
