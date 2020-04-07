package main

import(
    "fmt"
    "flag"
    "strconv"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

func rtxt(prompt string) string {

    var tmp string

    fmt.Printf("%s: ", prompt)
	fmt.Scanln(&tmp)

    if len(tmp) < 1 {
        panic("No data read: exiting")
    }

    return tmp
}

func main() {

    dbptr := flag.String("d", rscore.DBNAME, "specify database to open")
    vptr := flag.Bool("v", rscore.VERBDEF, "verbose mode")
    flag.Parse()

    db := rsdb.Open(*dbptr)
    defer db.Close()

    settings := rsdb.Rsettings(db)

    a := rscore.Smtp{
        Admin:      rtxt("admin"),
        Server:     rtxt("server"),
        User:       rtxt("user"),
        Pass:       rtxt("pass"),
    }

    p, e := strconv.Atoi(rtxt("port"))
    if e != nil { panic("Could not read port number") }

    a.Port = p

    settings.Smtp = a
    if *vptr { fmt.Println(a) }

    rsdb.Wrsettings(db, settings)
}
