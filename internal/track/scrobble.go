package track

import (
	"fmt"

	"github.com/sonjek/go-lastfm/lastfm"
)

func Scrobble(api *lastfm.Api, tracks []Track, colours *bool) (
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
