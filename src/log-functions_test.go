package main

import "testing"

const (
	ZERO_OFFSET    = "0h"
	TEST_TIMESTAMP = "1579643462"
)

func TestLogLineToTrackSkip(t *testing.T) {
	line := "50 Cent	Get Rich Or Die Tryin'	In Da Club	2	179	S	1579643462"
	_, bool := logLineToTrack(line, ZERO_OFFSET)

	if bool {
		t.Errorf("logLineToTrack Should Have returned false due to skipped file!")
	}
}

func TestLogLineToTrack(t *testing.T) {
	expecting := Track{
		artist:    "50 Cent",
		album:     "Get Rich Or Die Tryin'",
		title:     "Many Men (Wish Death)",
		timestamp: "1579643462",
	}

	line := "50 Cent	Get Rich Or Die Tryin'	Many Men (Wish Death)	2	179	L	1579643462"

	gotTrack, result := logLineToTrack(line, ZERO_OFFSET)

	if !result {
		t.Errorf("Track was listened! logLineToTrack should have returned true")
	}

	structEqual := true

	if !(expecting.artist == gotTrack.artist) {
		structEqual = false
	} else if !(expecting.album == gotTrack.album) {
		structEqual = false
	} else if !(expecting.title == gotTrack.title) {
		structEqual = false
	} else if !(expecting.timestamp == gotTrack.timestamp) {
		structEqual = false
	}

	if !structEqual {
		t.Errorf("Track object did not match expected!")
	}

}

func TestConvertTimeStamp(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp string
		offset    string
		expected  string
	}{
		{"ForwardFromUTC", TEST_TIMESTAMP, "+10h", "1579607462"},
		{"ForwardFromUTCLiteral", TEST_TIMESTAMP, "10h", "1579607462"},
		{"BackFromUTC", TEST_TIMESTAMP, "-10h", "1579679462"},
		{"ForwardFromUTCHalfHour", TEST_TIMESTAMP, "+0.5h", "1579641662"},
		{"BackFromUTCHalfHour", TEST_TIMESTAMP, "-0.5h", "1579645262"},
		{"MinuteInput", TEST_TIMESTAMP, "-30m", "1579645262"},
	}

	for _, test := range testCases {
		result := convertTimeStamp(test.timestamp, test.offset)

		if result != test.expected {
			t.Errorf("Test %q failed. Expected %q, got %q\n", test.name, test.expected, result)
		}
	}
}
