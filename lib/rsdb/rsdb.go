package rsdb

/*

    Package of database operations
    Used by randomslide

*/

import (
    "fmt"
    "encoding/json"
    "github.com/boltdb/bolt"
    "github.com/bnaucler/randomslide/lib/rscore"
)

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

