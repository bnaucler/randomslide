package main

import (
    "os"
    "fmt"
    "image"
    "image/png"
    "image/jpeg"
    "image/gif"
    "io/ioutil"

    "github.com/bnaucler/randomslide/lib/rscore"
)

func init() {

    image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
    image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}

func getclass(x int, y int) int {

    div := y / 10
    nx := x / div
    ny := y / div

    fmt.Printf("DEBUG: %dx%d\n", nx, ny)

    switch {
    case nx > 20:
        return 6 // Ultrawide
    case nx > 16:
        return 5 // Wide
    case nx > 14:
        return 4 // Normal
    case nx > 12:
        return 3 // Almost-box
    case nx > 10:
        return 2 // Box
    case nx > 5:
        return 1 // Portrait
    }

    return 0 // Ultra portrait
}

func main() {

    flist, e := ioutil.ReadDir(os.Args[1])
    rscore.Cherr(e)

    for _, fname := range flist {

        fp := fmt.Sprintf("%s/%s", os.Args[1], fname.Name())
        f, e := os.Open(fp)
        rscore.Cherr(e)
        defer f.Close()

        i, _, e := image.Decode(f)
        rscore.Cherr(e)
        b := i.Bounds()
        c := getclass(b.Max.X, b.Max.Y)

        fmt.Printf("%s: %dx%d (class %d)\n", fname.Name(), b.Max.X, b.Max.Y, c)
    }
}
