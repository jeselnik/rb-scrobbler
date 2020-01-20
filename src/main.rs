mod api_keys;
mod auth;
mod init;
mod log;
extern crate rustfm_scrobble;
use std::{fs, io, string};

fn main() {
    let app = init::app_info();

    let arguments = app.get_matches();
    let file_path = arguments
        .value_of("file")
        .expect("Failed to parse file path!");
    /* Get argument value and unwrap to type str, then parse the string and unwrap str for
     * conversion to f32*/
    let timezone_offset: f32 = arguments
        .value_of("offset")
        .unwrap()
        .parse()
        .expect("Offset not a number!");

    let mut scrobbler = rustfm_scrobble::Scrobbler::new(api_keys::API_KEY, api_keys::API_SECRET);

    if arguments.is_present("auth") {
        let auth_args = arguments
            .subcommand_matches("auth")
            .expect("Couldn't unwrap subcommand matches!");
        let username = auth_args
            .value_of("user")
            .expect("Couldn't unwrap username!");
        let password = auth_args
            .value_of("pass")
            .expect("Couldn't unwrap password!");
        auth::initial_authentication(&mut scrobbler, username, password);
    }

    auth::authenticate_key(&mut scrobbler);

    let scrobbles = log::as_scrobbles(file_path, (timezone_offset * 60.0) as i64);

    for individual_scrobble in scrobbles {
        let scrob_response = scrobbler.scrobble(&individual_scrobble);
        if scrob_response.is_ok() {
            println!(
                "[OK] {} - {}",
                individual_scrobble.artist(),
                individual_scrobble.track()
            );
        } else {
            println!(
                "[FAIL] {} - {}",
                individual_scrobble.artist(),
                individual_scrobble.track()
            );
        }
    }

    println!("Delete file \"{}\"?", file_path);

    let mut user_choice = string::String::new();
    io::stdin()
        .read_line(&mut user_choice)
        .expect("Failed to read from stdout!");

    if user_choice.to_lowercase().starts_with("y") {
        fs::remove_file(file_path).expect("I/O Error!");
        println!("\"{}\" deleted!", file_path);
    }
}
