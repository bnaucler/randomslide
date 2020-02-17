package main

import (
    "os"
    "fmt"
    "strconv"
    "encoding/json"

    "github.com/boltdb/bolt"
    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func retrtxt(db *bolt.DB, mxind int, buc []byte) []rscore.Textobj {

    var ret []rscore.Textobj

    for i := 0; i < mxind; i++ {

        cobj := rscore.Textobj{}

        k := []byte(strconv.Itoa(i))
        v, e := rsdb.Rdb(db, k, buc)
        rscore.Cherr(e)

        json.Unmarshal(v, &cobj)
        ret = append(ret, cobj)
    }

    return ret
}

func gettaglist(db *bolt.DB, tl []string, buc []byte) []rscore.Tag {

    var ret []rscore.Tag
    ctag := rscore.Tag{}

    for _, t := range tl {
        k := []byte(t)
        v, e := rsdb.Rdb(db, k, buc)
        rscore.Cherr(e)

        json.Unmarshal(v, &ctag)
        ret = append(ret, ctag)
    }

    return ret
}

func main() {

    if len(os.Args) < 2 { os.Exit(1) }

    db, e := bolt.Open(os.Args[1], 0640, nil)
    rscore.Cherr(e)
    defer db.Close()

    settings := rsdb.Rsettings(db)

    // DUMP SETTINGS
    fmt.Printf("SETTINGS: %+v\n", settings)

    // DUMP TAG LISTS
    ttl := gettaglist(db, settings.Taglist, rscore.TBUC)
    fmt.Printf("TTEXT TAG LIST: %+v\n", ttl)

    btl := gettaglist(db, settings.Taglist, rscore.BBUC)
    fmt.Printf("BTEXT TAG LIST: %+v\n", btl)

    // DUMP TTEXT
    if settings.Tmax > 0 {
        ttext := retrtxt(db, settings.Tmax, rscore.TBUC)
        fmt.Printf("TTEXT: %+v\n", ttext)
    }

    // DUMP BTEXT
    if settings.Bmax > 0 {
        btext := retrtxt(db, settings.Bmax, rscore.BBUC)
        fmt.Printf("BTEXT: %+v\n", btext)
    }
}
