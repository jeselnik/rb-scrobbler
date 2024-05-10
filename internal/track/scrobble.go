package track

import (
	"fmt"
	"strings"

	"github.com/sonjek/go-lastfm/lastfm"
)

const (
	GREEN = "\u001b[32;1m"
	RED   = "\u001b[31;1m"
	CLEAR = "\u001b[0m"
)

func PrintResult(success bool, colours *bool, track Track) {
	var msg strings.Builder

	if success && *colours {
		msg.WriteString(GREEN)
	} else if !success && *colours {
		msg.WriteString(RED)
	}

	if success {
		msg.WriteString("[OK] ")
	} else {
		msg.WriteString("[FAIL] ")
	}

	if *colours {
		msg.WriteString(CLEAR)
	}

	msg.WriteString(track.artist)
	msg.WriteString(" - ")
	msg.WriteString(track.title)

	fmt.Println(msg.String())
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
