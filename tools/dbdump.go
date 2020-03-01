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

func retrdeck(db *bolt.DB, mxind int, buc []byte) []rscore.Deck {

    var ret []rscore.Deck

    for i := 0; i < mxind; i++ {

        cobj := rscore.Deck{}

        k := []byte(strconv.Itoa(i))
        v, e := rsdb.Rdb(db, k, buc)
        rscore.Cherr(e)

        json.Unmarshal(v, &cobj)
        ret = append(ret, cobj)
    }

    return ret
}

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

func gettaglist(db *bolt.DB, tl []string, buc []byte) []string {

    var ret []string
    ctag := rscore.Iindex{}

    tl = append(tl, rscore.IKEY...)
    for _, t := range tl {

        k := []byte(t)
        v, e := rsdb.Rdb(db, k, buc)
        rscore.Cherr(e)

        json.Unmarshal(v, &ctag)

        cstr := fmt.Sprintf("%s:%s - %+v", t, string(buc), ctag)
        ret = append(ret, cstr)
    }

    return ret
}

func main() {

    if len(os.Args) < 2 { os.Exit(1) }

    db := rsdb.Open(os.Args[1])

    settings := rsdb.Rsettings(db)

    // DUMP SETTINGS
    fmt.Printf("SETTINGS: %+v\n", settings)

    // DUMP TAG LISTS
    if settings.Tmax > 0 {
        ttl := gettaglist(db, settings.Taglist, rscore.TBUC)
        fmt.Printf("TTEXT TAG LIST: %+v\n", ttl)
    }

    if settings.Bmax > 0 {
        btl := gettaglist(db, settings.Taglist, rscore.BBUC)
        fmt.Printf("BTEXT TAG LIST: %+v\n", btl)
    }

    if settings.Imax > 0 {
        il := gettaglist(db, settings.Taglist, rscore.IBUC)
        fmt.Printf("IMAGE TAG LIST: %+v\n", il)
    }

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

    // DUMP DECKS
    if settings.Dmax > 0 {
        decks := retrtxt(db, settings.Dmax, rscore.DBUC)
        fmt.Printf("DECKS: %+v\n", decks)
    }

}
