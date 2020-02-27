package main

import (
    "os"
    "fmt"
    "time"
    "flag"
    "sort"
    "image"
    "bytes"
    "image/png"
    "image/jpeg"
    "image/gif"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/boltdb/bolt"
    "golang.org/x/crypto/bcrypt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func init() {
    rand.Seed(time.Now().UnixNano())

    image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
    image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
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

// Determines an appropriate slide type to generate
func setslidetype(i int) rscore.Slidetype {

    st := rscore.Slidetype{}

    // We always start with type 0 (big title)
    if i == 0 {
        st.Type = 0

    } else {
        st.Type = rand.Intn(rscore.STYPES)
    }

    // TODO: Make proper index objects
    switch st.Type {

    case 0: // Big title
        st.TT = true
        st.BT = false
        st.IMG = true

    case 1: // Full screen picture
        st.TT = false
        st.BT = false
        st.IMG = true

    case 2: // Big number
        st.TT = false
        st.BT = false
        st.IMG = false

    case 3: // Bullet point list
        st.TT = false
        st.BT = false
        st.IMG = true

    case 4: // Title, body & img
        st.TT = true
        st.BT = true
        st.IMG = true

    case 5: // Inspirational quote
        st.TT = true
        st.BT = false
        st.IMG = false

    case 6: // Picture with text
        st.TT = false
        st.BT = true
        st.IMG = true

    }

    return st
}

// Generate bullet point list for slide type 3
func bpgen(db *bolt.DB, tags []string, settings rscore.Settings) []string {

    var bps []string
    var bp string

    rnd := rscore.BPMAX - rscore.BPMIN
    n := rand.Intn(rnd)
    n += rscore.BPMIN

    for i := 0 ; i < n ; i ++ {
        for len(bp) == 0 || len(bp) > rscore.BPOINTMAX {
            bp = getrndtextobj(db, settings.Tmax, tags, rscore.TBUC)
        }
        bps = append(bps, bp)
        bp = ""
    }

    return bps
}

// Flips a coin, returns random bool
func coin() bool {

    rnd := rand.Intn(2)

    if rnd == 0 { return false }
    return true
}

// Generate numbers of slide type 2
func numgen() string {

    p := byte(' ')
    s := byte(' ')

    n := rand.Intn(rscore.RNUMBMAX)

    for i := 0; i < rscore.RNUMEMAX; i ++ {
        if coin() { n *= 10 }
    }

    plen := len(rscore.NUMPREF)
    if coin() { p = rscore.NUMPREF[rand.Intn(plen)] }

    slen := len(rscore.NUMSUFF)
    if coin() { s = rscore.NUMSUFF[rand.Intn(slen)] }

    ret := fmt.Sprintf("%s%d%s", string(p), n, string(s))

    return strings.TrimSpace(ret)
}

// Retrieves relevant data based on slide type
func getslide(db *bolt.DB, st rscore.Slidetype, settings rscore.Settings,
    req rscore.Deckreq) rscore.Slide {

    slide := rscore.Slide{Type: st.Type}

    if st.TT { slide.Title = getrndtextobj(db, settings.Tmax, req.Tags, rscore.TBUC) }
    if st.BT { slide.Btext = getrndtextobj(db, settings.Bmax, req.Tags, rscore.BBUC) }
    if st.IMG { slide.Img = getrndimg(db, settings.Imax, req.Tags, rscore.IBUC) }

    if st.Type == 2 {
        slide.Title = numgen()
        fmt.Println(slide.Title)

    } else if st.Type == 3 {
        slide.Bpts = bpgen(db, req.Tags, settings)
    }

    return slide
}

// Returns a new slide deck according to request
func mkdeck(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    for i := 0; i < req.N; i++ {
        st := setslidetype(i)
        slide := getslide(db, st, settings, req)
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

// Wrapper for tag status responses
func sendtagstatus(r int, w http.ResponseWriter) {

    var sstr string
    if r != 0 { sstr = fmt.Sprintf("%d new tag(s) added", r) }
    rscore.Sendstatus(rscore.C_OK, sstr, w)
}

// Returns a slice of cleaned tags from http request
func gettagsfromreq(r *http.Request) []string {

    var ret []string
    rtags := r.FormValue("tags")
    itags := strings.Split(rtags, " ")

    for _, s := range itags {
        ret = append(ret, rscore.Cleanstring(s, rscore.RXTAGS))
    }

    return ret
}

// Conditionally returns image size type & true if fitting classification
func getimgtype(w int, h int) (int, bool) {

    i := 3

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

    // Add object to db
    mobj, e := json.Marshal(iobj)
    k := []byte(strconv.Itoa(settings.Imax))
    e = rsdb.Wrdb(db, k, mobj, rscore.IBUC)
    rscore.Cherr(e)

    // Update relevant tags
    nt, settings := rsdb.Tagstoindex(tags, settings)
    sendtagstatus(nt, w)
    rsdb.Updatetaglists(db, tags, settings.Imax, rscore.IBUC)
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

    fc, e := ioutil.ReadAll(f)
    rscore.Cherr(e)
    tmpf.Write(fc)

    ibuf, e := ioutil.ReadFile(tmpf.Name())
    rscore.Cherr(e)
    fszr := bytes.NewReader(ibuf)
    ic, _, e := image.DecodeConfig(fszr)
    rscore.Cherr(e)

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
            Bpoint: r.FormValue("bpoint"),
            Tags: tags }

    ltxt, e := json.Marshal(tr)
    rscore.Addlog(rscore.L_REQ, ltxt, r)

    nt, settings := rsdb.Tagstoindex(tags, settings)
    sendtagstatus(nt, w)

    if len(tr.Ttext) > 1 && len(tr.Ttext) < rscore.TTEXTMAX {
        rsdb.Addtextwtags(tr.Ttext, tags, db, settings.Tmax, rscore.TBUC)
        settings.Tmax++
    }

    if len(tr.Btext) > 1 && len(tr.Btext) < rscore.BTEXTMAX {
        rsdb.Addtextwtags(tr.Btext, tags, db, settings.Bmax, rscore.BBUC)
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

    wipe := r.FormValue("wipe")

    rscore.Addlog(rscore.L_SHUTDOWN, []byte(""), r)
    rscore.Sendstatus(rscore.C_OK, "", w)

    rsdb.Wrsettings(db, settings)

    if wipe == "yes" {
        db.Close()
        os.Remove(rscore.DBNAME)
        rscore.Rmall(rscore.IMGDIR)
    }

    rscore.Shutdown(settings)
}

// Returns true if key returns something from database
func isindb(db *bolt.DB, k []byte, buc []byte) bool {

    v, e := rsdb.Rdb(db, k, buc)

    if len(v) == 0 || e != nil { return false }
    return true
}

// Write user index to database
func updateuserindex(db *bolt.DB, uname string,
    settings rscore.Settings) rscore.Settings {

    users := rscore.Uindex{}
    var mindex []byte

    if settings.Umax > 0 {
        mindex, e := rsdb.Rdb(db, rscore.INDEX, rscore.UBUC)
        rscore.Cherr(e)
        e = json.Unmarshal(mindex, &users)
        rscore.Cherr(e)
    }

    users.Names = append(users.Names, uname)
    sort.Strings(users.Names)

    mindex, e := json.Marshal(users)

    fmt.Printf("DEBUG: %+v\n", users.Names)

    e = rsdb.Wrdb(db, rscore.INDEX, mindex, rscore.UBUC)
    rscore.Cherr(e)

    settings.Umax++
    return settings

}

// Handles incoming requests for user registrations
func reghandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    u := rscore.User{}
    ur := rscore.Uresp{}

    user := r.FormValue("user")
    pass := r.FormValue("pass")

    hash, e := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    rscore.Cherr(e)

    // Username already takem - registration not possible
    if settings.Umax > 0 {
        if isindb(db, []byte(user), rscore.UBUC) {
            rscore.Sendstatus(rscore.C_UIDB, "Username already in db", w)
            return settings
        }
    }

    // Username includes illegal characters
    if user != rscore.Cleanstring(user, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return settings
    }

    u.Name = user
    u.Skey = rscore.Randstr(rscore.SKEYLEN)
    ur.Name = u.Name
    ur.Skey = u.Skey
    u.Pass = hash

    settings = updateuserindex(db, user, settings)

    mu, e := json.Marshal(u)
    rscore.Cherr(e)
    e = rsdb.Wrdb(db, []byte(user), mu, rscore.UBUC)
    rscore.Cherr(e)

    mr, e := json.Marshal(ur)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mr, r)

    enc := json.NewEncoder(w)
    enc.Encode(ur)

    rsdb.Wrsettings(db, settings)
    return settings
}

func main() {

    pptr := flag.Int("p", rscore.DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", false, "verbose mode")
    flag.Parse()

    db := rsdb.Open(*dbptr)
    defer db.Close()

    settings := rsdb.Rsettings(db)
    settings.Verb = *vptr
    settings = rscore.Rsinit(settings)

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

    // User registration
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        settings = reghandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e := http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}
