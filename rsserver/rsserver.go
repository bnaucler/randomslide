package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "flag"
    "sort"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

// Initialize logger
func initlog(prgname string) {

    logfile := fmt.Sprintf("%s/%s.log", rscore.LOGPATH, prgname)

    f, e := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    rscore.Cherr(e)

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Creates valid selection list from tags
func mksel(db *bolt.DB, tags []string, buc []byte) []int {

    var sel []int
    ctags := rscore.Tag{}

    for _, t := range tags {
        bt := []byte(t)
        mtags, e := rsdb.Rdb(db, bt, buc)
        rscore.Cherr(e)

        json.Unmarshal(mtags, &ctags)
        sel = append(sel, ctags.Ids...)
    }

    return sel
}

// Returns a random key based on tag list
func getkeyfromsel(db *bolt.DB, tags []string, buc []byte, kmax int) []byte {

    sel := mksel(db, tags, buc)
    smax := len(sel)

    var k []byte

    if smax > 0 {
        ki := rand.Intn(smax)
        k = []byte(strconv.Itoa(sel[ki]))

    } else {
        ki := rand.Intn(kmax)
        k = []byte(strconv.Itoa(ki))
    }

    return k
}

// Sends random text object from database, based on requested tags
func getrndtextobj(db *bolt.DB, kmax int, tags []string, buc []byte) string {


    if kmax < 2 { return "" }

    k := getkeyfromsel(db, tags, buc, kmax)

    txt := rscore.Textobj{}
    mtxt, e := rsdb.Rdb(db, k, buc)
    rscore.Cherr(e)
    e = json.Unmarshal(mtxt, &txt)
    rscore.Cherr(e)

    return txt.Text
}

// Sends random image url from database, based on requested tags
func getrndimgobj(db *bolt.DB, kmax int, tags []string, buc []byte) string {

    if kmax < 2 { return "" }

    k := getkeyfromsel(db, tags, buc, kmax)

    img := rscore.Imgobj{}
    mimg, e := rsdb.Rdb(db, k, buc)
    rscore.Cherr(e)
    e = json.Unmarshal(mimg, &img)
    rscore.Cherr(e)

    return img.Fname
}

// Returns a full slide deck according to request
func mkdeck(req rscore.Deckreq, db *bolt.DB,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    deck := rscore.Deck{
            N: req.N,
            Lang: req.Lang }

    for i := 0; i < req.N; i++ {
        slide := rscore.Slide{
            Title: getrndtextobj(db, settings.Tmax, req.Tags, rscore.TBUC),
            Btext: getrndtextobj(db, settings.Bmax, req.Tags, rscore.BBUC),
            Imgur: getrndimgobj(db, settings.Imax, req.Tags, rscore.IBUC) }

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

    tags := gettagsfromreq(r)
    req := rscore.Deckreq{
            N: n,
            Lang: r.FormValue("lang"),
            Tags: tags }

    mreq, e := json.Marshal(req)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, mreq, r)

    deck, settings := mkdeck(req, db, settings)

    mdeck, e := json.Marshal(deck)
    rscore.Addlog(rscore.L_RESP, mdeck, r)

    enc := json.NewEncoder(w)
    enc.Encode(deck)

    return settings
}

// Updates index to include new tags
func addtagstoindex(tags []string, settings rscore.Settings,
    w http.ResponseWriter) rscore.Settings {

    r := 0

    for _, t := range tags {
        if rscore.Findstrinslice(t, settings.Taglist) == false {
            settings.Taglist = append(settings.Taglist, t)
            r++
        }
    }

    if r != 0 { sort.Strings(settings.Taglist) }

    var sstr string
    if r != 0 { sstr = fmt.Sprintf("%d new tag(s) added", r) }
    rscore.Sendstatus(rscore.C_OK, sstr, w)

    return settings
}

// Updates all relevant tag lists
func updatetaglists(db *bolt.DB, tags []string, i int, buc []byte) {

    for _, s := range tags {
        ctag := rscore.Tag{}
        key := []byte(s)

        resp, e := rsdb.Rdb(db, key, buc)
        rscore.Cherr(e)

        json.Unmarshal(resp, &ctag)
        ctag.Ids = append(ctag.Ids, i)

        dbw, e := json.Marshal(ctag)
        e = rsdb.Wrdb(db, key, dbw, buc)
        rscore.Cherr(e)
    }
}

// Conditionally adds tagged text to database
func addtextwtags(text string, tags []string, db *bolt.DB,
    mxindex int, buc []byte) {

    to := rscore.Textobj{
            Id: mxindex,
            Text: text,
            Tags: tags }

    // Storing the object in db
    key := []byte(strconv.Itoa(mxindex))
    mtxt, e := json.Marshal(to)
    e = rsdb.Wrdb(db, key, mtxt, buc)
    rscore.Cherr(e)

    // Update all relevant tag lists
    updatetaglists(db, tags, mxindex, buc)
}

// Returns a slice of cleaned tags from http request
func gettagsfromreq(r *http.Request) []string {

    var ret []string
    rtags := r.FormValue("tags")
    itags := strings.Split(rtags, " ")

    for _, s := range itags {
        ret = append(ret, rscore.Cleanstring(s))
    }

    return ret
}

// Handles incoming requests to add images
func imgreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    r.ParseMultipartForm(10 << 20)
    f, hlr, e := r.FormFile("file")
    e = rscore.Cherr(e)
    if e != nil { return settings }
    defer f.Close()

    mt := hlr.Header["Content-Type"][0]

    if rscore.Findstrinslice(mt, rscore.IMGMIME) == false {
        rscore.Sendstatus(rscore.C_WRFF,
            "Unknown image format - file not uploaded", w)
        return settings
    }

    lmsg := fmt.Sprintf("File: %+v(%s) - %+v",
        hlr.Filename, rscore.Prettyfsize(hlr.Size), mt)
    rscore.Addlog(rscore.L_REQ, []byte(lmsg), r)

    ext := filepath.Ext(hlr.Filename)
    fformat := fmt.Sprintf("img-*%s", ext)
    tmpf, e := ioutil.TempFile("static/img", fformat)
    rscore.Cherr(e)
    defer tmpf.Close()

    fc, e := ioutil.ReadAll(f)
    rscore.Cherr(e)
    tmpf.Write(fc)

    tags := gettagsfromreq(r)
    ofn := filepath.Base(tmpf.Name())
    iobj := rscore.Imgobj{
        Id: settings.Imax,
        Fname: ofn,
        Tags: tags }

    mobj, e := json.Marshal(iobj)
    k := []byte(strconv.Itoa(settings.Imax))
    e = rsdb.Wrdb(db, k, mobj, rscore.IBUC)
    rscore.Cherr(e)

    settings = addtagstoindex(tags, settings, w)
    updatetaglists(db, tags, settings.Imax, rscore.IBUC)

    settings.Imax++
    rsdb.Wrsettings(db, settings)

    rscore.Sendstatus(rscore.C_OK, "", w)

    return settings
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    tags := gettagsfromreq(r)

    tr := rscore.Textreq{
            Ttext: r.FormValue("ttext"),
            Btext: r.FormValue("btext"),
            Tags: tags }

    ltxt, e := json.Marshal(tr)
    rscore.Addlog(rscore.L_REQ, ltxt, r)

    settings = addtagstoindex(tags, settings, w)

    if len(tr.Ttext) > 1 {
        addtextwtags(tr.Ttext, tags, db, settings.Tmax, rscore.TBUC)
        settings.Tmax++
    }

    if len(tr.Btext) > 1 {
        addtextwtags(tr.Btext, tags, db, settings.Bmax, rscore.BBUC)
        settings.Bmax++
    }

    rsdb.Wrsettings(db, settings)

    return settings
}

// Handles incoming requests for tag index
func tagreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    resp := rscore.Tagresp{}
    var ttag rscore.Rtag

    for _, t := range settings.Taglist {
        ttag.Name = t

        if settings.Tmax > 0 {
            ttag.TN = rsdb.Countobj(db, t, rscore.TBUC)
        }

        if settings.Bmax > 0 {
            ttag.BN = rsdb.Countobj(db, t, rscore.BBUC)
        }

        if settings.Imax > 0 {
            ttag.IN = rsdb.Countobj(db, t, rscore.IBUC)
        }

        resp.Tags = append(resp.Tags, ttag)
    }

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for shutdowns
func shutdownhandler (w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    rscore.Addlog(rscore.L_SHUTDOWN, []byte(""), r)
    rscore.Sendstatus(rscore.C_OK, "", w)

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

    settings.Pidfile = fmt.Sprintf("%s%s.pid", rscore.PIDFILEPATH, prgname)
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

    // Slide deck requests
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

    // Upload images
    http.HandleFunc("/addimg", func(w http.ResponseWriter, r *http.Request) {
        settings = imgreqhandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e = http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}
