mod log;
extern crate clap;
use clap::{App, Arg};

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
        );

    let arguments = app.get_matches();

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
