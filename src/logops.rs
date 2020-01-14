use std::fs;

const AUDIOSCROBBLER_HEADER: &str = "#AUDIOSCROBBLER/";

pub fn log_to_vec(path: String) {
    let file = fs::read_to_string(path)
    .expect("Error Opening File");

    let file_as_lines: Vec<&str> = file.lines().collect();
}
