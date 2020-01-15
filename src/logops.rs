use std::fs;
use std::result;
use std::vec;

const AUDIOSCROBBLER_HEADER: &str = "#AUDIOSCROBBLER/";
const TRACKS_BEGIN_INDEX: usize = 3;
const LISTENED: &str = "\tL\t";

struct Track {
    artist: String,
    album: String,
    title: String,
    timestamp: u64,
}

pub fn log_to_vec(path: &str) -> Result<&str, &str> {
    let file = fs::read_to_string(path)
    .expect("Error Opening File");

    if !file.contains(AUDIOSCROBBLER_HEADER) {
        return result::Result::Err("Not a valid .scrobbler.log!");
    } else {
        //let mut tracks = vec::Vec::new();
        let file_as_lines: Vec<&str> = file.lines().collect();
        for index in TRACKS_BEGIN_INDEX..file_as_lines.len() {
            println!("{}", file_as_lines[index]);
        }


        return result::Result::Ok("Yeah it's working'");
    }
}
