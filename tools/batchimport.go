package main

import (
    "os"
    "fmt"
    "time"
    "flag"
    "bufio"
    "image"
    "bytes"
    "strconv"
    "image/png"
    "image/jpeg"
    "image/gif"
    "math/rand"
    "io/ioutil"
    "path/filepath"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func init() {
    rand.Seed(time.Now().UnixNano())

    image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
    image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}


func readtxtfile(db *bolt.DB, fn string, tags []string,
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

func readimg(db *bolt.DB, opath string, fl []string, tags []string,
    settings rscore.Settings) (int, rscore.Settings) {

    n := 0

    for _, fn := range fl {

        ext := filepath.Ext(fn)
        fnp := fmt.Sprintf("%s/%s", opath, fn)

        ibuf, e := ioutil.ReadFile(fnp)
        rscore.Cherr(e)

        fszr := bytes.NewReader(ibuf)
        ic, _, e := image.DecodeConfig(fszr)
        if e != nil { continue }

        isz, szok := rscore.Getimgtype(ic.Width, ic.Height)
        if !szok { continue }

        nfn := fmt.Sprintf("%s%s", rscore.Randstr(rscore.RFNLEN), ext)
        nfnp := fmt.Sprintf("%s%s", rscore.IMGDIR, nfn)


        _, e = rscore.Cp(fnp, nfnp)
        rscore.Cherr(e)

        // Append the appropriate size tag to slice
        ttags := append(tags, rscore.IKEY[isz])
        img := rscore.Mkimgobj(nfn, ttags, ic.Width, ic.Height, isz, settings)
        id := []byte(strconv.Itoa(settings.Imax))
        fmt.Printf("DEBUG: %+v\n", img)

        n++
        rsdb.Wrimage(db, id, img)
        rsdb.Uilists(db, ttags, settings.Imax, rscore.IBUC)
        settings.Imax++
    }

    fmt.Printf("DEBUG: %+v\n", tags)
    return n, settings
}

func readimgdir(db *bolt.DB, dns []string, tags []string,
    settings rscore.Settings, verb bool) rscore.Settings {

    var n int

    for _, dn := range dns {
        d, e := os.Open(dn)
        rscore.Cherr(e)
        defer d.Close()

        if verb { fmt.Printf("Importing images from %s...\n", dn) }
        fl, e := d.Readdirnames(-1)
        rscore.Cherr(e)
        n, settings = readimg(db, dn, fl, tags, settings)
        if verb { fmt.Printf("%d images imported\n", n) }
    }

    _, settings = rsdb.Tagstoindex(tags, settings)
    return settings
}

func readtextdir(db *bolt.DB, fn []string, tags []string,
    settings rscore.Settings, verb bool) rscore.Settings {

    var n int
    for _, f := range fn {
        if verb { fmt.Printf("Importing %s...\n", f) }
        settings, n = readtxtfile(db, f, tags, settings)
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
        settings = readimgdir(db, fn, tags, settings, *vptr)

    } else {
        settings = readtextdir(db, fn, tags, settings, *vptr)
    }

    rsdb.Wrsettings(db, settings)
}
