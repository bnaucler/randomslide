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

var DBUC = []byte("dbuc")       // Deck bucket
var TBUC = []byte("tbuc")       // Title text bucket
var BBUC = []byte("bbuc")       // Body text bucket
var IBUC = []byte("ibuc")       // Image bucket
var SBUC = []byte("sbuc")       // Settings bucket

var SETTINGSKEY = []byte("skey")

type Settings struct {
    Verb bool                   // Verbosity level
    Cid int                     // Max id TODO: remove
    Pidfile string              // Location of pidfile
    Taglist []string            // List of all existing tags
    Btmax map[string]int        // Maximal value for randomization
    Ttmax map[string]int        // Maximal value for randomization
}

type Deckreq struct {
    N int                       // Number of slides to generate
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Tags string                 // Tags on which to base the deck
}

type Textreq struct {
    Ttext string                // Title text object to add to db
    Btext string                // Body text object to add to db
    Tags string                 // whitespace separated tags for indexing
}

type Tagresp struct {
    Tags []string               // Array of tags for indexing
}

type Textobj struct {
    Id int                      // # for random selection
    Text string                 // The text itself
    Tags []string               // All tags where object exists (for associative decks)
}

type Deck struct {
    N int                       // Total number of slides in deck
    Lang string                 // Languge code, 'en', 'de', 'se', etc
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

