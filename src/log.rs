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

/* Track data begins after this line in the log */
const TRACKS_BEGIN_INDEX: usize = 3;
/* Indexes used when iterating over track data */
const ARTIST_INDEX: usize = 0;
const ALBUM_INDEX: usize = 1;
const TITLE_INDEX: usize = 2;
const TIMESTAMP_INDEX: usize = 6;

pub fn as_vec(path: &str, offset: i64) -> vec::Vec<Scrobble> {
    let file = fs::read_to_string(path).expect("Error Opening File");

    if !file.contains(AUDIOSCROBBLER_HEADER) {
        panic!("Not a valid .scrobbler.log!");
    } else {
        let mut scrobbles = vec::Vec::new();
        let file_as_lines: Vec<&str> = file.lines().collect();

        for index in TRACKS_BEGIN_INDEX..file_as_lines.len() {
            if file_as_lines[index].contains(LISTENED) {
                let track_str: Vec<&str> = file_as_lines[index].split(SEPARATOR).collect();

                let mut scrobble = Scrobble::new(
                    track_str[ARTIST_INDEX],
                    track_str[TITLE_INDEX],
                    track_str[ALBUM_INDEX],
                );

                if offset != 0 {
                    scrobble.with_timestamp(convert_time(offset, track_str[TIMESTAMP_INDEX]));
                }

                &scrobbles.push(scrobble);
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
