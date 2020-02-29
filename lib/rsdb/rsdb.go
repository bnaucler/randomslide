package rsdb

/*

    Package of database operations
    Used by randomslide

*/

import (
    "fmt"
    "sort"
    "strconv"
    "math/rand"
    "encoding/json"
    "github.com/boltdb/bolt"
    "github.com/bnaucler/randomslide/lib/rscore"
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

// Writes image to database
func Wrimage(db *bolt.DB, k []byte, img rscore.Imgobj) {

    mimg, e := json.Marshal(img)
    e = Wrdb(db, k, mimg, rscore.IBUC)
    rscore.Cherr(e)
}

// Write user index to database
func Wruindex(db *bolt.DB, uname string,
    settings rscore.Settings) rscore.Settings {

    users := rscore.Uindex{}
    var mindex []byte

    if settings.Umax > 0 {
        mindex, e := Rdb(db, rscore.INDEX, rscore.UBUC)
        rscore.Cherr(e)
        e = json.Unmarshal(mindex, &users)
        rscore.Cherr(e)
    }

    users.Names = append(users.Names, uname)
    sort.Strings(users.Names)

    mindex, e := json.Marshal(users)

    e = Wrdb(db, rscore.INDEX, mindex, rscore.UBUC)
    rscore.Cherr(e)

    settings.Umax++
    return settings
}

// Returns deck from database
func Getdeckfdb(db *bolt.DB, deck rscore.Deck, req rscore.Deckreq,
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
    mxindex int, buc []byte) {

    to := rscore.Textobj{
            Id: mxindex,
            Text: text,
            Tags: tags }

    // Storing the object in db
    key := []byte(strconv.Itoa(mxindex))
    mtxt, e := json.Marshal(to)
    e = Wrdb(db, key, mtxt, buc)
    rscore.Cherr(e)

    // Update all relevant tag lists
    Uilists(db, tags, mxindex, buc)
}
