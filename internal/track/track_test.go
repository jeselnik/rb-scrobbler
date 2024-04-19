package track

import (
	"strconv"
	"testing"
	"time"
)

const (
	ZERO_OFFSET    = 0
	TEST_TIMESTAMP = "1579643462"
)

func TestLogLineToTrack(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		offset        float64
		expectedTrack Track
	}{
		{"SkipTrack",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "In Da Club", "2", "179", "S", "1579643462"},
			ZERO_OFFSET,
			Track{
				artist:    "50 Cent",
				album:     "Get Rich Or Die Tryin'",
				title:     "In Da Club",
				timestamp: "1579643462"},
		},
		{"TimelessSupport",
			[]string{"CHVRCHES", "The Bones of What You Believe", "The Mother We Share", "1", "120", "L", "0"},
			ZERO_OFFSET,
			Track{
				artist:    "CHVRCHES",
				album:     "The Bones of What You Believe",
				title:     "The Mother We Share",
				timestamp: strconv.FormatInt(time.Now().Unix(), 10)},
		},
		{"LogLineToTrack",
			[]string{"50 Cent", "Get Rich Or Die Tryin'", "Many Men (Wish Death)", "2", "179", "L", "1579643462"},
			ZERO_OFFSET,
			Track{
				artist:    "50 Cent",
				album:     "Get Rich Or Die Tryin'",
				title:     "Many Men (Wish Death)",
				timestamp: "1579643462"},
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
	testCases := []struct {
		name      string
		timestamp string
		offset    float64
		expected  string
	}{
		{"ForwardFromUTC", TEST_TIMESTAMP, +10, "1579607462"},
		{"ForwardFromUTCLiteral", TEST_TIMESTAMP, 10, "1579607462"},
		{"BackFromUTC", TEST_TIMESTAMP, -10, "1579679462"},
		{"ForwardFromUTCHalfHour", TEST_TIMESTAMP, +0.5, "1579641662"},
		{"BackFromUTCHalfHour", TEST_TIMESTAMP, -0.5, "1579645262"},
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
