package rsdb

/*

    Package of database operations
    Used by randomslide

*/

import (
    "fmt"
    "sort"
    "strconv"
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

// Returns number of text objects per tag from db
func Countobj(db *bolt.DB, tn string, buc []byte) int {

    ttag := rscore.Tag{}
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

// Updates all relevant tag lists
func Updatetaglists(db *bolt.DB, tags []string, i int, buc []byte) {

    for _, s := range tags {
        ctag := rscore.Tag{}
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
    Updatetaglists(db, tags, mxindex, buc)
}
