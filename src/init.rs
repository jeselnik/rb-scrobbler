extern crate clap;
use clap::{App, Arg};

pub fn app_info() -> clap::App<'static, 'static> {
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
                .help(
                    "Your last.fm username (only required for 
            initial authentication)",
                )
                .required(false)
                .takes_value(true)
                .default_value(""),
        )
        .arg(
            Arg::with_name("pass")
                .short("p")
                .long("password")
                .value_name(
                    "PASSWORD (only required for 
            initial authentication)",
                )
                .required(false)
                .takes_value(true)
                .default_value(""),
        );

    return app;
}
