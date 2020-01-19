mod auth;
mod init;
mod log;
extern crate rustfm_scrobble;

fn main() {
    const API_KEY: &str = "";
    const API_SECRET: &str = "";

    let app = init::app_info();

    let arguments = app.get_matches();
    let file_path = arguments.value_of("file")
        .expect("Failed to parse file path!");
    /* Get argument value and unwrap to type str, then parse the string and unwrap str for
     * conversion to f32*/
    let timezone_offset: f32 = arguments.value_of("offset").unwrap().parse().unwrap();

    let mut scrobbler = rustfm_scrobble::Scrobbler::new(API_KEY, API_SECRET);

   if arguments.is_present("auth") {
       auth::initial_authentication(&mut scrobbler, username, password);
    } 

    let scrobbles = log::as_scrobbles(file_path, (timezone_offset * 60.0) as i64);

    for index in 0..scrobbles.len() {
        let scrobble_response = scrobbler.scrobble(&scrobbles[index]);
        println!("{:?}", scrobble_response);
    }

    /* Ask user if they want to delete file */
}
