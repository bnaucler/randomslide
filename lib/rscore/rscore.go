package rscore

/*

    Package of core operations and data types
    Used by randomslide

*/

import (
    "os"
    "fmt"
    "log"
    "net"
    "time"
    "regexp"
    "strconv"
    "net/http"
    "io/ioutil"
    "os/signal"
    "encoding/json"
    "path/filepath"
)

const DEFAULTPORT = 6291        // Default port can also be supplied with -p flag
const DBNAME = "./data/rs.db"   // Database location (change at own peril)
const LOGPATH = "./static/log/" // Logs should be accessible from frontend
const IMGDIR = "./static/img/"  // Image directory
const PIDFILEPATH = "./data/"   // Base directory for storage of PID file

const VOLATILEMODE = true

const TTEXTMAX = 35             // Max length for title text
const BTEXTMAX = 80             // Max length for body text

const STYPES = 7                // Number of slide types available

// Logging codes parsed by Addlog()
const L_REQ = 0                 // Request log
const L_RESP = 1                // Response log
const L_SHUTDOWN = 2            // Server shutdown request log

// Status response codes sent to client
const C_OK = 0                  // OK
const C_WRFF = 1                // Incorrect file format
const C_WRSZ = 2                // Not able to classify image size

// Min bounds for image sizes (w, h)
var IMGMIN = [][]int{
    {150, 150},                 // 0: Small
    {500, 500},                 // 1: Medium
    {1000, 1000},               // 2: Large
    {1920, 1080},               // 3: X Large
}

// Max bounds for image sizes (w, h)
var IMGMAX = [][]int{
    {499, 499},                 // 0: Small
    {999, 799},                 // 1: Medium
    {1919, 1079},               // 2: Large
    {3000, 3000},               // 3: X Large
}

var DBUC = []byte("dbuc")       // Deck bucket
var TBUC = []byte("tbuc")       // Title text bucket
var BBUC = []byte("bbuc")       // Body text bucket
var IBUC = []byte("ibuc")       // Image bucket
var SBUC = []byte("sbuc")       // Settings bucket

var INDEX = []byte(".index")    // Untouchable database index position

// For verification of image mime types
var IMGMIME = []string{
    "image/jpeg",
    "image/png",
    "image/gif" }

type Settings struct {
    Verb bool                   // Verbosity level
    Dmax int                    // Max id of decks
    Tmax int                    // Max id of title objects
    Bmax int                    // Max id of body objects
    Imax int                    // Max id of image objects
    Pidfile string              // Location of pidfile
    Taglist []string            // List of all existing tags TODO: Make map w ID
}

type Deckreq struct {
    Id int                      // Deck ID for db
    N int                       // Number of slides to generate
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Tags []string               // Slice of tags on which to base search
    Isidreq bool                // true if request has specified ID
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
    Size int                    // S (0), M (1), L (2) or XL (3)
    H int                       // Image height
    W int                       // Image width
}

type Deck struct {
    Id int                      // Deck ID for db
    N int                       // Total number of slides in deck
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Slides []Slide              // Slice of Slide objects
}

type Slidetype struct {
    Type int                    // The type ID number
    TT bool                     // Includes title text
    BT bool                     // Includes body text
    IMG bool                    // Includes image
}

type Slide struct {
    Type int                    // See CONTRIBUTING.md for type chart
    Title string                // Slide title
    Img Imgobj                  // Image object
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

// Enables clean shutdown. Needs delay for caller to send response
func Shutdown(settings Settings) {

    go func() {
        time.Sleep(1 * time.Second)
        os.Remove(settings.Pidfile)
        os.Exit(0)
    }()
}

// Setting up signal handler
func Sighandler(settings Settings) {

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func(){
        for sig := range sigc {
            fmt.Printf("Caught %+v - cleaning up.\n", sig)
            Shutdown(settings)
        }
    }()
}

// Creating PID file
func Mkpidfile(settings Settings, prgname string, pid int) Settings {

    settings.Pidfile = fmt.Sprintf("%s%s.pid", PIDFILEPATH, prgname)
    e := ioutil.WriteFile(settings.Pidfile, []byte(strconv.Itoa(pid)), 0644)
    Cherr(e)

    return settings
}

// Initialize logger
func Initlog(prgname string) {

    logfile := fmt.Sprintf("%s/%s.log", LOGPATH, prgname)

    f, e := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    Cherr(e)

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Creates pid file, signal handler and starts logging
func Rsinit(settings Settings) Settings {

    prgname := filepath.Base(os.Args[0])
    pid := os.Getpid()

    settings = Mkpidfile(settings, prgname, pid)
    Sighandler(settings)
    Initlog(prgname)

    log.Printf("%s started with PID: %d\n", prgname, pid)

    return settings
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

// Removes all files residing in dir (except .gitkeep)
func Rmall(dir string) {

    d, e := os.Open(dir)
    Cherr(e)
    defer d.Close()

    fl, e := d.Readdirnames(-1)
    Cherr(e)

    for _, fn := range fl {
        if fn == ".gitkeep" { continue }
        e = os.RemoveAll(filepath.Join(dir, fn))
        Cherr(e)
    }
}

// Sends a status code response as JSON object
func Sendstatus(code int, text string, w http.ResponseWriter) {

    resp := Statusresp{
            Code: code,
            Text: text }

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}
