use std::fs;
extern crate rustfm_scrobble;
use rustfm_scrobble::Scrobbler;

pub fn initial_authentication(scrob: &mut Scrobbler, user: &str, pass: &str) {
    let auth_res = scrob.authenticate_with_password(user, pass);
    if auth_res.is_ok() {
        println!("Authenticated!");

        let session_key = scrob.session_key()
            .expect("Couldn't get session key!");

    } else {
        panic!("Authentication failed!");
    }
}
