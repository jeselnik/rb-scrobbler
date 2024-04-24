package track

import (
	"fmt"

	"github.com/sonjek/go-lastfm/lastfm"
)

const (
	SUCCESS_STR = "%s[OK]%s %s - %s\n"
	FAIL_STR    = "%s[FAIL]%s %s - %s\n"

	GREEN = "\u001b[32;1m"
	RED   = "\u001b[31;1m"
	CLEAR = "\u001b[0m"
)

func PrintResult(success bool, colours *bool, track Track) {
	var (
		pattern, colour string
	)
	clear := CLEAR

	if success {
		pattern = SUCCESS_STR
		colour = GREEN
	} else {
		pattern = FAIL_STR
		colour = RED
	}

	if !(*colours) {
		colour = ""
		clear = ""
	}

	fmt.Printf(pattern, colour, clear, track.artist, track.title)
}

func Scrobble(api *lastfm.Api, tracks []Track, colours *bool) (
	success, fail uint) {

	for _, track := range tracks {
		p := lastfm.P{"artist": track.artist, "album": track.album,
			"track": track.title, "timestamp": track.timestamp}

		isSuccess := false

		res, err := api.Track.Scrobble(p)
		if err != nil || res.Ignored != "0" {
			fail++
		} else {
			isSuccess = true
			success++
		}

		PrintResult(isSuccess, colours, track)
	}

	return
}
