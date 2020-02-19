package rscore

/*

    Package of core operations and data types
    Used by randomslide

*/

import (
    "fmt"
    "log"
    "net"
    "regexp"
    "net/http"
)

const DEFAULTPORT = 6291
const DBNAME = "./data/rs.db"
const LOGPATH = "./static/log/"
const PIDFILEPATH = "./data/"

const VOLATILEMODE = true

const L_REQ = 0
const L_RESP = 1
const L_SHUTDOWN = 2

const C_OK = 0                  // Status OK
const C_WRFF = 1                // Error code 1 - incorrect file format

var DBUC = []byte("dbuc")       // Deck bucket
var TBUC = []byte("tbuc")       // Title text bucket
var BBUC = []byte("bbuc")       // Body text bucket
var IBUC = []byte("ibuc")       // Image bucket
var SBUC = []byte("sbuc")       // Settings bucket

var INDEX = []byte(".index")

var IMGMIME = []string{
    "image/jpeg",
    "image/png",
    "image/gif" }

type Settings struct {
    Verb bool                   // Verbosity level
    Tmax int                    // Max id of title objects
    Bmax int                    // Max id of body objects
    Imax int                    // Max id of image objects
    Pidfile string              // Location of pidfile
    Taglist []string            // List of all existing tags TODO: Make map w ID
}

type Deckreq struct {
    N int                       // Number of slides to generate
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Tags []string               // Slice of tags on which to base search
}

type Textreq struct {
    Ttext string                // Title text object to add to db
    Btext string                // Body text object to add to db
    Tags []string               // Tags for indexing
}

type Tag struct {
    Ids    []int                // All IDs associated with tag
}

type Rtag struct {
    Name string                  // Tag name
    TN int                       // Number of title objects in db
    BN int                       // Number of body objects in db
    IN int                       // Number of image objects in db
}

type Tagresp struct {
    Tags []Rtag                 // Array of tags for indexing
}

type Textobj struct {
    Id int                      // Index number
    Text string                 // The text itself
    Tags []string               // All tags where object exists (for associative decks)
}

type Imgobj struct {
    Id int                      // Index number
    Fname string                // File name
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
func Cherr(e error) error {
    if e != nil { log.Fatal(e) }
    return e
}

// Removes whitespace and special characters from string
func Cleanstring(src string) string {

    rx, e := regexp.Compile("[^a-zåäöüæø]+")
    Cherr(e)

    dst := rx.ReplaceAllString(src, "")

    return dst
}

// Returns true if string is present in list
func Findstrinslice(v string, list []string) bool {

    for _, t := range list {
        if v == t { return true }
    }

    return false
}

// Retrieves client IP address from http request
func getclientip(r *http.Request) string {

    ip, _, e := net.SplitHostPort(r.RemoteAddr)
    Cherr(e)

    return ip
}

// Log file wrapper
func Addlog(ltype int, msg []byte, r *http.Request) {

    ip := getclientip(r)
    var lentry string

    switch ltype {
        case L_REQ:
            lentry = fmt.Sprintf("REQ from %s: %s", ip, msg)

        case L_RESP:
            lentry = fmt.Sprintf("RESP to %s: %s", ip, msg)

        case L_SHUTDOWN:
            lentry = fmt.Sprintf("Server shutdown requested from %s", ip)

        default:
            lentry = fmt.Sprintf("Something happened, but I don't know how to log it!")
    }

    log.Println(lentry)
}

// Convert file size to human readable format
func Prettyfsize(b int64) string {

    k := b / 1024
    m := k / 1024
    var ret string

    if m == 0 {
        ret = fmt.Sprintf("%d.%dKB", k, (b % 1024) / 100)
    } else {
        ret = fmt.Sprintf("%d.%dMB", m, k / 100)
    }

    return ret
}
