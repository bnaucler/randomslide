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
const PIDFILEPATH = "./data/"

const DEBUG = 1                 // set to 1 for verbose output TODO: replace with flag

var dbuc = []byte("dbuc")       // deck bucket
var ibuc = []byte("ibuc")       // image bucket
var sbuc = []byte("sbuc")       // settings bucket

type Resp struct {
    Data string
    Id int
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

// Generate a deck based on request
func deckreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB, cid int) int {

    e := r.ParseForm()
    cherr(e)

    deck := Deck{
            N: cid,
            Lang: "en",
            Slides: make([]Slide, cid) }

    key := []byte(strconv.Itoa(cid)) // TODO: make this make sense somehow

    mdeck, e := json.Marshal(deck)

    e = wrdb(db, key, mdeck, dbuc)
    cherr(e)

    v, e := rdb(db, key, dbuc)
    cherr(e)

    rdeck := Deck{}
    e = json.Unmarshal(v, &rdeck)

    enc := json.NewEncoder(w)
    enc.Encode(rdeck)

    cid++
    return cid
}

func main() {

    rand.Seed(time.Now().UnixNano())

	pptr := flag.Int("p", DEFAULTPORT, "port number to listen")
	dbptr := flag.String("d", DBNAME, "specify database to open")
	flag.Parse()

    db, e := bolt.Open(*dbptr, 0640, nil)
    cherr(e)
    defer db.Close()

    pid := os.Getpid()
    prgname := filepath.Base(os.Args[0])
    pidfile := fmt.Sprintf("%s/%s.pid", PIDFILEPATH, prgname)
    e = ioutil.WriteFile(pidfile, []byte(strconv.Itoa(pid)), 0644)

    if DEBUG != 0 {
        fmt.Printf("DEBUG: %s started with PID: %d\n", prgname, pid)
    }

    cid := 0

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Slide requests
    http.HandleFunc("/getdeck", func(w http.ResponseWriter, r *http.Request) {
        cid = deckreqhandler(w, r, db, cid)
        fmt.Printf("DEBUG: %+v\n", cid)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    cherr(e)
}
