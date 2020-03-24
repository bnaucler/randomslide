package main

import (
    "os"
    "fmt"
    "flag"
    "bufio"
    "image"
    "bytes"
    "strconv"
    "io/ioutil"
    "path/filepath"

    "github.com/boltdb/bolt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
    "github.com/bnaucler/randomslide/lib/rsimage"
)

func readtxtfile(db *bolt.DB, fn string, tags []string) int {

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
            rsdb.Addtextwtags(raw, tags, db, rscore.Set.Bmax, rscore.BBUC)
            rscore.Set.Bmax++

        } else {
            rsdb.Addtextwtags(raw, tags, db, rscore.Set.Tmax, rscore.TBUC)
            rscore.Set.Tmax++
        }

        ret++
    }

    _, rscore.Set = rsdb.Tagstoindex(tags, rscore.Set)

    return ret
}

func readimg(db *bolt.DB, opath string, fl []string, tags []string) int {

    n := 0

    for _, fn := range fl {

        ext := filepath.Ext(fn)
        fnp := fmt.Sprintf("%s/%s", opath, fn)

        ibuf, e := ioutil.ReadFile(fnp)
        rscore.Cherr(e)

        fszr := bytes.NewReader(ibuf)
        i, _, e := image.Decode(fszr)
        if e != nil { continue }

        b := i.Bounds()

        isz, szok := rsimage.Getimgtype(b.Max.X, b.Max.Y)
        if !szok { continue }

        nfn := fmt.Sprintf("%s%s", rscore.Randstr(rscore.RFNLEN), ext)
        nfnp := fmt.Sprintf("%s%s", rscore.IMGDIR, nfn)

        rsz, b := rsimage.Scaleimage(i, isz, nfnp)

        if !rsz {
            _, e = rscore.Cp(fnp, nfnp)
            rscore.Cherr(e)
        }

        // Append the appropriate size tag to slice
        var suf []string
        suf = append(suf, rscore.SUFINDEX[isz])
        ttags := rscore.Addtagsuf(tags, suf)

        img := rsimage.Mkimgobj(nfn, ttags, b.Max.X, b.Max.Y, isz, rscore.Set)
        id := []byte(strconv.Itoa(rscore.Set.Imax))
        fmt.Printf("IMPORTING: %+v\n", img)

        n++
        rsdb.Wrimage(db, id, img)
        rsdb.Uilists(db, ttags, rscore.Set.Imax, rscore.IBUC)
        rscore.Set.Imax++
    }

    fmt.Printf("TAGS: %+v\n", tags)

    return n
}

func readimgdir(db *bolt.DB, dns []string, tags []string, verb bool)  {

    var n int

    for _, dn := range dns {
        d, e := os.Open(dn)
        rscore.Cherr(e)
        defer d.Close()

        if verb { fmt.Printf("Importing images from %s...\n", dn) }
        fl, e := d.Readdirnames(-1)
        rscore.Cherr(e)
        n = readimg(db, dn, fl, tags)
        if verb { fmt.Printf("%d images imported\n", n) }
    }

    _, rscore.Set = rsdb.Tagstoindex(tags, rscore.Set)
}

func readtextdir(db *bolt.DB, fn []string, tags []string, verb bool) {

    var n int
    for _, f := range fn {
        if verb { fmt.Printf("Importing %s...\n", f) }
        n = readtxtfile(db, f, tags)
        if verb { fmt.Printf("%d lines imported\n", n) }
    }
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
    rscore.Set = rsdb.Rsettings(db)
    tags := rscore.Formattags(*tptr)

    if *iptr == true {
        readimgdir(db, fn, tags, *vptr)

    } else {
        readtextdir(db, fn, tags, *vptr)
    }

    rsdb.Updatetindex(db)
    rsdb.Wrsettings(db, rscore.Set)
}
