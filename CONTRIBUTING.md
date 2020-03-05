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
static/log/     - Server and monitor log files
tools/          - Tools for maintenence etc.
```

## TODO
* Refactoring image handler
* User account access levels
* Security checks with session keys
* Proper automated tests
* Check image repo size requirements
* Verbosity levels
* Image upload status reporting

## Text types
Title (ttext): 1-35 characters  
Body (btext): 5-80 characters

## Slide types
```
Index   Img sz    Type                        Description
0       XL/16:9   Big title                   Could be good to start a slide set with
1       XL/16:9   Full screen picture
2       NULL      Big number                  A slide just saying things like '+12%'
3       9:16/NULL Bullet point list
4       16:9      Title, pic, body text       What we already have in alpha
5       NULL      'Inspirational quote'       Soo much potential here
6       16:9/9:16 Picture with text           Theme can decide if text goes under, next to etc
7       NULL      Graph                       No good slideshow is complete without a graph
```

# Image sizes
Images will either be classified as XL (fullscreen), 16:9 (width:height) or 9:16 (width:height).
```

## IRC
The main developers can be found in the channel #ljusdal at EFNet.

## Pull requests
Are appreciated!
