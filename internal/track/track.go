package track

import (
	"strconv"
	"time"
)

const (
	LISTENED         = "L"
	ARTIST_INDEX     = 0
	ALBUM_INDEX      = 1
	TITLE_INDEX      = 2
	RATING_INDEX     = 5
	TIMESTAMP_INDEX  = 6
	TIMESTAMP_NO_RTC = "0"
)

type Track struct {
	artist, album, title, timestamp string
}

func StringToTrack(line []string, offset int) (Track, error) {
	var (
		timestamp string = line[TIMESTAMP_INDEX]
		err       error  = nil
	)

	/* If user has a player with no Real Time Clock, the log file gives it
	a timestamp of 0. Last.fm API doesn't accept scrobbles dated that far
	into the past so in the interests of at least having the tracks sent,
	date them with current local time */
	if timestamp == TIMESTAMP_NO_RTC {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	/* Time conversion - the API wants it in UTC timezone */
	if offset != 0 {
		timestamp, err = convertTimeStamp(timestamp, offset)
	}

	track := Track{
		artist:    line[ARTIST_INDEX],
		album:     line[ALBUM_INDEX],
		title:     line[TITLE_INDEX],
		timestamp: timestamp,
	}

	return track, err

}

/* Convert back/to UTC from localtime */
func convertTimeStamp(timestamp string, offset int) (string, error) {
	timestampFlt, err := strconv.Atoi(timestamp)
	if err != nil {
		return "", err
	}

	converted := timestampFlt - offset

	return strconv.Itoa(converted), nil
}
