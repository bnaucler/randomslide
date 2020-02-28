package main

import (
    "os"
    "fmt"
    "flag"
    "bufio"

    "github.com/boltdb/bolt"
    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func importfile(db *bolt.DB, fn string, tags []string,
    settings rscore.Settings) (rscore.Settings, int) {

    f, e := os.Open(fn)
    rscore.Cherr(e)
    defer f.Close()

    ret := 0

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        raw := scanner.Text()
        tlen := len(raw)

        if tlen > rscore.BTEXTMAX {
            continue

        } else if tlen > rscore.TTEXTMAX {
            rsdb.Addtextwtags(raw, tags, db, settings.Bmax, rscore.BBUC)
            settings.Bmax++

        } else {
            rsdb.Addtextwtags(raw, tags, db, settings.Tmax, rscore.TBUC)
            settings.Tmax++
        }

        ret++
    }

    _, settings = rsdb.Tagstoindex(tags, settings)

    return settings, ret
}

func readimg(db *bolt.DB, fn []string, tags []string,
    settings rscore.Settings, verb bool) rscore.Settings {


    return settings
}

func readtext(db *bolt.DB, fn []string, tags []string,
    settings rscore.Settings, verb bool) rscore.Settings {

    var n int
    for _, f := range fn {
        if verb { fmt.Printf("Importing %s...\n", f) }
        settings, n = importfile(db, f, tags, settings)
        if verb { fmt.Printf("%d lines imported\n", n) }
    }

    return settings
}

func main() {

    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    tptr := flag.String("t", "", "Tags to associate")
    vptr := flag.Bool("v", false, "verbose mode")
    iptr := flag.Bool("i", false, "Image dir import")
    flag.Parse()
    fn := flag.Args()

    db := rsdb.Open(*dbptr)
    defer db.Close()
    settings := rsdb.Rsettings(db)
    tags := rscore.Formattags(*tptr)

    if *iptr == true {
        settings = readimg(db, fn, tags, settings, *vptr)

    } else {
        settings = readtext(db, fn, tags, settings, *vptr)
    }

    rsdb.Wrsettings(db, settings)
}
