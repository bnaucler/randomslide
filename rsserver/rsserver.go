package main

import (
    "fmt"
    "log"
    "time"
    "flag"
    "strconv"
    "net/http"
    "math/rand"
    "encoding/json"

    "github.com/boltdb/bolt"
)

const DEFAULTPORT = 6291
const DBNAME = "./db/random.db"

var tbuc = []byte("pbuc")       // text bucket
var ibuc = []byte("ibuc")       // image bucket
var sbuc = []byte("sbuc")       // settings bucket

type Resp struct {
    Data string
    Id int
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

// Handling requests for text objects
func txtreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB, cid int) int {

    e := r.ParseForm()
    cherr(e)

    key := []byte(strconv.Itoa(cid))

    e = wrdb(db, key, []byte(r.FormValue("request")), tbuc)
    cherr(e)

    v, e := rdb(db, key, tbuc)
    cherr(e)

    resp := Resp{
        Data: string(v),
        Id: cid}

    enc := json.NewEncoder(w)
    enc.Encode(resp)

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

    cid := 0

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Text requests
    http.HandleFunc("/gettext", func(w http.ResponseWriter, r *http.Request) {
        cid = txtreqhandler(w, r, db, cid)
        fmt.Printf("DEBUG: %+v\n", cid)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    cherr(e)
}
