package main

import (
    "fmt"
    "flag"
    "strconv"
    "encoding/json"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
    "github.com/bnaucler/randomslide/lib/rsimage"
)

// Retrieves all decks from database
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

// Retrieves text objects from database
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

// Retrieves image objects from database
func retrimg(db *bolt.DB, mxind int, buc []byte) []rscore.Imgobj {

    var ret []rscore.Imgobj

    for i := 0; i < mxind; i++ {

        cobj := rscore.Imgobj{}

        k := []byte(strconv.Itoa(i))
        v, e := rsdb.Rdb(db, k, buc)
        rscore.Cherr(e)

        json.Unmarshal(v, &cobj)
        ret = append(ret, cobj)
    }

    return ret
}

// Retrieves tag lists from database
func gettaglist(db *bolt.DB, tl []string, buc []byte) []string {

    var ret []string
    ctag := rscore.Iindex{}

    if rscore.Identicalbs(buc, rscore.IBUC) {
        suf := rsimage.Getallsuf()
        tl = rscore.Addtagsuf(tl, suf)
    }

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

// Retrieves user objects from database
func retrusers(db *bolt.DB, settings rscore.Settings, buc []byte) []rscore.User {

    var ret []rscore.User

    index := rsdb.Ruindex(db, settings)

    for _, v := range index.Names {
        u := rsdb.Ruser(db, v)
        ret = append(ret, u)
    }

    return ret
}

func main() {

    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    tlptr := flag.Bool("l", false, "tag lists")
    ttptr := flag.Bool("t", false, "title objects")
    btptr := flag.Bool("b", false, "body text objects")
    imptr := flag.Bool("i", false, "image objects")
    sptr := flag.Bool("s", false, "settings")
    dptr := flag.Bool("k", false, "decks")
    uptr := flag.Bool("u", false, "users")
    flag.Parse()

    db := rsdb.Open(*dbptr)

    if !*tlptr && !*ttptr && !*btptr && !*imptr &&
       !*sptr && !*dptr && !*uptr {
       fmt.Println("Exiting: No operation selected.")
        return
    }

    settings := rsdb.Rsettings(db)

    // DUMP SETTINGS
    if *sptr { fmt.Printf("SETTINGS: %+v\n", settings) }

    // DUMP TAG LISTS
    if settings.Tmax > 0 && *tlptr {
        ttl := gettaglist(db, settings.Taglist, rscore.TBUC)
        fmt.Printf("TTEXT TAG LIST: %+v\n", ttl)
    }

    if settings.Bmax > 0 && *tlptr {
        btl := gettaglist(db, settings.Taglist, rscore.BBUC)
        fmt.Printf("BTEXT TAG LIST: %+v\n", btl)
    }

    if settings.Imax > 0 && *tlptr {
        il := gettaglist(db, settings.Taglist, rscore.IBUC)
        fmt.Printf("IMAGE TAG LIST: %+v\n", il)
    }

    // DUMP TTEXT
    if settings.Tmax > 0 && *ttptr {
        ttext := retrtxt(db, settings.Tmax, rscore.TBUC)
        fmt.Printf("TTEXT: %+v\n", ttext)
    }

    // DUMP BTEXT
    if settings.Bmax > 0 && *btptr {
        btext := retrtxt(db, settings.Bmax, rscore.BBUC)
        fmt.Printf("BTEXT: %+v\n", btext)
    }

    // DUMP IMAGES
    if settings.Imax > 0 && *imptr {
        imgs := retrimg(db, settings.Bmax, rscore.IBUC)
        fmt.Printf("IMAGES: %+v\n", imgs)
    }

    // DUMP DECKS
    if settings.Dmax > 0 && *dptr {
        decks := retrtxt(db, settings.Dmax, rscore.DBUC)
        fmt.Printf("DECKS: %+v\n", decks)
    }

    // DUMP USERS
    if settings.Umax > 0 && *uptr {
        users := retrusers(db, settings, rscore.UBUC)
        fmt.Printf("USERS: %+v\n", users)
    }
}
