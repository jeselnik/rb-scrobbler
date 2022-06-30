package track

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
	Artist, Album, Title, Timestamp string
}

type Tracks []Track

func Scrobble(api *lastfm.Api, tracks Tracks, colours *bool) (
	success, fail uint) {

	var successString, failString string

	if *colours {
		successString = SUCCESS_STR_COLOURED
		failString = FAIL_STR_COLOURED
	} else {
		successString = SUCCESS_STR
		failString = FAIL_STR
	}

	for _, track := range tracks {
		p := lastfm.P{"artist": track.Artist, "album": track.Album,
			"track": track.Title, "timestamp": track.Timestamp}

		res, err := api.Track.Scrobble(p)
		if err != nil || res.Ignored != "0" {
			fmt.Printf(failString, track.Artist, track.Title)
			fail++
		} else {
			fmt.Printf(successString, track.Artist, track.Title)
			success++
		}
	}

	return
}
