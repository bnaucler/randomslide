package rsdb

/*

    Package of database operations
    Used by randomslide

*/

import (
    "fmt"
    "sort"
    "strconv"
    "net/http"
    "math/rand"
    "encoding/json"
    "path/filepath"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsimage"
)

// Opens the database
func Open(dbname string) *bolt.DB {

    db, e := bolt.Open(dbname, 0640, nil)
    rscore.Cherr(e)

    return db
}

// Write JSON encoded byte slice to DB
func Wrdb(db *bolt.DB, k []byte, v []byte, cbuc []byte) (e error) {

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
func Rdb(db *bolt.DB, k []byte, cbuc []byte) (v []byte, e error) {

    e = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(cbuc)
        if b == nil { return fmt.Errorf("No bucket!") }

        v = b.Get(k)
        return nil
    })
    return
}

// Unconditionaly removes k/v pair from database
func Rmkv(db *bolt.DB, k []byte, cbuc []byte) (e error) {

    e = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(cbuc)
        if b == nil { return fmt.Errorf("No bucket!") }

        e = b.Delete(k)
        if e != nil { return e }

        return nil
    })
    return
}

// Wrapper for writing settings to database
func Wrsettings(db *bolt.DB, settings rscore.Settings) {

    mset, e := json.Marshal(settings)
    rscore.Cherr(e)

    e = Wrdb(db, rscore.INDEX, mset, rscore.SBUC)
    rscore.Cherr(e)
}

// Wrapper for reading settings from database
func Rsettings(db *bolt.DB) rscore.Settings {

    settings := rscore.Settings{}

    mset, e := Rdb(db, rscore.INDEX, rscore.SBUC)
    if e != nil { return rscore.Settings{} }

    e = json.Unmarshal(mset, &settings)
    rscore.Cherr(e)

    return settings
}

// Writes user to database
func Wruser(db *bolt.DB, u rscore.User) {

    mu, e := json.Marshal(u)
    rscore.Cherr(e)
    e = Wrdb(db, []byte(u.Name), mu, rscore.UBUC)
    rscore.Cherr(e)
}

// Returns user account from database
func Ruser(db *bolt.DB, uname string) rscore.User {

    u := rscore.User{}

    mu, e := Rdb(db, []byte(uname), rscore.UBUC)
    rscore.Cherr(e)
    e = json.Unmarshal(mu, &u)

    if e == nil { return u }
    return rscore.User{}
}

// Writes deck to database
func Wrdeck(db *bolt.DB, deck rscore.Deck) {

    k := []byte(strconv.Itoa(deck.Id))
    mdeck, e := json.Marshal(deck)
    rscore.Cherr(e)
    e = Wrdb(db, k, mdeck, rscore.DBUC)
}

// Writes image to database
func Wrimage(db *bolt.DB, k []byte, img rscore.Imgobj) {

    mimg, e := json.Marshal(img)
    e = Wrdb(db, k, mimg, rscore.IBUC)
    rscore.Cherr(e)
}

// Returns user index
func Ruindex(db *bolt.DB, settings rscore.Settings) rscore.Uindex {

    users := rscore.Uindex{}

    if settings.Umax > 0 {
        mindex, e := Rdb(db, rscore.INDEX, rscore.UBUC)
        rscore.Cherr(e)
        e = json.Unmarshal(mindex, &users)
        rscore.Cherr(e)
    }

    return users
}

// Writes tag index to db
func Wrtindex(db *bolt.DB, tindex rscore.Tagresp) { // TODO rename tagresp

    mti, e := json.Marshal(tindex)
    rscore.Cherr(e)

    e = Wrdb(db, rscore.TINDEX, mti, rscore.SBUC)
    rscore.Cherr(e)
}

// Reads tag index from db
func Rtindex(db *bolt.DB) rscore.Tagresp {

    t := rscore.Tagresp{}

    mti, e := Rdb(db, rscore.TINDEX, rscore.SBUC)
    rscore.Cherr(e)
    e = json.Unmarshal(mti, &t)

    if e == nil { return t }
    return rscore.Tagresp{}
}

// Returns updated rtag object
func Updatertag(db *bolt.DB, t string) rscore.Rtag {

    var rtag rscore.Rtag

    rtag.Name = t

    if rscore.Set.Tmax > 0 { rtag.TN = Countobj(db, t, rscore.TBUC) }
    if rscore.Set.Bmax > 0 { rtag.BN = Countobj(db, t, rscore.BBUC) }
    if rscore.Set.Imax > 0 { rtag.IN = Imgobjctr(db, t) }

    return rtag
}

// Updates tag index
func Updatetindex(db *bolt.DB) {

    tindex := rscore.Tagresp{}

    rscore.Smut.Lock()

    for _, t := range rscore.Set.Taglist {
        rtag := Updatertag(db, t)
        tindex.Tags = append(tindex.Tags, rtag)
    }

    Wrtindex(db, tindex)
    rscore.Smut.Unlock()
}

// Write user index to database
func Wruindex(db *bolt.DB, users rscore.Uindex,
    settings rscore.Settings) rscore.Settings {

    mindex, e := json.Marshal(users)
    e = Wrdb(db, rscore.INDEX, mindex, rscore.UBUC)
    rscore.Cherr(e)

    settings.Umax = len(users.Names)
    return settings
}

// Append uname to database index
func Addutoindex(db *bolt.DB, uname string,
    settings rscore.Settings) rscore.Settings {

    users := Ruindex(db, settings)

    users.Names = append(users.Names, uname)
    users.Names = rscore.Rmdupstrfslice(users.Names)
    sort.Strings(users.Names)

    rscore.Smut.Lock()
    settings = Wruindex(db, users, settings)
    rscore.Smut.Unlock()

    return settings
}

// Removes user name from index
func Rmufrindex(db *bolt.DB, uname string,
    settings rscore.Settings) rscore.Settings {

    users := Ruindex(db, settings)
    nusers := rscore.Uindex{}

    for _, v := range users.Names {
        if !rscore.Findstrinslice(v, nusers.Names) && v != uname {
            nusers.Names = append(nusers.Names, v)
        }
    }

    settings = Wruindex(db, nusers, settings)

    return settings
}

// Returns deck from database
func Rdeck(db *bolt.DB, deck rscore.Deck, id int,
    settings rscore.Settings) rscore.Deck {

    if id >= settings.Dmax || settings.Dmax < 1 {
        return rscore.Deck{}
    }

    bk := []byte(strconv.Itoa(id))
    mdeck, e := Rdb(db, bk, rscore.DBUC)
    rscore.Cherr(e)

    e = json.Unmarshal(mdeck, &deck)
    rscore.Cherr(e)

    return deck
}

// Creates valid selection list from tags
func Mksel(db *bolt.DB, tags []string, buc []byte) []int {

    var sel []int
    ctags := rscore.Iindex{}

    for _, t := range tags {
        bt := []byte(t)
        mtags, e := Rdb(db, bt, buc)
        rscore.Cherr(e)

        json.Unmarshal(mtags, &ctags)
        sel = append(sel, ctags.Ids...)
    }

    return sel
}

// Returns a random key based on tag list
func Getkeyfromsel(db *bolt.DB, tags []string, buc []byte, kmax int) []byte {

    sel := Mksel(db, tags, buc)
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
func Getrndtxt(db *bolt.DB, kmax int, tags []string, buc []byte) string {

    if kmax < 2 { return "" }

    k := Getkeyfromsel(db, tags, buc, kmax)

    txt := rscore.Textobj{}
    mtxt, e := Rdb(db, k, buc)
    rscore.Cherr(e)
    e = json.Unmarshal(mtxt, &txt)
    rscore.Cherr(e)

    return txt.Text
}

// Sends random image url from database, based on requested tags
func Getrndimg(db *bolt.DB, kmax int, tags []string, buc []byte) rscore.Imgobj {

    if kmax < 2 { return rscore.Imgobj{} }

    k := Getkeyfromsel(db, tags, buc, kmax)

    img := rscore.Imgobj{}
    mimg, e := Rdb(db, k, buc)
    rscore.Cherr(e)
    e = json.Unmarshal(mimg, &img)
    rscore.Cherr(e)

    return img
}

// Returns true if key returns something from database
func Isindb(db *bolt.DB, k []byte, buc []byte) bool {

    v, e := Rdb(db, k, buc)

    if len(v) == 0 || e != nil { return false }
    return true
}

// Returns number of text objects per tag from db
func Countobj(db *bolt.DB, tn string, buc []byte) int {

    ttag := rscore.Iindex{}
    k := []byte(tn)

    v, e := Rdb(db, k, buc)
    rscore.Cherr(e)

    json.Unmarshal(v, &ttag)

    return len(ttag.Ids)
}

// Updates index to include new tags
func Tagstoindex(tags []string, settings rscore.Settings) (int, rscore.Settings) {

    r := 0

    for _, t := range tags {
        if rscore.Findstrinslice(t, settings.Taglist) == false {
            settings.Taglist = append(settings.Taglist, t)
            r++
        }
    }

    if r != 0 { sort.Strings(settings.Taglist) }

    return r, settings
}

// Updates all relevant index lists
func Uilists(db *bolt.DB, tags []string, i int, buc []byte) {

    for _, s := range tags {
        ctag := rscore.Iindex{}
        key := []byte(s)

        resp, e := Rdb(db, key, buc)
        rscore.Cherr(e)

        json.Unmarshal(resp, &ctag)
        ctag.Ids = append(ctag.Ids, i)

        dbw, e := json.Marshal(ctag)
        e = Wrdb(db, key, dbw, buc)
        rscore.Cherr(e)
    }
}

// Conditionally adds tagged text to database
func Addtextwtags(text string, tags []string, db *bolt.DB,
    uname string, mxindex int, buc []byte) {

    to := rscore.Textobj{
            Id: mxindex,
            Text: text,
            Contr: uname,
            Tags: tags }

    // Storing the object in db
    key := []byte(strconv.Itoa(mxindex))
    mtxt, e := json.Marshal(to)
    e = Wrdb(db, key, mtxt, buc)
    rscore.Cherr(e)

    // Update all relevant tag lists
    Uilists(db, tags, mxindex, buc)
}

// Stores image object in database TODO make work for batchimport
func Addimgwtags(db *bolt.DB, fn string, iw int, ih int, isz int, contr string,
    tags []string, w http.ResponseWriter, settings rscore.Settings) rscore.Settings {

    ofn := filepath.Base(fn)

    // Write image object to database
    ttags := append(tags, rscore.SUFINDEX[isz])
    img := rsimage.Mkimgobj(ofn, ttags, iw, ih, isz, contr, settings)
    k := []byte(strconv.Itoa(img.Id))
    Wrimage(db, k, img)

    // Update relevant tags
    Uilists(db, ttags, settings.Imax, rscore.IBUC)
    settings.Imax++

    return settings
}

// Image wrapper for rsdb.Countobj()
func Imgobjctr(db *bolt.DB, t string) int {

    var tl []string
    var ret int

    tl = append(tl, t)
    suf := rsimage.Getallsuf()
    stl := rscore.Addtagsuf(tl, suf)

    for _, st := range stl {
        ret += Countobj(db, st, rscore.IBUC)
    }

    return ret
}
