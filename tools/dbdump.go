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

func main() {

    if len(os.Args) < 2 { os.Exit(1) }

    db, e := bolt.Open(os.Args[1], 0640, nil)
    rscore.Cherr(e)
    defer db.Close()

    settings := rsdb.Rsettings(db)

    // DUMP SETTINGS
    fmt.Printf("SETTINGS: %+v\n", settings)

    // DUMP TTEXT
    ttext := retrtxt(db, settings.Tmax, rscore.TBUC)
    fmt.Printf("TTEXT: %+v\n", ttext)

    // DUMP BTEXT
    btext := retrtxt(db, settings.Bmax, rscore.BBUC)
    fmt.Printf("BTEXT: %+v\n", btext)


}
