package track

import (
	"strconv"
	"testing"
	"time"

	"github.com/twoscott/gobble-fm/lastfm"
)

const (
	ZERO_OFFSET    = 0
	TEST_TIMESTAMP = "1579643462"
)

func TestLogLineToTrack(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		offset        int
		expectedTrack lastfm.ScrobbleParams
	}{
		{"SkipTrack",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "In Da Club", "2", "179", "S", "1579643462"},
			ZERO_OFFSET,
			lastfm.ScrobbleParams{
				Artist:      "50 Cent",
				Album:       "Get Rich Or Die Tryin'",
				Track:       "In Da Club",
				TrackNumber: 2,
				Duration:    lastfm.DurationSeconds(179),
				timestamp:   "1579643462",
				MBID:        "",
				AlbumArtist: "",
			},
		},
		{"TimelessSupport",
			[]string{"CHVRCHES", "The Bones of What You Believe", "The Mother We Share", "1", "120", "L", "0"},
			ZERO_OFFSET,
			Track{
				artist:      "CHVRCHES",
				album:       "The Bones of What You Believe",
				title:       "The Mother We Share",
				position:    "1",
				duration:    "120",
				rating:      "L",
				timestamp:   strconv.FormatInt(time.Now().Unix(), 10),
				mbid:        "",
				albumArtist: "",
			},
		},
		{"LogLineToTrack",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "Many Men (Wish Death)", "2", "179", "L", "1579643462", "8588c220-1cab-418e-9b99-37e115755463"},
			ZERO_OFFSET,
			Track{
				artist:      "50 Cent",
				album:       "Get Rich Or Die Tryin'",
				title:       "Many Men (Wish Death)",
				position:    "2",
				duration:    "179",
				rating:      "L",
				timestamp:   "1579643462",
				mbid:        "8588c220-1cab-418e-9b99-37e115755463",
				albumArtist: "",
			},
		},
		{"ExtendedAudioscrobblerLog",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "Many Men (Wish Death)", "2", "179", "L", "1579643462", "8588c220-1cab-418e-9b99-37e115755463", "50 Cent"},
			ZERO_OFFSET,
			Track{
				artist:      "50 Cent",
				album:       "Get Rich Or Die Tryin'",
				title:       "Many Men (Wish Death)",
				position:    "2",
				duration:    "179",
				rating:      "L",
				timestamp:   "1579643462",
				mbid:        "8588c220-1cab-418e-9b99-37e115755463",
				albumArtist: "50 Cent",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotTrack, _ := StringToTrack(test.input, test.offset)
			trackEqual := true
			if !(test.expectedTrack.artist == gotTrack.artist) {
				trackEqual = false
			} else if !(test.expectedTrack.album == gotTrack.album) {
				trackEqual = false
			} else if !(test.expectedTrack.title == gotTrack.title) {
				trackEqual = false
			} else if !(test.expectedTrack.timestamp ==
				gotTrack.timestamp) {
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
		offset    int
		expected  string
	}{
		{"ForwardFromUTC", TEST_TIMESTAMP, int(+10 * SECONDS_IN_HOUR), "1579607462"},
		{"ForwardFromUTCLiteral", TEST_TIMESTAMP, int(10 * SECONDS_IN_HOUR), "1579607462"},
		{"BackFromUTC", TEST_TIMESTAMP, int(-10 * SECONDS_IN_HOUR), "1579679462"},
		{"ForwardFromUTCHalfHour", TEST_TIMESTAMP, int(+0.5 * SECONDS_IN_HOUR), "1579641662"},
		{"BackFromUTCHalfHour", TEST_TIMESTAMP, int(-0.5 * SECONDS_IN_HOUR), "1579645262"},
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
