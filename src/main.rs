mod api_keys;
mod auth;
mod init;
mod log;
extern crate rustfm_scrobble;

fn main() {
    let app = init::app_info();

    let arguments = app.get_matches();
    let file_path = arguments
        .value_of("file")
        .expect("Failed to parse file path!");
    /* Get argument value and unwrap to type str, then parse the string and unwrap str for
     * conversion to f32*/
    let timezone_offset: f32 = arguments.value_of("offset").unwrap().parse().unwrap();

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

    /* Ask user if they want to delete file */
}
