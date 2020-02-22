# Contributing to randomslide
Contributions are very welcome! At this early stage, this file will double as a TODO list.

## Code style
* Keep variable & function names short and in lower case (whenever possible)
* JSON format used for data exchange
* Limit to 80 columns
* Limit to three levels of indentation
* Indent by four spaces

## Defaults
* Backend requests are sent as POST
* Server listens at port 6291 by default

## Folder structure
```
bin/            - Script and compiled binaries
data/           - Home of database & pidfiles
lib/            - Custom golang libraries
rsserver/       - Source file(s) of server
static/         - HTML, CSS, JS and images
static/log      - Server and monitor log files
tests/          - Scripts for server tests
tools/          - Tools for maintenence etc.
```

## TODO
* Generate slides w. numbers
* Different slide types
* Handling of image types/sizes?
* Refactoring image handler

## IRC
The main developers can be found in the channel #ljusdal at EFNet.

## Pull requests
Are appreciated!
