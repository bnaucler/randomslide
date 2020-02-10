package main

import (
    "fmt"
    "log"
    "time"
    "flag"
    "strconv"
    "net/http"
    "math/rand"

    "github.com/boltdb/bolt"
)

const DEFAULTPORT = 6291
const DBNAME = "./db/random.db"

var tbuc = []byte("pbuc")       // text bucket
var ibuc = []byte("ibuc")       // image bucket
var sbuc = []byte("sbuc")       // settings bucket

// Log all errors to console
func cherr(e error) {
    if e != nil { log.Fatal(e) }
}

// Write JSON encoded byte slice to DB
func wrdb(db *bolt.DB, k int, v []byte, cbuc []byte) (e error) {

    e = db.Update(func(tx *bolt.Tx) error {
        b, e := tx.CreateBucketIfNotExists(cbuc)
        if e != nil { return e }

        e = b.Put([]byte(strconv.Itoa(k)), v)
        if e != nil { return e }

        return nil
    })
    return
}

// Return JSON encoded byte slice from DB
func rdb(db *bolt.DB, k int, cbuc []byte) (v []byte, e error) {

    e = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(cbuc)
        if b == nil { return fmt.Errorf("No bucket!") }

        v = b.Get([]byte(strconv.Itoa(k)))
        return nil
    })
    return
}

func txtreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) string {

    e := r.ParseForm()
    cherr(e)

    return r.FormValue("request")
}

func main() {

    rand.Seed(time.Now().UnixNano())

	pptr := flag.Int("p", DEFAULTPORT, "port number to listen")
	dbptr := flag.String("d", DBNAME, "specify database to open")
	flag.Parse()

    db, e := bolt.Open(*dbptr, 0640, nil)
    cherr(e)
    defer db.Close()

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    http.HandleFunc("/gettext", func(w http.ResponseWriter, r *http.Request) {
        data := txtreqhandler(w, r, db)
        fmt.Printf("DEBUG: %+v\n", data)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    cherr(e)
}
