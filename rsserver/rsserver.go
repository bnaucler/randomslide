package main

import (
    "os"
    "fmt"
    "net"
    "log"
    "time"
    "flag"
    "sort"
    "regexp"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/rsserver/rscore"
    "github.com/bnaucler/randomslide/rsserver/rsdb"
)

// Retrieves client IP address from http request
func getclientip(r *http.Request) string {

    ip, _, e := net.SplitHostPort(r.RemoteAddr)
    rscore.Cherr(e)

    return ip
}

// Log file wrapper
// TODO: Use interface() instead of []byte and ltype
//       log levels
func addlog(ltype int, msg []byte, r *http.Request) {

    ip := getclientip(r)
    var lentry string

    switch ltype {
        case rscore.L_REQ:
            lentry = fmt.Sprintf("REQ from %s: %s", ip, msg)

        case rscore.L_RESP:
            lentry = fmt.Sprintf("RESP to %s: %s", ip, msg)

        case rscore.L_SHUTDOWN:
            lentry = fmt.Sprintf("Server shutdown requested from %s", ip)

        default:
            lentry = fmt.Sprintf("Something happened, but I don't know how to log it!")
    }

    log.Println(lentry)
}

// Initialize logger
func initlog(prgname string) {

    logfile := fmt.Sprintf("%s/%s.log", rscore.LOGPATH, prgname)

    f, e := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    rscore.Cherr(e)

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Sends random text object from database
func getrndtextobj(db *bolt.DB, kmax int) string {

    txtreq := rscore.Textobj{}
    k := []byte(strconv.Itoa(rand.Intn(kmax)))

    mtxt, e := rsdb.Rdb(db, k, rscore.TBUC)
    rscore.Cherr(e)

    json.Unmarshal(mtxt, &txtreq)

    return txtreq.Text
}

// Returns a full slide deck according to request
func mkdeck(req rscore.Deckreq, db *bolt.DB,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    deck := rscore.Deck{
            N: req.N,
            Lang: req.Lang }

    for i := 0; i < req.N; i++ {
        slide := rscore.Slide{
            Title: getrndtextobj(db, settings.Cid),
            Btext: getrndtextobj(db, settings.Cid) }

        deck.Slides = append(deck.Slides, slide)
    }

    return deck, settings
}

// Handles incoming requests for decks
func deckreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    n, e := strconv.Atoi(r.FormValue("amount"))
    rscore.Cherr(e)

    req := rscore.Deckreq{
            N: n,
            Lang: r.FormValue("lang"),
            Tags: r.FormValue("category") }

    mreq, e := json.Marshal(req)
    rscore.Cherr(e)
    addlog(rscore.L_REQ, mreq, r)

    deck, settings := mkdeck(req, db, settings)

    mdeck, e := json.Marshal(deck)
    addlog(rscore.L_RESP, mdeck, r)

    enc := json.NewEncoder(w)
    enc.Encode(deck)

    return settings
}

// Removes whitespace and special characters from string
func cleanstring(src string) string {

    rx, e := regexp.Compile("[^a-z]+")
    rscore.Cherr(e)

    dst := rx.ReplaceAllString(src, "")

    return dst
}

// Returns true if string is present in list
func findstrinslice(v string, list []string) bool {

    for _, t := range list {
        if v == t { return true }
    }

    return false
}

// Updates index to include new tags
func addtagstoindex(tags []string, settings rscore.Settings) (rscore.Settings, int) {

    r := 0

    for _, t := range tags {
        if findstrinslice(t, settings.Taglist) == false {
            settings.Taglist = append(settings.Taglist, t)
            r++
        }
    }

    if r != 0 { sort.Strings(settings.Taglist) }

    return settings, r
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    tr := rscore.Textreq{
            Text: r.FormValue("text"),
            Tags: r.FormValue("tags") }

    ltxt, e := json.Marshal(tr)
    addlog(rscore.L_REQ, ltxt, r)

    itags := strings.Split(tr.Tags, " ")

    to := rscore.Textobj{
            Id: settings.Cid,
            Text: tr.Text }

    for _, s := range itags {
        to.Tags = append(to.Tags, cleanstring(s))
    }

    for _, s := range to.Tags {
        key := []byte(s)
        mtxt, e := json.Marshal(to)
        e = rsdb.Wrdb(db, key, mtxt, rscore.TBUC)
        rscore.Cherr(e)
    }

    settings, cngs := addtagstoindex(to.Tags, settings)
    rsdb.Wrsettings(db, settings)

    var sstr string
    if cngs != 0 { sstr = fmt.Sprintf("%d new tags added", cngs) }
    sendstatus(rscore.C_OK, sstr, w)

    return settings
}

// Sends a status code response as JSON object
func sendstatus(code int, text string, w http.ResponseWriter) {

    resp := rscore.Statusresp{
            Code: code,
            Text: text }

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for tag index
func tagreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    resp := rscore.Tagresp{
            Tags: settings.Taglist }

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    addlog(rscore.L_RESP, mresp, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for shutdowns
func shutdownhandler (w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    addlog(rscore.L_SHUTDOWN, []byte(""), r)
    sendstatus(rscore.C_OK, "", w)

    rsdb.Wrsettings(db, settings)

    go func() {
        time.Sleep(1 * time.Second)
        os.Remove(settings.Pidfile)
        os.Exit(0)
    }()
}

// Creates pid file and starts logging
func rsinit(settings rscore.Settings) rscore.Settings {

    prgname := filepath.Base(os.Args[0])
    pid := os.Getpid()

    settings.Pidfile = fmt.Sprintf("%s/%s.pid", rscore.PIDFILEPATH, prgname)
    e := ioutil.WriteFile(settings.Pidfile, []byte(strconv.Itoa(pid)), 0644)
    rscore.Cherr(e)

    initlog(prgname)

    if settings.Verb { log.Printf("%s started with PID: %d\n", prgname, pid) }

    return settings
}

func main() {

    rand.Seed(time.Now().UnixNano())

    pptr := flag.Int("p", rscore.DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", false, "verbose mode")
    flag.Parse()

    db, e := bolt.Open(*dbptr, 0640, nil)
    rscore.Cherr(e)
    defer db.Close()

    settings := rsdb.Rsettings(db)
    settings.Verb = *vptr
    settings = rsinit(settings)

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    if rscore.VOLATILEMODE == true {
        http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
            shutdownhandler(w, r, db, settings)
        })
    }

    // Slide requests
    http.HandleFunc("/getdeck", func(w http.ResponseWriter, r *http.Request) {
        settings = deckreqhandler(w, r, db, settings)
    })

    // Tags requests
    http.HandleFunc("/gettags", func(w http.ResponseWriter, r *http.Request) {
        tagreqhandler(w, r, db, settings)
    })

    // Add text to db
    http.HandleFunc("/addtext", func(w http.ResponseWriter, r *http.Request) {
        settings = textreqhandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}
