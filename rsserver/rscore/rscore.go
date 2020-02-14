package rscore

/*

    Package of core operations and data types
    Used by randomslide

*/

import (
    "log"
)

const DEFAULTPORT = 6291
const DBNAME = "./data/random.db"
const LOGPATH = "./static/log/"
const PIDFILEPATH = "./data/"

const VOLATILEMODE = true

const L_REQ = 0
const L_RESP = 1
const L_SHUTDOWN = 2

const C_OK = 0

var DBUC = []byte("dbuc")       // deck bucket
var TBUC = []byte("tbuc")       // text bucket
var IBUC = []byte("ibuc")       // image bucket
var SBUC = []byte("sbuc")       // settings bucket

var SETTINGSKEY = []byte("skey")

type Settings struct {
    Verb bool
    Cid int
    Pidfile string
}

type Deckreq struct {
    N int
    Lang string
    Cat string
}

type Textreq struct {
    Text string                 // The text object to add to db
    Tags string                 // whitespace separated tags for indexing
}

type Deck struct {
    N int                       // Total number of slides in deck
    Lang string                 // Deck language, 'en', 'de', 'se', etc
    Slides []Slide              // Slice of Slide objects
}

type Slide struct {
    Title string                // Slide title
    Imgur string                // URL to image
    Btext string                // Body text
}

type Statusresp struct {
    Code int                    // Error code to be parsed in frontend
    Text string                 // Additional related data
}

// Log all errors to file
func Cherr(e error) {
    if e != nil { log.Fatal(e) }
}

