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

## Themes
Deck themes are defined by css files located in `static/css/themes/`.  
Look at the files `empty.css` and `white.css` to get an idea of what does what.

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

## TODO (Backend)
* Refactor getimgtype()
* Handler for text file batch processing
* Enable goroutine for tag index mutex
* Handler for SMTP setup
* Proper automated tests
* Improve API call logging
* Web scraper - generate data from url
* Add logging to rmhandler
* Improved bullet point handling
* Put untagged data in 'random' tag
* Rebuild setslidetype()
* Anonymize skeys when logging


## TODO (Frontend)
* Share (deckid) on social media
* CSS for all slides (maybe not 1 and 5)
* CSS for the main site
* JS-removing unused tags from DOM after slideshow has started


## Text types
Title (ttext): 1-35 characters  
Body (btext): 5-80 characters

## Slide types
```
Index   Img sz    Type                        Description
0       0/1       Big title                   Could be good to start a slide set with
1       0         Full screen image
2       NULL      Big number                  A slide just saying things like '+12%'
3       3         Bullet point list
4       1         Title, image & subtitle
5       NULL      'Inspirational quote'       Soo much potential here
6       1/2/3     Picture with text           Theme can decide if text goes under, next to etc
7       NULL      Graph                       No good slideshow is complete without a graph
```
See `Image sizes` section below for image reference sizes.

# Image sizes
```
Index   Type            Aspect ratio        Min size        Max size
0       XL              20:10-13:10         1920x1080       1920x1080
1       Landscape       20:10-13:10         640x360         1920x1080
2       Box             12:10-9:10          360x360         1080x1080
3       Portrait        8:10-5:10           360x640         1080x1920
```
Images larger than the max size will be automatically scaled to this size. Please note that this might have undesirable effects for picture quality, but has been deemed necessary to improve loading times.

## IRC
The main developers can be found in the channel #ljusdal at EFNet.

## Pull requests
Are appreciated!
