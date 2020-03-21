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
    "net/smtp"
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

// Sends simple plaintext email
func sendmail(addr string, subj string, body string) {

    var to []string
    var s = rscore.Set.Smtp

    to = append(to, addr)

    a := smtp.PlainAuth("", s.User, s.Pass, s.Server)

    m := fmt.Sprintf("To: %s\r\n" +
                      "Subject: %s\r\n\r\n" +
                      "%s\r\n", addr, subj, body)

    swp := fmt.Sprintf("%s:%d", s.Server, s.Port)

    e := smtp.SendMail(swp, a, s.From, to, []byte(m))
    rscore.Cherr(e)
}

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
func bpgen(db *bolt.DB, tags []string) []string {

    var bps []string
    var bp string

    rnd := rscore.BPMAX - rscore.BPMIN
    n := rand.Intn(rnd)
    n += rscore.BPMIN

    var ctr int

    for i := 0 ; i < n ; i++ {
        ctr = 0
        for len(bp) == 0 || len(bp) > rscore.BPOINTMAX {
            bp = rsdb.Getrndtxt(db, rscore.Set.Tmax, tags, rscore.TBUC)
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
func getslide(db *bolt.DB, st rscore.Slidetype, req rscore.Deckreq) rscore.Slide {

    slide := rscore.Slide{Type: st.Type}

    if st.TT { slide.Title = rsdb.Getrndtxt(db, rscore.Set.Tmax, req.Tags, rscore.TBUC) }
    if st.BT { slide.Btext = rsdb.Getrndtxt(db, rscore.Set.Bmax, req.Tags, rscore.BBUC) }
    if len(st.IMG) > 0 {
        suf := rsimage.Mkimgsuflist(st.Type)
        stags := rscore.Addtagsuf(req.Tags, suf)
        slide.Img = rsdb.Getrndimg(db, rscore.Set.Imax, stags, rscore.IBUC)
    }

    switch st.Type {

    case 2:
        slide.Title = numgen()

    case 3:
        slide.Bpts = bpgen(db, req.Tags)

    case 7:
        slide.Dpts = dpgen()
        slide.Ctype = rand.Intn(3)
    }

    return slide
}

// Returns a new slide deck according to request
func mkdeck(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq) rscore.Deck {

    var st rscore.Slidetype
    sprob := rscore.SPROB

    for i := 0; i < req.N; i++ {
        st, sprob = setslidetype(i, sprob)
        slide := getslide(db, st, req)
        deck.Slides = append(deck.Slides, slide)
    }

    deck.Id = rscore.Set.Dmax
    rsdb.Wrdeck(db, deck)

    // Mutex test
    rscore.Smut.Lock()
    rscore.Set.Dmax++
    rsdb.Wrsettings(db, rscore.Set)
    rscore.Smut.Unlock()

    return deck
}

// Sets basic params & determines if new deck should be built
func getdeck(req rscore.Deckreq, db *bolt.DB) rscore.Deck {

    deck := rscore.Deck{
            Id: req.Id,
            N: req.N,
            Lang: req.Lang }

    if req.Isidreq {
        deck = rsdb.Rdeck(db, deck, req.Id, rscore.Set)

    } else {
        deck = mkdeck(db, deck, req)

    }

    return deck
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
func deckreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)

    var e error
    n := 0

    if len(c.Amount) > 0 {
        n, e = strconv.Atoi(c.Amount)
        rscore.Cherr(e)
    }

    id, isidr := isidreq(r)
    tags := rscore.Formattags(c.Tags)
    if len(tags) < 1 || tags[0] == "" { tags = rscore.Set.Taglist }

    req := rscore.Deckreq{
            Id: id,
            Isidreq: isidr,
            N: n,
            Lang: c.Lang,
            Tags: tags }

    mreq, e := json.Marshal(req)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, mreq, rscore.Set.Llev, rscore.User{}, r)

    deck := getdeck(req, db)

    mdeck, e := json.Marshal(deck)
    rscore.Addlog(rscore.L_RESP, mdeck, rscore.Set.Llev, rscore.User{}, r)

    enc := json.NewEncoder(w)
    enc.Encode(deck)
}

// addlog() wrapper for file requests
func logfreq(hlr *multipart.FileHeader, mt string, llev int, r *http.Request) {

    lmsg := fmt.Sprintf("File: %+v(%s) - %+v",
        hlr.Filename, rscore.Prettyfsize(hlr.Size), mt)
    rscore.Addlog(rscore.L_REQ, []byte(lmsg), llev, rscore.User{}, r) // TODO
}

// Wrapper for file mime type check
func chkimgmime(hlr *multipart.FileHeader, w http.ResponseWriter) (bool, string) {

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
func imgreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    r.ParseMultipartForm(10 << 20)
    sf, hlr, e := r.FormFile("file")
    e = rscore.Cherr(e)
    if e != nil { return }
    defer sf.Close()

    c := getcall(r)

    ok, _ := rsuser.Userv(db, w, rscore.Set.Umax, c, rscore.ALEV_CONTRIB)
    if !ok { return }

    ok, mt := chkimgmime(hlr, w)
    if !ok { return }
    logfreq(hlr, mt, rscore.Set.Llev, r)

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
        return
    }

    _, b = rsimage.Scaleimage(i, isz, fnp)

    ok, tags := gettags(r, w)
    if !ok { return }

    var nt int
    nt, rscore.Set = rsdb.Tagstoindex(tags, rscore.Set) // TODO mutex
    rscore.Sendtagstatus(nt, w)

    var suf []string
    suf = append(suf, rscore.SUFINDEX[isz])
    stags := rscore.Addtagsuf(tags, suf)

    // Mutex test
    rscore.Smut.Lock()
    rscore.Set = rsdb.Addimgwtags(db, fn, b.Max.X, b.Max.Y, isz, stags, w, rscore.Set)
    rsdb.Updatetindex(db)
    rsdb.Wrsettings(db, rscore.Set)
    rscore.Smut.Unlock()

    rscore.Sendstatus(rscore.C_OK, "", w)
}

// Handles incoming requests to add text
func textreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)

    ok, u := rsuser.Userv(db, w, rscore.Set.Umax, c, rscore.ALEV_CONTRIB)
    if !ok { return }

    ok, tags := gettags(r, w)
    if !ok { return }

    tr := rscore.Textreq{
            Ttext: c.Ttext,
            Btext: c.Btext,
            Bpoint: c.Bpoint,
            Tags: tags }

    ltxt, e := json.Marshal(tr)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_REQ, ltxt, rscore.Set.Llev, u, r)

    var nt int
    nt, rscore.Set = rsdb.Tagstoindex(tags, rscore.Set) // TODO mutex
    rscore.Sendtagstatus(nt, w)

    rscore.Smut.Lock() // Mutex test
    if len(tr.Ttext) > 1 && len(tr.Ttext) < rscore.TTEXTMAX {
        rsdb.Addtextwtags(tr.Ttext, tags, db, rscore.Set.Tmax, rscore.TBUC)
        rscore.Set.Tmax++
    }

    if len(tr.Btext) > 1 && len(tr.Btext) < rscore.BTEXTMAX {
        rsdb.Addtextwtags(tr.Btext, tags, db, rscore.Set.Bmax, rscore.BBUC)
        rscore.Set.Bmax++
    }

    rsdb.Updatetindex(db)
    rsdb.Wrsettings(db, rscore.Set)
    rscore.Smut.Unlock()
}

// Handles incoming requests for user index
func userreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    resp := rsdb.Ruindex(db, rscore.Set)

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, rscore.Set.Llev, rscore.User{}, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for tag index
func tagreqhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    resp := rsdb.Rtindex(db)

    mresp, e := json.Marshal(resp)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, mresp, rscore.Set.Llev, rscore.User{}, r)

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Handles incoming requests for shutdowns
func shutdownhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)

    ok, _ := rsuser.Userv(db, w, rscore.Set.Umax, c, rscore.ALEV_ADMIN)
    if !ok { return }

    rscore.Addlog(rscore.L_SHUTDOWN, []byte(""), rscore.Set.Llev, rscore.User{}, r)
    rscore.Sendstatus(rscore.C_OK, "", w)

    rsdb.Wrsettings(db, rscore.Set)

    if c.Wipe == "yes" {
        db.Close()
        os.Remove(rscore.DBNAME)
        rscore.Rmall(rscore.IMGDIR)
    }

    rscore.Shutdown(rscore.Set)
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

// Sets new password and sends by email
func pwdreset(db *bolt.DB, c rscore.Apicall,
    w http.ResponseWriter) (bool, rscore.User) {

    u := rsdb.Ruser(db, c.User)

    if u.Email != c.Email {
        rscore.Sendstatus(rscore.C_WEMA, "Incorrect email address for user", w)
        return false, u
    }

    np := rscore.Randstr(rscore.RPWDLEN)
    ok, u := rsuser.Setpass(u, np)
    if !ok { return false, u }

    msg := fmt.Sprintf("Your new randomslide password: %s", np)
    go sendmail(u.Email, "randomslide password reset", msg)

    return ok, u
}

// Changes user settings
func cuhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)
    ok, op := getop(c.Rop, w)
    tu := rsdb.Ruser(db, c.Tuser)

    if tu.Name != c.Tuser || c.Tuser == "" {
        rscore.Sendstatus(rscore.C_NOSU, "No such target user", w)
        return
    }

    switch {
    case op == rscore.CU_MKADM || op == rscore.CU_RMADM:
        ok, tu = rsuser.Chadminstatus(db, op, rscore.Set.Umax, c, tu, w)

    case op == rscore.CU_CPASS:
        ok, tu = rsuser.Chpass(db, rscore.Set, c, tu, w)

    case op == rscore.CU_RMUSR:
        ok, rscore.Set = rsuser.Rmuser(db, rscore.Set, c, tu, w)

    case op == rscore.CU_PWDRS:
        ok, tu = pwdreset(db, c, w)

    default:
        rscore.Sendstatus(rscore.C_NSOP, "No such operation", w)
        return
    }

    if ok && op == rscore.CU_CPASS {
        rsuser.Senduser(tu, r, w, rscore.Set)
        rsdb.Wruser(db, tu)

    } else if ok {
        rscore.Sendstatus(rscore.C_OK, "", w)
        rsdb.Wruser(db, tu)
    }
}

// Returns true if email address looks valid
func valemail(addr string) bool {

    rx := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")

    if rx.MatchString(addr) { return true }
    return false
}

// Handles incoming requests for user registrations
func reghandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)
    u := rscore.User{}

    if !valemail(c.Email) {
        rscore.Sendstatus(rscore.C_IEMA, "Invalid email address", w)
        return
    }

    u.Name = c.User
    u.Email = c.Email

    if rscore.Set.Umax ==  0 {
        u.Alev = rscore.ALEV_ADMIN // Auto admin for first user to register

    } else {
        if rsdb.Isindb(db, []byte(u.Name), rscore.UBUC) {
            rscore.Sendstatus(rscore.C_UIDB, "Username already in db", w)
            return
        }

        u.Alev = rscore.ALEV_CONTRIB
    }

    if u.Name != rscore.Cleanstring(u.Name, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return
    }

    ok, u := rsuser.Setpass(u, c.Pass)
    if !ok {
        rscore.Sendstatus(rscore.C_USPW, "Unsafe password", w)
        return
    }

    rscore.Set = rsdb.Addutoindex(db, u.Name, rscore.Set)
    rsdb.Wruser(db, u)
    rsuser.Senduser(u, r, w, rscore.Set)

    rsdb.Wrsettings(db, rscore.Set)
    return
}

// Handles incoming login requests
func loginhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)

    if rscore.Set.Umax < 1 {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return
    }

    if c.User != rscore.Cleanstring(c.User, rscore.RXUSER) {
            rscore.Sendstatus(rscore.C_UICH,
                "Username includes illegal characters", w)
            return
    }

    uindex := rsdb.Ruindex(db, rscore.Set)
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
    rscore.Addlog(rscore.L_RESP, ml, rscore.Set.Llev, u, r)

    enc := json.NewEncoder(w)
    enc.Encode(li)
}

// Receives feedback data and saves to file
func feedbackhandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {

    c := getcall(r)

    ok, u := rsuser.Userv(db, w, rscore.Set.Umax, c, rscore.ALEV_CONTRIB)
    if !ok { return }

    d := fmt.Sprintf("%s (%s): %s\n", u.Name, u.Email, c.Fb)
    rscore.Appendfile(rscore.FBFILE, d)
    rscore.Sendstatus(rscore.C_OK, "", w)
}

// Launches mapped handler functions
func starthlr(url string, fn rscore.Hfn, db *bolt.DB) {

    http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
        fn(w, r, db)
    })
}

func main() {

    pptr := flag.Int("p", rscore.DEFAULTPORT, "port number to listen")
    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", rscore.VERBDEF, "increase log level")
    flag.Parse()

    db := rsdb.Open(*dbptr)
    defer db.Close()

    rscore.Set = rsdb.Rsettings(db)
    if *vptr { rscore.Set.Llev = 1 }
    rscore.Set = rscore.Rsinit(rscore.Set)

    // Assuming first start, ensuring creation of SBUC
    if rscore.Set.Umax == 0 { rsdb.Wrsettings(db, rscore.Set) }

    // Static content
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Map endpoints to handlers
    var hlrs = map[string]rscore.Hfn {
        "/restart":     shutdownhandler,
        "/getdeck":     deckreqhandler,
        "/gettags":     tagreqhandler,
        "/getusers":    userreqhandler,
        "/addtext":     textreqhandler,
        "/addimg":      imgreqhandler,
        "/register":    reghandler,
        "/login":       loginhandler,
        "/chuser":      cuhandler,
        "/feedback":    feedbackhandler,
    }

    // Launch all handlers
    for url, fn := range hlrs { starthlr(url, fn, db) }

    lport := fmt.Sprintf(":%d", *pptr)
    e := http.ListenAndServe(lport, nil)
    rscore.Cherr(e)
}
