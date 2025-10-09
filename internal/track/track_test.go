package track

import (
	"testing"
	"time"

	"github.com/twoscott/gobble-fm/lastfm"
)

const (
	ZERO_OFFSET    = float64(0)
	TEST_TIMESTAMP = "1579643462"
)

func TestLogLineToTrack(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		offset        float64
		expectedTrack lastfm.ScrobbleParams
	}{
		{"TimelessSupport",
			[]string{"CHVRCHES", "The Bones of What You Believe", "The Mother We Share", "1", "120", "L", "0"},
			ZERO_OFFSET,
			lastfm.ScrobbleParams{
				Artist:      "CHVRCHES",
				Album:       "The Bones of What You Believe",
				Track:       "The Mother We Share",
				TrackNumber: 1,
				Duration:    lastfm.DurationSeconds(120),
				Time:        time.Now(),
				MBID:        "",
				AlbumArtist: "",
			},
		},
		{"LogLineToTrack",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "Many Men (Wish Death)", "2", "179", "L", "1579643462", "8588c220-1cab-418e-9b99-37e115755463"},
			ZERO_OFFSET,
			lastfm.ScrobbleParams{
				Artist:      "50 Cent",
				Album:       "Get Rich Or Die Tryin'",
				Track:       "Many Men (Wish Death)",
				TrackNumber: 2,
				Duration:    lastfm.DurationSeconds(179),
				Time:        time.Unix(1579643462, 0),
				MBID:        "8588c220-1cab-418e-9b99-37e115755463",
				AlbumArtist: "",
			},
		},
		{"ExtendedAudioscrobblerLog",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "Many Men (Wish Death)", "2", "179", "L", "1579643462", "8588c220-1cab-418e-9b99-37e115755463", "50 Cent"},
			ZERO_OFFSET,
			lastfm.ScrobbleParams{
				Artist:      "50 Cent",
				Album:       "Get Rich Or Die Tryin'",
				Track:       "Many Men (Wish Death)",
				TrackNumber: 2,
				Duration:    lastfm.DurationSeconds(179),
				Time:        time.Unix(1579643462, 0),
				MBID:        "8588c220-1cab-418e-9b99-37e115755463",
				AlbumArtist: "50 Cent",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotTrack, _ := StringToTrack(test.input, test.offset)
			trackEqual := true
			if !(test.expectedTrack.Artist == gotTrack.Artist) {
				trackEqual = false
			} else if !(test.expectedTrack.Album == gotTrack.Album) {
				trackEqual = false
			} else if !(test.expectedTrack.Track == gotTrack.Track) {
				trackEqual = false
			}
			if !trackEqual {
				t.Errorf("Created track was not equal to expected!\n")
			}

		})
	}

}

func TestConvertTimeStamp(t *testing.T) {
	const SECONDS_IN_HOUR = 3600
	testCases := []struct {
		name      string
		timestamp string
		offset    float64
		expected  time.Time
	}{
		{"ForwardFromUTC", TEST_TIMESTAMP, float64(+10 * SECONDS_IN_HOUR), time.Unix(1579607462, 0)},
		{"ForwardFromUTCLiteral", TEST_TIMESTAMP, float64(10 * SECONDS_IN_HOUR), time.Unix(1579607462, 0)},
		{"BackFromUTC", TEST_TIMESTAMP, float64(-10 * SECONDS_IN_HOUR), time.Unix(1579679462, 0)},
		{"ForwardFromUTCHalfHour", TEST_TIMESTAMP, float64(+0.5 * SECONDS_IN_HOUR), time.Unix(1579641662, 0)},
		{"BackFromUTCHalfHour", TEST_TIMESTAMP, float64(-0.5 * SECONDS_IN_HOUR), time.Unix(1579645262, 0)},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result, _ := convertTimeStamp(test.timestamp, test.offset)

			if result != test.expected {
				t.Errorf("Test %q failed. Expected %q, got %q\n",
					test.name, test.expected, result)
			}
		})
	}
}
