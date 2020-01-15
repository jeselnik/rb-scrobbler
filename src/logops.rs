use std::fs;
use std::string;
use std::vec;

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

pub fn log_to_vec(path: &str) -> vec::Vec<Track> {
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
