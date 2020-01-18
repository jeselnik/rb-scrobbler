use std::fs;
extern crate app_dirs;
use app_dirs::*;
extern crate rustfm_scrobble;
use rustfm_scrobble::Scrobbler;

const APP_INFO: AppInfo = app_dirs::AppInfo {
    name: "rb-scrobbler",
    author: "Eddie Jeselnik",
};

pub fn initial_authentication(scrob: &mut Scrobbler, user: &str, pass: &str) {
    let auth_res = scrob.authenticate_with_password(user, pass);
    if auth_res.is_ok() {
        println!("Authenticated!");
        let session_key = scrob.session_key()
            .expect("Couldn't get session key!");
        fs::write("placeholder-path", session_key)
            .expect("Couldn't save session key!'");
    } else {
        panic!("Authentication failed!");
    }
}
