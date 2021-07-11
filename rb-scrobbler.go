package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

const (
	AUDIOSCROBBLER_HEADER  = "#AUDIOSCROBBLER/"
	LISTENED               = "\tL\t"
	SEPARATOR              = "\t"
	FIRST_TRACK_LINE_INDEX = 3
	ARTIST_INDEX           = 0
	ALBUM_INDEX            = 1
	TITLE_INDEX            = 2
	TIMESTAMP_INDEX        = 6
)

type Track struct {
	artist    string
	album     string
	title     string
	timestamp string
}

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	offset := flag.String("o", "0h", "Offset from UTC")
	flag.Parse()

	scrobblerLog, err := importLog(logPath)
	if err != nil {
		log.Fatal(err)
	}

	var tracks []Track
	for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
		if strings.Contains(scrobblerLog[i], LISTENED) {
			tracks = append(tracks, logLineToTrack(scrobblerLog[i], *offset))
		}
	}

	for i := 0; i < len(tracks); i++ {
		fmt.Println(tracks[i].timestamp)
	}
}
