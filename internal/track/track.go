package track

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shkh/lastfm-go/lastfm"
)

const (
	SUCCESS_STR          = "[OK] %s - %s\n"
	FAIL_STR             = "[FAIL] %s - %s\n"
	SUCCESS_STR_COLOURED = "\u001b[32;1m[OK]\u001b[0m %s - %s\n"
	FAIL_STR_COLOURED    = "\u001b[31;1m[FAIL]\u001b[0m %s - %s\n"

	SEPARATOR        = '\t'
	LISTENED         = "L"
	ARTIST_INDEX     = 0
	ALBUM_INDEX      = 1
	TITLE_INDEX      = 2
	RATING_INDEX     = 5
	TIMESTAMP_INDEX  = 6
	TIMESTAMP_NO_RTC = "0"

	SECONDS_IN_HOUR = 3600
)

type Track struct {
	artist, album, title, timestamp string
}

type Tracks []Track

/* Take a string, split it, convert time if needed and return a track */
func StringToTrack(line []string, offset string) Track {
	/* Check the "RATING" index instead of looking for "\tL\t" in a line,
	just in case a track or album is named "L". If anything like this exists
	and was skipped the old method would false positive it as listened
	and then it'd be submitted */
	var timestamp string = line[TIMESTAMP_INDEX]

	/* If user has a player with no Real Time Clock, the log file gives it
	a timestamp of 0. Last.fm API doesn't accept scrobbles dated that far
	into the past so in the interests of at least having the tracks sent,
	date them with current local time */
	if timestamp == TIMESTAMP_NO_RTC {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	/* Time conversion - the API wants it in UTC timezone */
	if offset != "0h" {
		timestamp = convertTimeStamp(timestamp, offset)
	}

	track := Track{
		artist:    line[ARTIST_INDEX],
		album:     line[ALBUM_INDEX],
		title:     line[TITLE_INDEX],
		timestamp: timestamp,
	}

	return track

}

/* Convert back/to UTC from localtime */
func convertTimeStamp(timestamp, offset string) string {
	/* Log stores it in unix epoch format. Convert to an int
	so it can be manipulated with the time package */
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	/* Convert the epoch into a Time type for conversion.
	The 0 stands for 0 milliseconds since epoch which isn't
	needed */
	trackTime := time.Unix(timestampInt, 0)

	/* Take the offset flag and convert it into the duration
	which will be added/subtracted */
	newOffset, err := time.ParseDuration(offset)
	if err != nil {
		log.Fatal(err)
	}

	/* The duration is negative so that entries behind UTC are
	brought forward while entries ahead are brought back */
	convertedTime := trackTime.Add(-newOffset)
	return strconv.FormatInt(convertedTime.Unix(), 10)
}

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
