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
    "github.com/bnaucler/randomslide/lib/rsimage"
)

func init() {

    image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
    image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
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
        c, szok := rsimage.Getimgtype(b.Max.X, b.Max.Y)
        ni, rsz := rsimage.Scaleimage(i, c)
        b = ni.Bounds()

        if szok {
            fmt.Printf("%s: %dx%d (class %d, scaled: %v)\n",
                fname.Name(), b.Max.X, b.Max.Y, c, rsz)

        } else {
            fmt.Printf("%s: %dx%d (Could not be classified)\n",
                fname.Name(), b.Max.X, b.Max.Y)
        }
    }
}
