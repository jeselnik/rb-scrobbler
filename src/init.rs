extern crate clap;
use clap::{App, Arg, SubCommand};

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
        .subcommand(
            SubCommand::with_name("auth")
                .about("Authenticate with last.fm (only required once)")
                .arg(
                    Arg::with_name("user")
                        .short("u")
                        .long("username")
                        .value_name("USER")
                        .help("Your last.fm username")
                        .required(true)
                        .takes_value(true),
                )
                .arg(
                    Arg::with_name("pass")
                        .short("p")
                        .long("password")
                        .value_name("PASS")
                        .help("Your last.fm password")
                        .required(true)
                        .takes_value(true),
                ),
        );

    return app;
}
