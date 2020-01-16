use std::fs;
use std::string;
use std::time::UNIX_EPOCH;
use std::vec;
extern crate chrono;
use chrono::Utc;
extern crate time;

const AUDIOSCROBBLER_HEADER: &str = "#AUDIOSCROBBLER/";
const LISTENED: &str = "\tL\t";
const SEPARATOR: &str = "\t";

/* Track data begins after this line in the log */
const TRACKS_BEGIN_INDEX: usize = 3;
/* Indexes used when iterating over track data */
const ARTIST_INDEX: usize = 0;
const ALBUM_INDEX: usize = 1;
const TITLE_INDEX: usize = 2;
const TIMESTAMP_INDEX: usize = 6;

pub struct Track {
    pub artist: String,
    pub album: String,
    pub title: String,
    pub timestamp: u64,
}

pub fn as_vec(path: &str) -> vec::Vec<Track> {
    let file = fs::read_to_string(path).expect("Error Opening File");

    if !file.contains(AUDIOSCROBBLER_HEADER) {
        panic!("Not a valid .scrobbler.log!");
    } else {
        let mut tracks = vec::Vec::new();
        let file_as_lines: Vec<&str> = file.lines().collect();

        for index in TRACKS_BEGIN_INDEX..file_as_lines.len() {
            if file_as_lines[index].contains(LISTENED) {
                let track_str: Vec<&str> = file_as_lines[index].split(SEPARATOR).collect();

                let track = Track {
                    artist: string::String::from(track_str[ARTIST_INDEX]),
                    album: string::String::from(track_str[ALBUM_INDEX]),
                    title: string::String::from(track_str[TITLE_INDEX]),
                    timestamp: track_str[TIMESTAMP_INDEX].parse().unwrap(),
                };

                &tracks.push(track);
            }
        }
        return tracks;
    }
}

impl Track {
    pub fn convert_time(&mut self, offset: i64) {
        let raw_time = UNIX_EPOCH + std::time::Duration::from_secs(self.timestamp);

        let in_datetime = chrono::prelude::DateTime::<Utc>::from(raw_time);

        let converted = in_datetime.checked_sub_signed(time::Duration::minutes(offset));

        let in_unix: i64 = converted.unwrap().timestamp();

        self.timestamp = in_unix as u64;
    }
}
