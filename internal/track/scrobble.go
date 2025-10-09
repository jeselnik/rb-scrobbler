package track

import (
	"fmt"
	"strings"

	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/gobble-fm/session"
)

const (
	GREEN = "\u001b[32;1m"
	RED   = "\u001b[31;1m"
	CLEAR = "\u001b[0m"
)

func PrintResult(success bool, colours *bool, track lastfm.ScrobbleParams) {
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

	msg.WriteString(track.Artist)
	msg.WriteString(" - ")
	msg.WriteString(track.Track)

	fmt.Println(msg.String())
}

func Scrobble(api *session.Client, tracks lastfm.ScrobbleMultiParams, colours *bool) (
	success, fail uint) {

	for _, track := range tracks {
		isSuccess := false

		res, err := api.Track.Scrobble(track)
		if err != nil {
			fail++
			PrintResult(isSuccess, colours, track)
		}

		if res.Scrobble.Ignored.Code == lastfm.ScrobbleNotIgnored {
			isSuccess = true
			success++
		} else {
			fail++
		}

		PrintResult(isSuccess, colours, track)
	}

	return
}
