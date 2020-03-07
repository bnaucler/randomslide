package main

import (
    "os"
    "fmt"
    "flag"
    "bytes"
    "image"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "encoding/json"

    "github.com/boltdb/bolt"
    "golang.org/x/crypto/bcrypt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
    "github.com/bnaucler/randomslide/lib/rsimage"
)

// Updates slilde probabilities
func setsprob(n int, sprob []int) []int {

    if sprob[n] > 1 {
        for i := range sprob {
            if i == n {
                sprob[i]--

            } else {
                if coin() { sprob[i]++ }

            }
        }
    }

    return sprob
}

// Sets a random slide type based on current probabilities
func rndslidetype(i int, sprob []int) (int, []int) {

    // We always start with a title slide
    if i == 0 { return 0, sprob }

    tot := 0
    for _, v := range sprob { tot += v }

    target := rand.Intn(tot)
    p := 0
    n := 0

    for {
        p += sprob[n]
        if p >= target { break }
        n++
    }

    sprob = setsprob(n, sprob)

    return n, sprob
}

// Determines an appropriate slide type to generate
func setslidetype(i int, sprob []int) (rscore.Slidetype, []int) {

    st := rscore.Slidetype{}

    // We always start with type 0 (big title)
    st.Type, sprob = rndslidetype(i, sprob)

    // TODO: Make proper index objects
    switch st.Type {

    case 0: // Big title
        st.TT = true
        st.BT = false

    case 1: // Full screen picture
        st.TT = false
        st.BT = false

    case 2: // Big number
        st.TT = false
        st.BT = false

    case 3: // Bullet point list
        st.TT = false
        st.BT = false

    case 4: // Title, body & img
        st.TT = true
        st.BT = true

    case 5: // Inspirational quote
        st.TT = true
        st.BT = false

    case 6: // Picture with text
        st.TT = false
        st.BT = true

    case 7: // Graph
        st.TT = true
        st.BT = false
    }

    st.IMG = rscore.ISZINDEX[st.Type]

    return st, sprob
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
            bp = rsdb.Getrndtxt(db, settings.Tmax, tags, rscore.TBUC)
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

// Flips coins to get random exponent
func rndexp(n int) int {

    for i := 0; i < rscore.RNUMEMAX; i++ {
        if coin() { n *= 10 }
    }

    return n
}

// Returns number with fixed exponent
func setexp(n int, e int) int {

    for i := 0; i < e; i ++ {
        n *= 10
    }

    return n
}

// Generate numbers of slide type 2
func numgen() string {

    p := byte(' ')
    s := byte(' ')

    b := rand.Intn(rscore.RNUMBMAX)
    n := rndexp(b)

    plen := len(rscore.NUMPREF)
    if coin() { p = rscore.NUMPREF[rand.Intn(plen)] }

    slen := len(rscore.NUMSUFF)
    if coin() { s = rscore.NUMSUFF[rand.Intn(slen)] }

    ret := fmt.Sprintf("%s%d%s", string(p), n, string(s))

    return strings.TrimSpace(ret)
}

func dpgen() []int {

    var ret []int

    dp := rand.Intn(4)
    dp += 2
    e := rand.Intn(rscore.RNUMEMAX)

    for i := 0; i < dp; i++ {
        b := rand.Intn(rscore.RNUMBMAX)
        n := setexp(b, e)
        ret = append(ret, n)
    }

    return ret
}

// Retrieves relevant data based on slide type
func getslide(db *bolt.DB, st rscore.Slidetype, settings rscore.Settings,
    req rscore.Deckreq) rscore.Slide {

    slide := rscore.Slide{Type: st.Type}

    if st.TT { slide.Title = rsdb.Getrndtxt(db, settings.Tmax, req.Tags, rscore.TBUC) }
    if st.BT { slide.Btext = rsdb.Getrndtxt(db, settings.Bmax, req.Tags, rscore.BBUC) }
    if len(st.IMG) > 0 {
        suf := rsimage.Mkimgsuflist(st.Type)
        stags := rscore.Addtagsuf(req.Tags, suf)
        slide.Img = rsdb.Getrndimg(db, settings.Imax, stags, rscore.IBUC)
    }

    switch st.Type {

    case 2:
        slide.Title = numgen()

    case 3:
        slide.Bpts = bpgen(db, req.Tags, settings)

    case 7:
        slide.Dpts = dpgen()
        slide.Ctype = rand.Intn(3)
    }

    return slide
}

// Returns a new slide deck according to request
func mkdeck(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq,
    settings rscore.Settings) (rscore.Deck, rscore.Settings) {

    var st rscore.Slidetype
    sprob := rscore.SPROB

    for i := 0; i < req.N; i++ {
        st, sprob = setslidetype(i, sprob)
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
        deck = rsdb.Getdeckfdb(db, deck, req, settings)

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
    tags := rscore.Formattags(r.FormValue("tags"))

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

// Handles incoming requests to add images
func imgreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    r.ParseMultipartForm(10 << 20)
    sf, hlr, e := r.FormFile("file")
    e = rscore.Cherr(e)
    if e != nil { return settings }
    defer sf.Close()

    mt := hlr.Header["Content-Type"][0]

    if rscore.Findstrinslice(mt, rscore.IMGMIME) == false {
        rscore.Sendstatus(rscore.C_WRFF,
            "Unknown image format - file not uploaded", w)
        return settings
    }

    lmsg := fmt.Sprintf("File: %+v(%s) - %+v",
        hlr.Filename, rscore.Prettyfsize(hlr.Size), mt)
    rscore.Addlog(rscore.L_REQ, []byte(lmsg), r)

    fn, fnp := rsimage.Newimagepath(hlr.Filename)

    e = rscore.Wrdatafile(fnp, sf)
    rscore.Cherr(e)

    ibuf, e := ioutil.ReadFile(fnp)
    rscore.Cherr(e)

    fszr := bytes.NewReader(ibuf)
    i, _, e := image.Decode(fszr)
    b := i.Bounds()

    isz, szok := rsimage.Getimgtype(b.Max.X, b.Max.Y)

    if szok == false {
        rscore.Sendstatus(rscore.C_WRSZ,
            "Image size to small or aspect ratio out of bounds", w)
        return settings
    }

    ni, rsz := rsimage.Scaleimage(i, isz)

    if rsz {
        b = ni.Bounds()
        os.RemoveAll(fnp)
        e = rsimage.Wrimagefile(ni, fnp)
        rscore.Cherr(e)
    }

    tags := rscore.Formattags(r.FormValue("tags"))

    nt, settings := rsdb.Tagstoindex(tags, settings)
    rscore.Sendtagstatus(nt, w)

    var suf []string
    suf = append(suf, rscore.IKEY[isz])
    stags := rscore.Addtagsuf(tags, suf)

    settings = rsdb.Addimgwtags(db, fn, b.Max.X, b.Max.Y, isz, stags, w, settings)
    rsdb.Wrsettings(db, settings)

    rscore.Sendstatus(rscore.C_OK, "", w)

    return settings
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    tags := rscore.Formattags(r.FormValue("tags"))

    tr := rscore.Textreq{
            Ttext: r.FormValue("ttext"),
            Btext: r.FormValue("btext"),
            Bpoint: r.FormValue("bpoint"),
            Tags: tags }

    ltxt, e := json.Marshal(tr)
    rscore.Addlog(rscore.L_REQ, ltxt, r)

    nt, settings := rsdb.Tagstoindex(tags, settings)
    rscore.Sendtagstatus(nt, w)

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

// Image wrapper for rsdb.Countobj()
func imgobjctr(db *bolt.DB, t string) int {

    var tl []string
    var ret int

    tl = append(tl, t)
    suf := rsimage.Getallsuf()
    stl := rscore.Addtagsuf(tl, suf)

    for _, st := range stl {
        ret += rsdb.Countobj(db, st, rscore.IBUC)
    }

    return ret
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
            ttag.IN = imgobjctr(db, t)
        }

        resp.Tags = append(resp.Tags, ttag)
    }

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Retrieves and validates user object
func userv(db *bolt.DB, w http.ResponseWriter, umax int, uname string,
    skey string, alevreq int) (bool, rscore.User) {

    if uname == "" {
        rscore.Sendstatus(rscore.C_NLOG,
            "User not logged in - no username provided", w)
    }

    if !rsdb.Isindb(db, []byte(uname), rscore.UBUC) {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return false, rscore.User{}
    }

    sok, u := valskey(db, uname, skey, umax)

    if !sok {
        rscore.Sendstatus(rscore.C_NLOG,
            "User not logged in - skey mismatch", w)
        return false, rscore.User{}
    }

    if u.Alev < alevreq {
        rscore.Sendstatus(rscore.C_ALEV,
            "User does not have sufficient access level", w)
        return false, u
    }

    return true, u
}

// Handles incoming requests for shutdowns
func shutdownhandler (w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    wipe := r.FormValue("wipe")
    uname := r.FormValue("user")
    skey := r.FormValue("skey")

    ok, _ := userv(db, w, settings.Umax, uname, skey, rscore.ALEV_ADMIN)
    if !ok { return }

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

// Handles incoming requests for user registrations
func reghandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    e := r.ParseForm()
    rscore.Cherr(e)

    u := rscore.User{}
    u.Name = r.FormValue("user")
    u.Email = r.FormValue("email") // TODO validate
    pass := r.FormValue("pass")

    if settings.Umax ==  0 {
        u.Alev = rscore.ALEV_ADMIN // Auto admin for first user to register

    } else {
        if rsdb.Isindb(db, []byte(u.Name), rscore.UBUC) {
            rscore.Sendstatus(rscore.C_UIDB, "Username already in db", w)
            return settings
        }

        u.Alev = rscore.ALEV_CONTRIB
    }

    if u.Name != rscore.Cleanstring(u.Name, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return settings
    }

    u.Skey = rscore.Randstr(rscore.SKEYLEN)
    u.Pass, e = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    rscore.Cherr(e)

    li := getloginobj(u)

    settings = rsdb.Wruindex(db, u.Name, settings)
    rsdb.Wruser(db, u)

    ml, e := json.Marshal(li)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, ml, r)

    enc := json.NewEncoder(w)
    enc.Encode(li)

    rsdb.Wrsettings(db, settings)
    return settings
}

// Creates login object
func getloginobj(u rscore.User) rscore.Login {

    ur := rscore.Login{}
    ur.Name = u.Name
    ur.Skey = u.Skey
    ur.Alev = u.Alev

    return ur
}

// Handles incoming login requests
func loginhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    e := r.ParseForm()
    rscore.Cherr(e)

    user := r.FormValue("user")
    pass := r.FormValue("pass")

    if settings.Umax < 1 {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return

    }
    // TODO refactor to new func
    if user != rscore.Cleanstring(user, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return
    }

    u := rsdb.Ruser(db, user)
    li := rscore.Login{}

    if rscore.Valuser(u, []byte(pass)) {
        u.Skey = rscore.Randstr(rscore.SKEYLEN)
        li = getloginobj(u)
        rsdb.Wruser(db, u)
    }

    ml, e := json.Marshal(li)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, ml, r)

    enc := json.NewEncoder(w)
    enc.Encode(li)
}

// Validates skey - returns true if user is logged in
func valskey(db *bolt.DB, uname string, skey string,
    umax int) (bool, rscore.User) {

    if umax < 1 { return false, rscore.User{} }

    u := rsdb.Ruser(db, uname)

    if skey == u.Skey { return true, u }
    return false, rscore.User{}
}

// Appends str to file at fname
func appendfile(fname string, str string) {

    f, e := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    rscore.Cherr(e)
    defer f.Close()

    _, e = f.WriteString(str)
    rscore.Cherr(e)
}

// Receives feedback data and saves to file
func feedbackhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    e := r.ParseForm()
    rscore.Cherr(e)

    uname := r.FormValue("user")
    skey := r.FormValue("skey")
    str := r.FormValue("fb")

    ok, u := userv(db, w, settings.Umax, uname, skey, rscore.ALEV_CONTRIB)
    if !ok { return }

    d := fmt.Sprintf("%s (%s): %s\n", u.Name, u.Email, str)
    appendfile(rscore.FBFILE, d)
    rscore.Sendstatus(rscore.C_OK, "", w)
}

func main() {

    pptr := flag.Int("p", rscore.DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", rscore.VERBDEF, "verbose mode")
    xptr := flag.Bool("x", rscore.VOLATILEDEF, "volatile mode")
    flag.Parse()

    db := rsdb.Open(*dbptr)
    defer db.Close()

    settings := rsdb.Rsettings(db)
    settings.Verb = *vptr
    settings = rscore.Rsinit(settings)

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    if *xptr {
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

    // User login
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        loginhandler(w, r, db, settings)
    })

    // Feedback
    http.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
        feedbackhandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e := http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}