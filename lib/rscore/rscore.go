package rscore

/*

    Package of core operations and data types
    Used by randomslide

*/

import (
    "os"
    "io"
    "fmt"
    "log"
    "sync"
    "time"
    "regexp"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "os/signal"
    "math/rand"
    "encoding/json"
    "path/filepath"

    "github.com/boltdb/bolt"
    "golang.org/x/crypto/bcrypt"
)

const DEFAULTPORT = 6291        // Port can also be supplied with -p flag

const DBNAME = "./data/rs.db"   // Database location (change at own peril)
const LOGPATH = "./static/log/" // Logs should be accessible from frontend
const IMGDIR = "./static/img/"  // Image directory
const PIDFILEPATH = "./data/"   // Base directory for storage of PID file
const FBFILE = "./data/fb.txt"  // Storage of feedback data
const THDIR = "./static/css/themes/" // CSS theme file directory

var Set Settings                // Global rsserver settings
var Smut sync.Mutex             // Settings mutex

const VERBDEF = false           // Verbose mode defaults to false

const TTEXTMAX = 35             // Max length for title text
const BTEXTMAX = 80             // Max length for body text
const BPOINTMAX = 20            // Max length for bullet point
const RNUMBMAX = 30             // Random number base max
const RNUMEMAX = 3              // Random number exponent max

const UNMINLEN = 2              // Minimum accepted lendth of user names
const PWMINLEN = 6              // Minimum accepted lendth of password
const RPWDLEN = 15              // Length of random generated passwords

const SKEYLEN = 40              // # of characters in a session key
const RFNLEN = 20               // Length of random file names (w/o .ext)

var NUMPREF = []byte("$+-")     // Potential number prefixes for slide type 2
var NUMSUFF = []byte("%!?")     // Potential number suffixes for slide type 2

const STYPES = 8                // Number of slide types available
const BPMIN = 3                 // Min number of bullet points for lists
const BPMAX = 8                 // Max number of bullet points for lists

var RXUSER = "[^a-z0-9]+"       // Regex for allowed user names
var RXTAGS = "[^a-zåäöüæø]+"    // Regex for allowed tags

// Logging codes parsed by Addlog()
const L_REQ = 0                 // Request log
const L_RESP = 1                // Response log
const L_SHUTDOWN = 2            // Server shutdown request log

// Status response codes sent to client
const C_OK = 0                  // OK
const C_WRFF = 1                // Incorrect file format
const C_WRSZ = 2                // Not able to classify image size
const C_UIDB = 3                // User already exists in database
const C_UICH = 4                // Username includes illegal characters
const C_NOSU = 5                // No such user
const C_NLOG = 6                // User not logged in
const C_ALEV = 7                // User does not have sufficient access level
const C_NTAG = 8                // No tags provided
const C_NSOP = 9                // No such operation
const C_UNKN = 10               // Unknown error
const C_USPW = 11               // Unsafe password
const C_IEMA = 12               // Invalid email address
const C_WEMA = 13               // Incorrect email address
const C_NOBJ = 14               // No such database object
const C_UNTS = 15               // User name too short

// User level definitions
const ALEV_USER = 0             // Regular user
const ALEV_CONTRIB = 1          // Contributor
const ALEV_ADMIN = 2            // BOFH

// Operation codes for user handling
const CU_MKADM = 0              // Makes specified user admin
const CU_RMADM = 1              // Removes admin status from user
const CU_CPASS = 2              // Password change request
const CU_RMUSR = 3              // Removes user account
const CU_PWDRS = 4              // Requests password reset email

// Probability chart for slide occurance. Higher number = higher probability.
var SPROB = []int{2, 6, 3, 4, 9, 6, 5, 4}

// Image type classifications
const IMG_XL = 0                // X Large
const IMG_LS = 1                // Landscape
const IMG_BO = 2                // Box-shaped
const IMG_PO = 3                // Portrait

// Min bounds for image sizes (w, h)
var IMGMIN = [][]uint{
    {1920, 1080},               // 0: X Large
    {640, 360},                 // 1: Landscape
    {360, 360},                 // 2: Box-shaped
    {360, 480},                 // 3: Portrait
}

// Max bounds for image sizes (w, h)
var IMGMAX = [][]uint{
    {1920, 1080},               // 0: X Large
    {1920, 1080},               // 1: Landscape
    {1080, 1080},               // 2: Box-shaped
    {1080, 1920},               // 3: Portrait
}

var DBUC = []byte("dbuc")       // Deck bucket
var TBUC = []byte("tbuc")       // Title text bucket
var BBUC = []byte("bbuc")       // Body text bucket
var IBUC = []byte("ibuc")       // Image bucket
var SBUC = []byte("sbuc")       // Settings bucket
var UBUC = []byte("ubuc")       // User bucket

var INDEX = []byte(".index")    // Untouchable database index position
var TINDEX = []byte(".tindex")  // Untouchable database tag index position

// Image size per slide type reference chart TODO map with TT & BT bools
var ISZINDEX = [][]int {
    {0, 1},                     // Big title
    {0},                        // Full screen image
    {0},                        // Big number
    {3},                        // Bullet point list
    {1, 2, 3},                  // Title, img & body
    {},                         // Inspirational quote
    {1, 2, 3},                  // Image with text
    {},                         // Graph
}

// Data object including references to all IKEYs
var ALLSUF = []int{0, 1, 2, 3}

// Index keys to be used for image size indexes
var SUFINDEX = []string {
    ".s0",
    ".s1",
    ".s2",
    ".s3",
}

// For verification of image mime types
var IMGMIME = []string{
    "image/jpeg",
    "image/png",
    "image/gif",
}

// Launcher of handle functions
type Hfn func(http.ResponseWriter, *http.Request, *bolt.DB)

type Settings struct {
    Llev int                    // Log level
    Dmax int                    // Max id of decks
    Tmax int                    // Max id of title objects
    Bmax int                    // Max id of body objects
    Imax int                    // Max id of image objects
    Umax int                    // Max user ID in database
    Pidfile string              // Location of pidfile
    Taglist []string            // List of all existing tags
    Smtp Smtp                   // Settings for sending email
}

type User struct {
    Name string                 // Username
    Pass []byte                 // Password hash
    Email string                // Email address
    Skey string                 // Session key
    Alev int                    // Access level
}

type Login struct {
    Name string                 // User name
    Skey string                 // Session key
    Alev int                    // Access level
}

type Apicall struct {
    User string                 // User name
    Pass string                 // Password
    Email string                // Email address
    Skey string                 // Session key
    Tuser string                // Target user
    Tags string                 // Tags
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Id string                   // Deck ID request
    Slide string                // Slide number for reference
    Amount string               // Amount of slides requested
    Ttext string                // Title test
    Btext string                // Body test
    Msg string                  // Feedback text
    Bpoint string               // Bullet point
    Rop string                  // Raw (unprocessed) operation
    Wipe string                 // Request to wipe database
    Type string                 // Object type reference
}

type Smtp struct {
    Admin string                // Admin email address
    Server string               // Which server to connect to
    User string                 // User account for logging in
    Pass string                 // Password for user auth
    Port int                    // Server port
}

type Deckreq struct {
    Id int                      // Deck ID for db
    N int                       // Number of slides to generate
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Tags []string               // Slice of tags on which to base search
    Isidreq bool                // true if request has specified ID
}

type Textreq struct {
    Ttext string                // Title text object to add to db
    Btext string                // Body text object to add to db
    Bpoint string               // Bullet points for making lists
    Tags []string               // Tags for indexing
}

type Thindex struct {
    Themes []string             // Theme filenames
}

type Iindex struct {
    Ids    []int                // All IDs associated with tag
}

type Uindex struct {
    Names []string              // All user names in database
}

type Rtag struct {
    Name string                 // Tag name
    TN int                      // Number of title objects in db
    BN int                      // Number of body objects in db
    IN int                      // Number of image objects in db
}

type Tagresp struct {
    Tags []Rtag                 // Array of tags for indexing
}

type Textobj struct {
    Id int                      // Index number
    Text string                 // The text itself
    Tags []string               // All tags where object exists (for associative decks)
    Contr string                // Text object contributor
}

type Imgobj struct {
    Id int                      // Index number
    Fname string                // File name
    Tags []string               // All tags where object exists (for associative decks)
    Size int                    // S (0), M (1), L (2) or XL (3)
    Contr string                // Image object contributor
    H int                       // Image height
    W int                       // Image width
}

type Deck struct {
    Id int                      // Deck ID for db
    N int                       // Total number of slides in deck
    Lang string                 // Languge code, 'en', 'de', 'se', etc
    Slides []Slide              // Slice of Slide objects
}

type Slidetype struct {
    Type int                    // The type ID number
    TT bool                     // Includes title text
    BT bool                     // Includes body text
    IMG []int                   // Image size preferences
}

type Slide struct {
    Type int                    // See CONTRIBUTING.md for type chart
    Title string                // Slide title
    Btext string                // Body text
    Bpts []string               // Bullet points
    Dpts []int                  // Graph data points
    Ctype int                   // Chart type
    Img Imgobj                  // Image object
}

type Statusresp struct {
    Code int                    // Error code to be parsed in frontend
    Text string                 // Additional related data
}

// Log all errors to file
func Cherr(e error) error {
    if e != nil { log.Println(e) }
    return e
}

// Enables clean shutdown. Needs delay for caller to send response
func Shutdown(settings Settings) {

    go func() {
        time.Sleep(1 * time.Second)
        os.Remove(settings.Pidfile)
        os.Exit(0)
    }()
}

// Setting up signal handler
func Sighandler(settings Settings) {

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func(){
        for sig := range sigc {
            fmt.Printf("Caught %+v - cleaning up.\n", sig)
            Shutdown(settings)
        }
    }()
}

// Creating PID file
func Mkpidfile(settings Settings, prgname string, pid int) Settings {

    settings.Pidfile = fmt.Sprintf("%s%s.pid", PIDFILEPATH, prgname)
    e := ioutil.WriteFile(settings.Pidfile, []byte(strconv.Itoa(pid)), 0644)
    Cherr(e)

    return settings
}

// Initialize logger
func Initlog(prgname string) {

    logfile := fmt.Sprintf("%s/%s.log", LOGPATH, prgname)

    f, e := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    Cherr(e)

    log.SetOutput(f)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Creates pid file, signal handler and starts logging
func Rsinit(settings Settings) Settings {

    prgname := filepath.Base(os.Args[0])
    pid := os.Getpid()

    settings = Mkpidfile(settings, prgname, pid)
    Sighandler(settings)
    Initlog(prgname)

    log.Printf("%s started with PID: %d\n", prgname, pid)

    return settings
}

// Returns random sting with length ln
func Randstr(ln int) (string){

    const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
    var cslen = len(charset)

    b := make([]byte, ln)
    for i := range b { b[i] = charset[rand.Intn(cslen)] }

    return string(b)
}

// Removes whitespace and special characters from string
func Cleanstring(src string, pat string) string {

    rx, e := regexp.Compile(pat)
    Cherr(e)

    dst := rx.ReplaceAllString(src, "")

    return dst
}

// Returns true if string is present in list
func Findstrinslice(v string, list []string) bool {

    for _, t := range list {
        if v == t { return true }
    }

    return false
}

// Compares two byte slices, returs true if identical
func Identicalbs(s1 []byte, s2 []byte) bool {

    if len(s1) != len(s2) { return false }

    for i := 0; i < len(s1); i++ {
        if s1[i] != s2[i] { return false }
    }

    return true
}

// Removes int from slice if present returns 1 if action taken
func Rmifrslice(n int, s []int) ([]int, bool) {

    var ud bool
    var ret []int

    for _, v := range s {
        if v != n {
            ret = append(ret, v)
        } else {
            ud = true
        }
    }

    return ret, ud
}

// Removes duplicate strings from slice
func Rmdupstrfslice(list []string) []string {

    var nlist []string

    for _, v := range list {
        if !Findstrinslice(v, nlist) {
            nlist = append(nlist, v)
        }
    }

    return nlist
}

// Retrieves client IP address from http request
func getclientip(r *http.Request) string {

    ip := r.Header.Get("X-Real-Ip")
    if ip == "" { ip = r.Header.Get("X-Forwarded-For") }
    if ip == "" { ip = r.RemoteAddr }

    return ip
}

// Returns true if user password validates
func Valuser(u User, pass []byte) bool {

    e := bcrypt.CompareHashAndPassword(u.Pass, pass)

    if e == nil { return true }
    return false
}

// Log file wrapper
func Addlog(ltype int, msg []byte, llev int, u User, r *http.Request) {

    ip := getclientip(r)
    var lentry string

    switch ltype {
    case L_REQ:
        if llev < 1 { return }

        if u.Name == "" {
            lentry = fmt.Sprintf("REQ from %s: %s", ip, msg)
        } else {
            lentry = fmt.Sprintf("REQ from %s (%s): %s", u.Name, ip, msg)
        }

    case L_RESP:
        if llev < 1 { return }

        if u.Name == "" {
            lentry = fmt.Sprintf("RESP to %s: %s", ip, msg)
        } else {
            lentry = fmt.Sprintf("RESP to %s (%s): %s", u.Name, ip, msg)
        }

    case L_SHUTDOWN:
        lentry = fmt.Sprintf("Server shutdown requested from %s (%s)", u.Name, ip)

    default:
        lentry = fmt.Sprintf("Unknown log request")
    }

    log.Println(lentry)
}

// Convert file size to human readable format
func Prettyfsize(b int64) string {

    k := b / 1024
    m := k / 1024
    var ret string

    if m == 0 {
        ret = fmt.Sprintf("%d.%dKB", k, (b % 1024) / 100)
    } else {
        ret = fmt.Sprintf("%d.%dMB", m, k / 100)
    }

    return ret
}

// Transforms whitespace separated string to tag slice
func Formattags(s string) []string {

    var ret []string

    itags := strings.Split(s, " ")

    for _, s := range itags {
        ret = append(ret, Cleanstring(s, RXTAGS))
    }

    return ret
}

// Makes a hard copy of a file on the file system
func Cp(s string, d string) (int64, error) {

    sf, e := os.Open(s)
    Cherr(e)
    defer sf.Close()

    df, e := os.Create(d)
    Cherr(e)
    defer df.Close()

    b, e := io.Copy(df, sf)

    return b, e
}

// Appends str to file at fname
func Appendfile(fname string, str string) {

    f, e := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    Cherr(e)
    defer f.Close()

    _, e = f.WriteString(str)
    Cherr(e)
}

// Removes all files residing in dir (except .gitkeep)
func Rmall(dir string) {

    d, e := os.Open(dir)
    Cherr(e)
    defer d.Close()

    fl, e := d.Readdirnames(-1)
    Cherr(e)

    for _, fn := range fl {
        if fn == ".gitkeep" { continue }
        e = os.RemoveAll(filepath.Join(dir, fn))
        Cherr(e)
    }
}

// Writes raw data to file
func Wrdatafile(fnp string, sf io.Reader) error {

    f, e := os.Create(fnp)
    Cherr(e)
    defer f.Close()

    Cherr(e)
    _, e = io.Copy(f, sf)

    return e
}

// Sends a status code response as JSON object
func Sendstatus(code int, text string, w http.ResponseWriter) {

    resp := Statusresp{
            Code: code,
            Text: text }

    enc := json.NewEncoder(w)
    enc.Encode(resp)
}

// Wrapper for tag status responses
func Sendtagstatus(r int, w http.ResponseWriter) {

    var sstr string
    if r != 0 { sstr = fmt.Sprintf("%d new tag(s) added", r) }
    Sendstatus(C_OK, sstr, w)
}

// Adds tag suffixes and returns new tag list
func Addtagsuf(tags []string, suf []string) []string {

    var ret []string

    for _, t := range tags {
        for _, s := range suf {
            tmp := t + s
            ret = append(ret, tmp)
        }
    }

    return ret
}

// Removes all tag suffixes and returns 'clean' list
func Striptagsuf(stags []string) []string {

    var ret []string

    for _, t := range stags {
        suf := filepath.Ext(t)
        ct := t[0:len(t) - len(suf)]
        ret = append(ret, ct)
    }

    return ret
}
