package track

import (
	"fmt"
	"strings"

	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/gobble-fm/session"
)

const (
	API_TRACK_SUB_LIMIT = 50
	GREEN               = "\u001b[32;1m"
	RED                 = "\u001b[31;1m"
	CLEAR               = "\u001b[0m"
)

func PrintResult(success bool, colours *bool, artist string, trackName string) {
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

	msg.WriteString(artist)
	msg.WriteString(" - ")
	msg.WriteString(trackName)

	fmt.Println(msg.String())
}

func Scrobble(api *session.Client, tracks lastfm.ScrobbleMultiParams, colours *bool) (
	success, fail int) {

	for i := 0; i < len(tracks); i += API_TRACK_SUB_LIMIT {
		end := i + API_TRACK_SUB_LIMIT
		if end > len(tracks) {
			end = len(tracks)
		}
		batch := tracks[i:end]

		res, err := api.Track.ScrobbleMulti(batch)
		if err != nil {
			fail += len(tracks)
			for _, scr := range tracks {
				PrintResult(false, colours, scr.Artist, scr.Track)
			}
			continue
		}

		success += res.Accepted
		fail += res.Ignored

		for _, scr := range res.Scrobbles {
			success := true
			if scr.Ignored.Code != lastfm.ScrobbleNotIgnored {
				success = false
			}
			PrintResult(success, colours, scr.Artist.Name, scr.Track.Title)
		}
	}

	return
}
