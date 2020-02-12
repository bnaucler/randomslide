package main

import (
    "os"
    "fmt"
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

var dbuc = []byte("dbuc")       // deck bucket
var ibuc = []byte("ibuc")       // image bucket
var sbuc = []byte("sbuc")       // settings bucket

type Settings struct {
    Verb bool
    Cid int
}

type Deckreq struct {
    N int
    Lang string
    Cat string
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

// Log all errors to console
func cherr(e error) {
    if e != nil { log.Fatal(e) }
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

// Returns a slide according to request
func mkslide(req Deckreq) Slide {

    slide := Slide{Title: randstr(10) }

    return slide
}

// Returns a full slide deck according to request
func mkdeck(req Deckreq) Deck {

    deck := Deck{
            N: req.N,
            Lang: req.Lang }

    for i := 0; i < req.N; i++ {
        deck.Slides = append(deck.Slides, mkslide(req))
    }

    return deck
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

    deck := mkdeck(req)

    key := []byte(strconv.Itoa(settings.Cid)) // TODO: make this make sense somehow

    mdeck, e := json.Marshal(deck)

    e = wrdb(db, key, mdeck, dbuc)
    cherr(e)

    v, e := rdb(db, key, dbuc)
    cherr(e)

    rdeck := Deck{}
    e = json.Unmarshal(v, &rdeck)

    enc := json.NewEncoder(w)
    enc.Encode(rdeck)

    settings.Cid++
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

    settings := Settings{}
    settings.Verb = *vptr
    settings.Cid = 0

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

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    cherr(e)
}
