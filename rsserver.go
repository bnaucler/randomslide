package main

import (
    "os"
    "fmt"
    "flag"
    "bytes"
    "image"
    "regexp"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "encoding/json"
    "mime/multipart"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
    "github.com/bnaucler/randomslide/lib/rsimage"
    "github.com/bnaucler/randomslide/lib/rsuser"
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

    var ctr int

    for i := 0 ; i < n ; i ++ {
        ctr = 0
        for len(bp) == 0 || len(bp) > rscore.BPOINTMAX {
            bp = rsdb.Getrndtxt(db, settings.Tmax, tags, rscore.TBUC)
            if ctr > 50 { return bps } // TODO
            ctr++
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

// Generate numbers of slide type 2 (big number)
func numgen() string {

    p := byte(' ')
    s := byte(' ')

    b := rand.Intn(rscore.RNUMBMAX)
    n := rndexp(b)

    plen := len(rscore.NUMPREF)
    if coin() { p = rscore.NUMPREF[rand.Intn(plen)] }

    slen := len(rscore.NUMSUFF)
    if coin() { s = rscore.NUMSUFF[rand.Intn(slen)] }

    if n == 0 && p == '-' { p = ' ' }       // '-0' looks a bit stupid
    if p == '$' && s == '%' { s = ' ' }     // and so does '$5%'

    ret := fmt.Sprintf("%s%d%s", string(p), n, string(s))

    return strings.TrimSpace(ret)
}

// Generate data points for slide type 7 (graph)
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
    rsdb.Wrdeck(db, deck)

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
        deck = rsdb.Rdeck(db, deck, req.Id, settings)

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

    c := getcall(r)

    var e error
    n := 0

    if len(c.Amount) > 0 {
        n, e = strconv.Atoi(c.Amount)
        rscore.Cherr(e)
    }

    id, isidr := isidreq(r)
    tags := rscore.Formattags(c.Tags)
    if len(tags) < 1 || tags[0] == "" { tags = settings.Taglist }

    req := rscore.Deckreq{
            Id: id,
            Isidreq: isidr,
            N: n,
            Lang: c.Lang,
            Tags: tags }

    mreq, e := json.Marshal(req)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, mreq, settings.Llev, r)

    deck, settings := getdeck(req, db, settings)

    mdeck, e := json.Marshal(deck)
    rscore.Addlog(rscore.L_RESP, mdeck, settings.Llev, r)

    enc := json.NewEncoder(w)
    enc.Encode(deck)

    return settings
}

// addlog() wrapper for file requests
func logfreq(hlr *multipart.FileHeader, mt string, llev int, r *http.Request) {

    lmsg := fmt.Sprintf("File: %+v(%s) - %+v",
        hlr.Filename, rscore.Prettyfsize(hlr.Size), mt)
    rscore.Addlog(rscore.L_REQ, []byte(lmsg), llev, r)
}

// Wrapper for file mime type check
func chkimgmime(hlr *multipart.FileHeader, settings rscore.Settings,
    w http.ResponseWriter) (bool, string) {

    mt := hlr.Header["Content-Type"][0]

    if rscore.Findstrinslice(mt, rscore.IMGMIME) == false {
        rscore.Sendstatus(rscore.C_WRFF,
            "Unknown image format - file not uploaded", w)
        return false, mt
    }

    return true, mt
}

// Retrieves and checks tags from request
func gettags(r *http.Request, w http.ResponseWriter) (bool, []string) {

    tags := rscore.Formattags(r.FormValue("tags"))

    if len(tags) < 1  || tags[0] == "" {
        rscore.Sendstatus(rscore.C_NTAG,
            "No tags provided - cannot add data", w)
        return false, tags
    }

    return true, tags
}

// Handles incoming requests to add images
func imgreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    r.ParseMultipartForm(10 << 20)
    sf, hlr, e := r.FormFile("file")
    e = rscore.Cherr(e)
    if e != nil { return settings }
    defer sf.Close()

    uname := r.FormValue("user")
    skey := r.FormValue("skey")

    ok, _ := rsuser.Userv(db, w, settings.Umax, uname, skey, rscore.ALEV_CONTRIB)
    if !ok { return settings }

    ok, mt := chkimgmime(hlr, settings, w)
    if !ok { return settings }
    logfreq(hlr, mt, settings.Llev, r)

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

    _, b = rsimage.Scaleimage(i, isz, fnp)

    ok, tags := gettags(r, w)
    if !ok { return settings }

    nt, settings := rsdb.Tagstoindex(tags, settings)
    rscore.Sendtagstatus(nt, w)

    var suf []string
    suf = append(suf, rscore.SUFINDEX[isz])
    stags := rscore.Addtagsuf(tags, suf)

    settings = rsdb.Addimgwtags(db, fn, b.Max.X, b.Max.Y, isz, stags, w, settings)
    rsdb.Wrsettings(db, settings)

    rscore.Sendstatus(rscore.C_OK, "", w)

    return settings
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    c := getcall(r)

    ok, _ := rsuser.Userv(db, w, settings.Umax, c.User, c.Skey, rscore.ALEV_CONTRIB)
    if !ok { return settings }

    ok, tags := gettags(r, w)
    if !ok { return settings }

    tr := rscore.Textreq{
            Ttext: c.Ttext,
            Btext: c.Btext,
            Bpoint: c.Bpoint,
            Tags: tags }

    ltxt, e := json.Marshal(tr)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, ltxt, settings.Llev, r)

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

// Handles incoming requests for user index
func userreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    resp := rsdb.Ruindex(db, settings)

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, settings.Llev, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
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
            ttag.IN = rsdb.Imgobjctr(db, t)
        }

        resp.Tags = append(resp.Tags, ttag)
    }

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, settings.Llev, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for shutdowns
func shutdownhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    c := getcall(r)

    ok, _ := rsuser.Userv(db, w, settings.Umax, c.User, c.Skey, rscore.ALEV_ADMIN)
    if !ok { return }

    rscore.Addlog(rscore.L_SHUTDOWN, []byte(""), settings.Llev, r)
    rscore.Sendstatus(rscore.C_OK, "", w)

    rsdb.Wrsettings(db, settings)

    if c.Wipe == "yes" {
        db.Close()
        os.Remove(rscore.DBNAME)
        rscore.Rmall(rscore.IMGDIR)
    }

    rscore.Shutdown(settings)
}

// Wrapper for basic checks of valid operations
func getop(rop string, w http.ResponseWriter) (bool, int) {

    if len(rop) < 1 {
        rscore.Sendstatus(rscore.C_NSOP, "No such operation", w)
        return false, 0
    }

    op, e := strconv.Atoi(rop)

    if e != nil {
        rscore.Sendstatus(rscore.C_NSOP, "No such operation", w)
        return false, 0
    }

    return true, op
}

// Parses API call and returns object
func getcall(r *http.Request) rscore.Apicall {

    e := r.ParseForm()
    rscore.Cherr(e)

    ret := rscore.Apicall{
        User:       r.FormValue("user"),
        Pass:       r.FormValue("pass"),
        Email:      r.FormValue("email"),
        Skey:       r.FormValue("skey"),
        Tuser:      r.FormValue("tuser"),
        Tags:       r.FormValue("tags"),
        Lang:       r.FormValue("lang"),
        Id:         r.FormValue("id"),
        Amount:     r.FormValue("amount"),
        Ttext:      r.FormValue("ttext"),
        Btext:      r.FormValue("btext"),
        Fb:         r.FormValue("fb"),
        Bpoint:     r.FormValue("bpoint"),
        Rop:        r.FormValue("op"),
        Wipe:       r.FormValue("wipe"),
    }

    return ret
}

// Changes user settings
func cuhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    c := getcall(r)
    ok, op := getop(c.Rop, w)
    tu := rsdb.Ruser(db, c.Tuser)

    if tu.Name != c.Tuser {
        rscore.Sendstatus(rscore.C_NOSU, "No such target user", w)
        return
    }

    switch {
    case op == rscore.CU_MKADM || op == rscore.CU_RMADM:
        ok, tu = rsuser.Chadminstatus(db, op, settings.Umax, c.User, c.Skey, tu, w)

    case op == rscore.CU_CPASS:
        ok, tu = rsuser.Chpass(db, settings, c.User, c.Skey, c.Pass, tu, w)

    case op == rscore.CU_RMUSR:
        ok, settings = rsuser.Rmuser(db, settings, c.User, c.Skey, tu, w)

    default:
        rscore.Sendstatus(rscore.C_NSOP, "No such operation", w)
        return
    }

    if ok && op == rscore.CU_CPASS {
        rsuser.Senduser(tu, r, w, settings)
        rsdb.Wruser(db, tu)

    } else if ok {
        rscore.Sendstatus(rscore.C_OK, "", w)
        rsdb.Wruser(db, tu)

    } else if !ok && op == rscore.CU_CPASS {
        rscore.Sendstatus(rscore.C_USPW, "Unsafe password", w)
    }
}

// Returns true if email address looks valid
func valemail(addr string) bool {

    rx := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")

    if rx.MatchString(addr) { return true }
    return false
}

// Handles incoming requests for user registrations
func reghandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) rscore.Settings {

    c := getcall(r)
    u := rscore.User{}

    if !valemail(c.Email) {
        rscore.Sendstatus(rscore.C_IEMA, "Invalid email address", w)
        return settings
    }

    u.Name = c.User
    u.Email = c.Email

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

    ok, u := rsuser.Setpass(u, c.Pass)
    if !ok {
        rscore.Sendstatus(rscore.C_USPW, "Unsafe password", w)
        return settings
    }

    settings = rsdb.Addutoindex(db, u.Name, settings)
    rsdb.Wruser(db, u)
    rsuser.Senduser(u, r, w, settings)

    rsdb.Wrsettings(db, settings)
    return settings
}

// Handles incoming login requests
func loginhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    c := getcall(r)

    if settings.Umax < 1 {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return
    }

    if c.User != rscore.Cleanstring(c.User, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return
    }

    uindex := rsdb.Ruindex(db, settings)
    if !rscore.Findstrinslice(c.User, uindex.Names) {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return
    }

    u := rsdb.Ruser(db, c.User)
    li := rscore.Login{}

    if rscore.Valuser(u, []byte(c.Pass)) {
        u.Skey = rscore.Randstr(rscore.SKEYLEN)
        li = rsuser.Getloginobj(u)
        rsdb.Wruser(db, u)
    }

    ml, e := json.Marshal(li)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, ml, settings.Llev, r)

    enc := json.NewEncoder(w)
    enc.Encode(li)
}

// Receives feedback data and saves to file
func feedbackhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB,
    settings rscore.Settings) {

    c := getcall(r)

    ok, u := rsuser.Userv(db, w, settings.Umax, c.User, c.Skey, rscore.ALEV_CONTRIB)
    if !ok { return }

    d := fmt.Sprintf("%s (%s): %s\n", u.Name, u.Email, c.Fb)
    rscore.Appendfile(rscore.FBFILE, d)
    rscore.Sendstatus(rscore.C_OK, "", w)
}

func main() {

    pptr := flag.Int("p", rscore.DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", rscore.VERBDEF, "increase log level")
    flag.Parse()

    db := rsdb.Open(*dbptr)
    defer db.Close()

    settings := rsdb.Rsettings(db)
    if *vptr { settings.Llev = 1 }
    settings = rscore.Rsinit(settings)

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Requests to shut down server
    http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
        shutdownhandler(w, r, db, settings)
    })

    // Slide deck requests
    http.HandleFunc("/getdeck", func(w http.ResponseWriter, r *http.Request) {
        settings = deckreqhandler(w, r, db, settings)
    })

    // Tags requests
    http.HandleFunc("/gettags", func(w http.ResponseWriter, r *http.Request) {
        tagreqhandler(w, r, db, settings)
    })

    // Tags requests
    http.HandleFunc("/getusers", func(w http.ResponseWriter, r *http.Request) {
        userreqhandler(w, r, db, settings)
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

    // Change user settings
    http.HandleFunc("/chuser", func(w http.ResponseWriter, r *http.Request) {
        cuhandler(w, r, db, settings)
    })

    // Feedback
    http.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
        feedbackhandler(w, r, db, settings)
    })

    lport := fmt.Sprintf(":%d", *pptr)
    e := http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}
