package main

import (
    "os"
    "fmt"
    "net"
    "log"
    "time"
    "flag"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/boltdb/bolt"
)

const DEFAULTPORT = 6291
const DBNAME = "./data/random.db"
const LOGPATH = "./static/log/"
const PIDFILEPATH = "./data/"

const VOLATILEMODE = true

const L_REQ = 0
const L_RESP = 1

var dbuc = []byte("dbuc")       // deck bucket
var tbuc = []byte("tbuc")       // text bucket
var ibuc = []byte("ibuc")       // image bucket
var sbuc = []byte("sbuc")       // settings bucket

var SETTINGSKEY = []byte("skey")

type Settings struct {
    Verb bool
    Cid int
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
    Tcolor string               // Text color in CSS-compatible hex code
    Bgcolor string              // Body color in CSS-compatible hex code
}

// TODO: Create a status response object - implement in addtext

// Log all errors to console
func cherr(e error) {
    if e != nil { log.Fatal(e) }
}

// Retrieves client IP address from http request
func getclientip(r *http.Request) string {

    ip, _, e := net.SplitHostPort(r.RemoteAddr)
    cherr(e)

    return ip
}

// Log file wrapper
// TODO: Use interface() instead of []byte and ltype
//       log levels
func addlog(ltype int, msg []byte, r *http.Request) {

    ip := getclientip(r)
    var lentry string

    switch ltype {
        case L_REQ:
            lentry = fmt.Sprintf("REQ from %s: %s", ip, msg)

        case L_RESP:
            lentry = fmt.Sprintf("RESP to %s: %s", ip, msg)

        default:
            lentry = fmt.Sprintf("Something happened, but I don't know how to log it!")
    }

    log.Println(lentry)
}

// Initialize logger
func initlog(prgname string) {

    logfile := fmt.Sprintf("%s/%s.log", LOGPATH, prgname)

    f, e := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    cherr(e)

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Write JSON encoded byte slice to DB
func wrdb(db *bolt.DB, k []byte, v []byte, cbuc []byte) (e error) {

    e = db.Update(func(tx *bolt.Tx) error {
        b, e := tx.CreateBucketIfNotExists(cbuc)
        if e != nil { return e }

        e = b.Put(k, v)
        if e != nil { return e }

        return nil
    })
    return
}

// Return JSON encoded byte slice from DB
func rdb(db *bolt.DB, k []byte, cbuc []byte) (v []byte, e error) {

    e = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(cbuc)
        if b == nil { return fmt.Errorf("No bucket!") }

        v = b.Get(k)
        return nil
    })
    return
}

// Create random string of length ln
func randstr(ln int) string {

    const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
    var cslen = len(charset)

    b := make([]byte, ln)
    for i := range b { b[i] = charset[rand.Intn(cslen)] }

    return string(b)
}

// Wrapper for writing settings to database
func wrsettings(db *bolt.DB, settings Settings) {

    mset, e := json.Marshal(settings)
    cherr(e)

    e = wrdb(db, SETTINGSKEY, mset, sbuc)
    cherr(e)
}

// Wrapper for reading settings from database
func rsettings(db *bolt.DB) Settings {

    settings := Settings{}

    mset, e := rdb(db, SETTINGSKEY, sbuc)
    if e != nil { return Settings{} }

    e = json.Unmarshal(mset, &settings)
    cherr(e)

    return settings
}

// Returns a full slide deck according to request
func mkdeck(req Deckreq, db *bolt.DB, settings Settings) (Deck, Settings) {

    deck := Deck{
            N: req.N,
            Lang: req.Lang }

    deck.Slides = make([]Slide, req.N)

    txtreq := Textreq{}

    for i := 0; i < req.N; i++ {
        mtxtobj, e := rdb(db, []byte(strconv.Itoa(rand.Intn(settings.Cid))), tbuc)
        json.Unmarshal(mtxtobj, &txtreq)
        deck.Slides[i].Title = txtreq.Text
        cherr(e)
    }

    return deck, settings
}

// Handles incoming requests for decks
func deckreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings Settings) Settings {

    e := r.ParseForm()
    cherr(e)

    n, e := strconv.Atoi(r.FormValue("amount"))
    cherr(e)

    req := Deckreq{
            N: n,
            Lang: r.FormValue("lang"),
            Cat: r.FormValue("category") }

    mreq, e := json.Marshal(req)
    cherr(e)
    addlog(L_REQ, mreq, r)

    deck, settings := mkdeck(req, db, settings)

    mdeck, e := json.Marshal(deck)
    addlog(L_RESP, mdeck, r)

    enc := json.NewEncoder(w)
    enc.Encode(deck)

    return settings
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings Settings) Settings {

    e := r.ParseForm()
    cherr(e)

    t := Textreq{
            Text: r.FormValue("text"),
            Tags: r.FormValue("tags") }

    // TODO: Create tag searchability
    // sptags := strings.Split(t.Tags, ",")

    key := []byte(strconv.Itoa(settings.Cid)) // TODO: make this make sense somehow
    mtxt, e := json.Marshal(t)

    e = wrdb(db, key, mtxt, tbuc)
    cherr(e)

    addlog(L_REQ, mtxt, r)

    settings.Cid++
    wrsettings(db, settings)
    return settings
}

func main() {

    rand.Seed(time.Now().UnixNano())

    pptr := flag.Int("p", DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", DBNAME, "specify database to open")
    vptr := flag.Bool("v", false, "verbose mode")
    flag.Parse()

    db, e := bolt.Open(*dbptr, 0640, nil)
    cherr(e)
    defer db.Close()

    settings := rsettings(db)
    settings.Verb = *vptr

    pid := os.Getpid()
    prgname := filepath.Base(os.Args[0])
    pidfile := fmt.Sprintf("%s/%s.pid", PIDFILEPATH, prgname)
    e = ioutil.WriteFile(pidfile, []byte(strconv.Itoa(pid)), 0644)

    initlog(prgname)

    if settings.Verb == true {
        log.Printf("DEBUG: %s started with PID: %d\n", prgname, pid)
    }

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    if VOLATILEMODE == true {
        http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
            log.Printf("Restart request received. Shutting down.\n")
            os.Exit(1)
        })
    }

    // Slide requests
    http.HandleFunc("/getdeck", func(w http.ResponseWriter, r *http.Request) {
        settings = deckreqhandler(w, r, db, settings)
    })

    // Add text to db
    http.HandleFunc("/addtext", func(w http.ResponseWriter, r *http.Request) {
        settings = textreqhandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    cherr(e)
}
