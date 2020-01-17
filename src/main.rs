mod init;
mod log;
extern crate app_dirs;
use app_dirs::*;
extern crate rustfm_scrobble;

fn main() {
    const APP_INFO: app_dirs::AppInfo = app_dirs::AppInfo{name:"rb-scrobbler", author:"Eddie Jeselnik"};
    const API_KEY: &str = "INSERT_API_KEY";
    const API_SECRET: &str = "INSERT_API_SECRET";

    let app = init::app_info();

    let arguments = app.get_matches();
    let username = arguments.value_of("user").unwrap();
    let password = arguments.value_of("pass").unwrap();
    let file_path = arguments.value_of("file").unwrap();
    /* Get argument value and unwrap to type str, then parse the string and unwrap str for
     * conversion to f32*/
    let timezone_offset: f32 = arguments.value_of("offset").unwrap().parse().unwrap();

    if username != "" && password != "" {
        let mut scrobbler = rustfm_scrobble::Scrobbler::new(API_KEY, API_SECRET);
        let response = scrobbler.authenticate_with_password(username, password);
        println!("{:?}", response);
    }

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
