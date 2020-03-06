package rsimage

/*

    Package of image operations
    Used by randomslide

*/

import (
    "os"
    "fmt"
    "image"
    "image/png"
    "image/jpeg"
    "image/gif"
    "path/filepath"

    "github.com/nfnt/resize"
    "github.com/bnaucler/randomslide/lib/rscore"
)

// Returns path of new image file
func Newimagepath(ofn string) (string, string) {

    ext := filepath.Ext(ofn)
    fn := fmt.Sprintf("%s%s", rscore.Randstr(rscore.RFNLEN), ext)
    fnp := fmt.Sprintf("%s%s", rscore.IMGDIR, fn)

    return fn, fnp
}

// Writes image object to file
func Wrimagefile(i image.Image, fnp string) error {

    var e error
    ext := filepath.Ext(fnp)

    f, e := os.Create(fnp)
    rscore.Cherr(e)
    defer f.Close()

    switch {
    case ext == ".jpg" || ext == ".jpeg":
        jpeg.Encode(f, i, nil)

    case ext == ".png":
        png.Encode(f, i)

    case ext == ".gif":
        gif.Encode(f, i, nil)

    }

    return e
}

// Conditionally returns image size type & true if fitting classification
// TODO refactor
func Getimgtype(iw int, ih int) (int, bool) {

    var t int
    var ok bool

    w := uint(iw)
    h := uint(ih)

    div := ih / 10
    nw := iw / div

    switch {
    case nw > 20:
        return 4, false

    case nw > 12: // Landscape
        if w > rscore.IMGMAX[0][0] || h > rscore.IMGMAX[0][1] {
            t = 0
            ok = true

        } else if w < rscore.IMGMIN[0][0] || h < rscore.IMGMIN[0][1] {
            ok = false

        } else {
            t = 1
        }

    case nw > 8: // Box-shaped
        t = 2

        if w < rscore.IMGMIN[2][0] || h < rscore.IMGMIN[2][1] { ok = false
        } else { ok = true }

    case nw > 5: // Portrait
        t = 3

        if w < rscore.IMGMIN[3][0] || h < rscore.IMGMIN[3][1] { ok = false
        } else { ok = true }

    default:
        return 4, false

    }

    return t, ok
}

func Mkimgobj(fn string, tags []string, iw int, ih int, szt int,
    settings rscore.Settings) rscore.Imgobj {

    img := rscore.Imgobj{
        Id: settings.Imax,
        Fname: fn,
        Tags: tags,
        W: iw,
        H: ih,
        Size: szt,
    }

    return img
}

// Scales image down to max dimensions allowed, returns true if image was scaled
func Scaleimage(i image.Image, t int) (image.Image, bool) {

    b := i.Bounds()

    if uint(b.Max.X) > rscore.IMGMAX[t][0] || uint(b.Max.Y) > rscore.IMGMAX[t][1] {
        rsz := resize.Thumbnail(rscore.IMGMAX[t][0], rscore.IMGMAX[t][1],
            i, resize.Lanczos3)
        return rsz, true
    }

    return i, false
}

