use std::io;
mod log;
extern crate app_dirs;
use app_dirs::*;
extern crate clap;
use clap::{App, Arg};
extern crate rustfm_scrobble;

fn main() {
    let app = App::new("Minimal Rockbox Scrobbler")
        .version("0.1")
        .author("Eddie Jeselnik <eddie@jeselnik.xyz>")
        .about("Submit .scrobbler.log files to last.fm")
        .arg(
            Arg::with_name("file")
                .short("f")
                .long("file")
                .value_name("FILE")
                .help("File path to your .scrobbler.log")
                .required(true)
                .takes_value(true),
        )
        .arg(
            Arg::with_name("offset")
                .short("o")
                .long("offset")
                .value_name("OFFSET")
                .help("Offset (in hours) of local time to UTC")
                .required(false)
                .takes_value(true)
                .default_value("0"),
        )
        .arg(
            Arg::with_name("user")
            .short("u")
            .long("username")
            .value_name("USERNAME")
            .help("Your last.fm username (only required for 
            initial authentication)")
            .required(false)
            .takes_value(true)
            .default_value("")
            )
        .arg(
            Arg::with_name("pass")
            .short("p")
            .long("password")
            .value_name("PASSWORD (only required for 
            initial authentication)")
            .required(false)
            .takes_value(true)
            .default_value(""));

    const APP_INFO: app_dirs::AppInfo = app_dirs::AppInfo{name:"rb-scrobbler", author:"Eddie Jeselnik"};
    let arguments = app.get_matches();

    let username = arguments.value_of("user").unwrap();
    let password = arguments.value_of("pass").unwrap();

    const API_KEY: &str = "INSERT_API_KEY";
    const API_SECRET: &str = "INSERT_API_SECRET";

    if username != "" && password != "" {
        let mut scrobbler = rustfm_scrobble::Scrobbler::new(API_KEY, API_SECRET);
        let response = scrobbler.authenticate_with_password(username, password);
        println!("{:?}", response);
    }

    let file_path = arguments.value_of("file").unwrap();
    /* Get argument value and unwrap to type str, then parse the string and unwrap str for
     * conversion to f32*/
    let timezone_offset: f32 = arguments.value_of("offset").unwrap().parse().unwrap();

    let mut tracks = log::as_vec(file_path);

    println!("{}", tracks[0].timestamp);

    if timezone_offset != 0.0 {
        for index in 0..tracks.len() {
            tracks[index].convert_time((timezone_offset * 60.0) as i64);
        }
    }

    println!("{}", tracks[0].timestamp);

    /* Scrobble */

    /* Ask user if they want to delete file */
}
