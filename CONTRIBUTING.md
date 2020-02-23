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
* Store generated decks in db for direct access
* Store image sizes in object
* Add image size to requests

## Slide types
* Big title
    - Could be good to start a slide set with

* Full screen picture

* Bullet point list
    - Needs new database object..

* Title, picture, body text
    - What we already have in alpha

* 'inspirational quote'
    - Soo much potential here

* Picture with adjacent text
    - Theme can decide if text goes under, next to etc

# Image sizes
Images used will be classified in the following sizes:
```
Size (px)   Min W       Min H       Max W       Max H
0 (S):
1 (M):
2 (L):
3 (XL)
```

## IRC
The main developers can be found in the channel #ljusdal at EFNet.

## Pull requests
Are appreciated!
