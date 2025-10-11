package track

import (
	"strconv"
	"time"

	"github.com/twoscott/gobble-fm/lastfm"
)

const (
	LISTENED         = "L"
	TIMESTAMP_NO_RTC = "0"
)

const (
	ARTIST_INDEX = iota
	ALBUM_INDEX
	TITLE_INDEX
	POSITION_INDEX
	DURATION_INDEX
	RATING_INDEX
	TIMESTAMP_INDEX
	MBID_INDEX
	ALBUM_ARTIST_INDEX
)

func StringToTrack(line []string, offset float64) (lastfm.ScrobbleParams, error) {
	var (
		albumArtist  string
		duration     int
		mbid         string
		position     int
		timestamp    time.Time
		durationRaw  string = line[DURATION_INDEX]
		positionRaw  string = line[POSITION_INDEX]
		timestampRaw string = line[TIMESTAMP_INDEX]
		err          error  = nil
	)

	position, err = strconv.Atoi(positionRaw)
	if err != nil {
		return lastfm.ScrobbleParams{}, err
	}

	duration, err = strconv.Atoi(durationRaw)
	if err != nil {
		return lastfm.ScrobbleParams{}, err
	}

	/* If the user has a player with no Real Time Clock, the log file gives it
	a timestamp of 0. The last.fm API doesn't accept scrobbles dated that far
	into the past so in the interests of at least having the tracks sent,
	date them with the current local time */
	if timestampRaw == TIMESTAMP_NO_RTC {
		timestamp = time.Now().UTC()
	}

	/* Time conversion - the API wants it in UTC timezone */
	if offset != 0 {
		timestamp, err = convertTimeStamp(timestampRaw, offset)
	}

	if len(line) > 7 {
		mbid = line[MBID_INDEX]
	}

	if len(line) > 8 {
		albumArtist = line[ALBUM_ARTIST_INDEX]
	}

	track := lastfm.ScrobbleParams{
		Artist:      line[ARTIST_INDEX],
		Track:       line[TITLE_INDEX],
		Time:        timestamp,
		Album:       line[ALBUM_INDEX],
		AlbumArtist: albumArtist,
		TrackNumber: position,
		Duration:    lastfm.DurationSeconds(duration),
		MBID:        mbid,
	}

	return track, err
}

/* Convert back/to UTC from localtime */
func convertTimeStamp(timestamp string, offset float64) (time.Time, error) {
	timestampFlt, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		return time.Time{}, err
	}

	converted := timestampFlt - offset
	return time.Unix(int64(converted), 0), nil
}
