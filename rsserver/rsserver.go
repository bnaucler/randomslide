package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "flag"
    "sort"
    "image"
    "image/png"
    "image/jpeg"
    "image/gif"
    "strings"
    "strconv"
    "net/http"
    "os/signal"
    "math/rand"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func init() {
    rand.Seed(time.Now().UnixNano())

    image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
    image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

}

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
func getrndimg(db *bolt.DB, kmax int, tags []string, buc []byte) rscore.Imgobj {

    if kmax < 2 { return rscore.Imgobj{} }

    k := getkeyfromsel(db, tags, buc, kmax)

    img := rscore.Imgobj{}
    mimg, e := rsdb.Rdb(db, k, buc)
    rscore.Cherr(e)
    e = json.Unmarshal(mimg, &img)
    rscore.Cherr(e)

    return img
}

// Returns deck from database
func getdeckfdb(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq,
    settings rscore.Settings) rscore.Deck {

    var k int

    if req.Id >= settings.Dmax {
        return rscore.Deck{}

    } else if settings.Dmax < 1 {
        return rscore.Deck{}

    } else {
        k = req.Id
    }

    bk := []byte(strconv.Itoa(k))
    mdeck, e := rsdb.Rdb(db, bk, rscore.DBUC)
    rscore.Cherr(e)

    e = json.Unmarshal(mdeck, &deck)
    rscore.Cherr(e)

    return deck
}

// Returns a new slide deck according to request
func mkdeck(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    for i := 0; i < req.N; i++ {
        slide := rscore.Slide{
            Type: rand.Intn(6), // TODO
            Title: getrndtextobj(db, settings.Tmax, req.Tags, rscore.TBUC),
            Btext: getrndtextobj(db, settings.Bmax, req.Tags, rscore.BBUC),
            Img: getrndimg(db, settings.Imax, req.Tags, rscore.IBUC) }

        deck.Slides = append(deck.Slides, slide)
    }

    deck.Id = settings.Dmax

    k := []byte(strconv.Itoa(deck.Id))
    mdeck, e := json.Marshal(deck)
    rscore.Cherr(e)
    e = rsdb.Wrdb(db, k, mdeck, rscore.DBUC)

    settings.Dmax++
    rsdb.Wrsettings(db, settings)

    return deck, settings
}

// Sets basic params & determines if new deck should be built
func getdeck(req rscore.Deckreq, db *bolt.DB,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    deck := rscore.Deck{
            Id: req.Id,
            N: req.N,
            Lang: req.Lang }

    if req.Isidreq {
        deck = getdeckfdb(db, deck, req, settings)

    } else {
        deck, settings = mkdeck(db, deck, req, settings)

    }

    return deck, settings
}

// Determines if deck requests specific id
func isidreq(r *http.Request) (int, bool) {

    fvid := r.FormValue("id")

    if len(fvid) < 1 { return 0, false }

    id, e := strconv.Atoi(fvid)

    if e == nil { return id, true }

    return 0, false
}

// Handles incoming requests for decks
func deckreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)
    var n int

    fvam := r.FormValue("amount")

    if len(fvam) < 1 {
        n = 0

    } else {
        n, e = strconv.Atoi(fvam)
        rscore.Cherr(e)
    }

    id, isidr := isidreq(r)
    tags := gettagsfromreq(r)

    req := rscore.Deckreq{
            Id: id,
            Isidreq: isidr,
            N: n,
            Lang: r.FormValue("lang"),
            Tags: tags }

    mreq, e := json.Marshal(req)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, mreq, r)

    deck, settings := getdeck(req, db, settings)

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

// Conditionally returns image size type & true if fitting classification
func getimgtype(w int, h int) (int, bool) {

    i := 3

    // Upper bounds check
    if w > rscore.IMGMAX[3][0] || h > rscore.IMGMAX[3][1] {
        return 0, false
    }

    for i >= 0 {
        if w > rscore.IMGMIN[i][0] && h > rscore.IMGMIN[i][1] &&
           w < rscore.IMGMAX[i][0] && h < rscore.IMGMAX[i][1] {
               return i, true
           }
        i--
    }

    return 0, false
}

// Stores image object in database
func addimgwtags(db *bolt.DB, fn string, iw int, ih int, tags []string,
    w http.ResponseWriter, settings rscore.Settings) rscore.Settings {

    ofn := filepath.Base(fn)

    isz, szok := getimgtype(iw, ih)

    fmt.Printf("DEBUG: %d:%d - %+v, %+v\n", iw, ih, isz, szok)

    if szok == false {
        rscore.Sendstatus(rscore.C_WRSZ,
            "Could not classify size - file not uploaded", w)
        return settings
    }

    iobj := rscore.Imgobj{
        Id: settings.Imax,
        Fname: ofn,
        Tags: tags,
        W: iw,
        H: ih,
        Size: isz }

    fmt.Printf("DEBUG: %+v\n", iobj)

    // Add object to db
    mobj, e := json.Marshal(iobj)
    k := []byte(strconv.Itoa(settings.Imax))
    e = rsdb.Wrdb(db, k, mobj, rscore.IBUC)
    rscore.Cherr(e)

    // Update relevant tags
    settings = addtagstoindex(tags, settings, w)
    updatetaglists(db, tags, settings.Imax, rscore.IBUC)
    settings.Imax++

    return settings
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
    tmpf, e := ioutil.TempFile(rscore.IMGDIR, fformat)
    rscore.Cherr(e)
    defer tmpf.Close()

    ic, _, e := image.DecodeConfig(f)
    rscore.Cherr(e)

    fc, e := ioutil.ReadAll(f)
    rscore.Cherr(e)
    tmpf.Write(fc)

    tags := gettagsfromreq(r)
    settings = addimgwtags(db, tmpf.Name(), ic.Width, ic.Height, tags, w, settings)
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

    if len(tr.Ttext) > 1 && len(tr.Ttext) < rscore.TTEXTMAX {
        addtextwtags(tr.Ttext, tags, db, settings.Tmax, rscore.TBUC)
        settings.Tmax++
    }

    if len(tr.Btext) > 1 && len(tr.Btext) < rscore.BTEXTMAX {
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

// Removes all files residing in dir
func rmallfromdir(dir string) {

    d, e := os.Open(dir)
    rscore.Cherr(e)
    defer d.Close()

    fl, e := d.Readdirnames(-1)
    rscore.Cherr(e)

    for _, fn := range fl {
        if fn == ".gitkeep" { continue }
        e = os.RemoveAll(filepath.Join(dir, fn))
        rscore.Cherr(e)
    }
}

func shutdown(settings rscore.Settings) {

    go func() {
        time.Sleep(1 * time.Second)
        os.Remove(settings.Pidfile)
        os.Exit(0)
    }()
}

// Handles incoming requests for shutdowns
func shutdownhandler (w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    wipe := r.FormValue("wipe")

    rscore.Addlog(rscore.L_SHUTDOWN, []byte(""), r)
    rscore.Sendstatus(rscore.C_OK, "", w)

    rsdb.Wrsettings(db, settings)

    if wipe == "yes" {
        db.Close()
        os.Remove(rscore.DBNAME)
        rmallfromdir(rscore.IMGDIR)
    }

    shutdown(settings)
}

// Creates pid file and starts logging
func rsinit(settings rscore.Settings) rscore.Settings {

    prgname := filepath.Base(os.Args[0])
    pid := os.Getpid()

    settings.Pidfile = fmt.Sprintf("%s%s.pid", rscore.PIDFILEPATH, prgname)
    e := ioutil.WriteFile(settings.Pidfile, []byte(strconv.Itoa(pid)), 0644)
    rscore.Cherr(e)

    initlog(prgname)

    // SIGINT handler
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func(){
        for sig := range sigc {
            fmt.Printf("Caught %+v - cleaning up.\n", sig)
            shutdown(settings)
        }
    }()

    if settings.Verb { log.Printf("%s started with PID: %d\n", prgname, pid) }

    return settings
}

func main() {

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
