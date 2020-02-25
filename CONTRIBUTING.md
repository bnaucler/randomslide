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
* Refactoring image handler
* User accounts with access levels
* Proper automated tests
* CLI tool for automated db import
* Verification of text lengths
* Check image repo size requirements
* Verbosity levels

## Text types
Title (ttext): 1-35 characters  
Body (btext): 5-80 characters

## Slide types
```
Index   Type                            Description
0       Big title                       Could be good to start a slide set with
1       Full screen picture
2       Big number                      A slide just saying things like '+12%'
3       Bullet point list
4       Title, pic, body text           What we already have in alpha
5       'Inspirational quote'           Soo much potential here
6       Picture with text               Theme can decide if text goes under, next to etc
```

# Image sizes
Images used will be classified in the following sizes:
```
Size (px)   Min W       Min H       Max W       Max H
0 (S):      150         150         499         499
1 (M):      500         500         999         799
2 (L):      1000        800         1919        1079
3 (XL)      1920        1080        3000+       3000+
```

## IRC
The main developers can be found in the channel #ljusdal at EFNet.

## Pull requests
Are appreciated!
