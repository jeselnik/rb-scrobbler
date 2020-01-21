use std::fs;
use std::time::UNIX_EPOCH;
use std::vec;
extern crate chrono;
use chrono::Utc;
extern crate rustfm_scrobble;
extern crate time;
use rustfm_scrobble::Scrobble;

const AUDIOSCROBBLER_HEADER: &str = "#AUDIOSCROBBLER/";
const LISTENED: &str = "\tL\t";
const SEPARATOR: &str = "\t";

/* Indexes used when iterating over track data */
const ARTIST_INDEX: usize = 0;
const ALBUM_INDEX: usize = 1;
const TITLE_INDEX: usize = 2;
const TIMESTAMP_INDEX: usize = 6;

pub fn as_scrobbles(path: &str, offset: i64) -> vec::Vec<Scrobble> {
    let file = fs::read_to_string(path).expect("Error Opening File");

    if !file.contains(AUDIOSCROBBLER_HEADER) {
        panic!("Not a valid .scrobbler.log!");
    } else {
        let mut scrobbles = Vec::new();
        let file_as_lines = file.lines();

        for line in file_as_lines {
            if line.contains(LISTENED) {
                let track_str: vec::Vec<&str> = line.split(SEPARATOR).collect();

                let mut scrobble = Scrobble::new(
                    track_str[ARTIST_INDEX],
                    track_str[TITLE_INDEX],
                    track_str[ALBUM_INDEX],
                );

                if offset != 0 {
                    scrobble.with_timestamp(convert_time(offset, track_str[TIMESTAMP_INDEX]));
                } else {
                    let timestamp: u64 = track_str[TIMESTAMP_INDEX].parse().unwrap();
                    scrobble.with_timestamp(timestamp);
                }
                scrobbles.push(scrobble);
            }
        }
        return scrobbles;
    }
}

fn convert_time(offset: i64, timestamp: &str) -> u64 {
    let casted_timestamp: u64 = timestamp.parse().unwrap();
    let raw_time = UNIX_EPOCH + std::time::Duration::from_secs(casted_timestamp);
    let in_datetime = chrono::prelude::DateTime::<Utc>::from(raw_time);
    let converted = in_datetime.checked_sub_signed(time::Duration::minutes(offset));
    return converted.unwrap().timestamp() as u64;
}
