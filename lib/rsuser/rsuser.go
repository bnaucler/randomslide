package rsuser

/*

    Package of user related operations
    Used by randomslide

*/

import (
    "net/http"
    "encoding/json"

    "github.com/boltdb/bolt"
    "golang.org/x/crypto/bcrypt"

    "github.com/bnaucler/randomslide/lib/rscore"
    "github.com/bnaucler/randomslide/lib/rsdb"
)

// Validates skey - returns true if user is logged in
func Valskey(db *bolt.DB, uname string, skey string,
    umax int) (bool, rscore.User) {

    if umax < 1 { return false, rscore.User{} }

    u := rsdb.Ruser(db, uname)

    if skey == u.Skey { return true, u }
    return false, rscore.User{}
}

// Retrieves and validates user object
func Userv(db *bolt.DB, w http.ResponseWriter, umax int, c *rscore.Apicall,
    alevreq int) (bool, rscore.User) {

    if c.User == "" {
        rscore.Sendstatus(rscore.C_NLOG,
            "User not logged in - no username provided", w)
        return false, rscore.User{}
    }

    if !rsdb.Isindb(db, []byte(c.User), rscore.UBUC) {
        rscore.Sendstatus(rscore.C_NOSU, "No such user", w)
        return false, rscore.User{}
    }

    sok, u := Valskey(db, c.User, c.Skey, umax)

    if !sok {
        rscore.Sendstatus(rscore.C_NLOG,
            "User not logged in - skey mismatch", w)
        return false, rscore.User{}
    }

    if u.Alev < alevreq {
        rscore.Sendstatus(rscore.C_ALEV,
            "User does not have sufficient access level", w)
        return false, u
    }

    return true, u
}

// Returns true if initiated by admin or operation applied to initiating user
func Isadminorme(db *bolt.DB, settings rscore.Settings, c *rscore.Apicall,
    tu rscore.User, w http.ResponseWriter) bool {

    var ok bool

    if tu.Name == c.User {
        ok, _ = Userv(db, w, settings.Umax, c, rscore.ALEV_CONTRIB)
    } else {
        ok, _ = Userv(db, w, settings.Umax, c, rscore.ALEV_ADMIN)
    }

    return ok
}

// Sets user password
func Setpass(u rscore.User, pass string) (bool, rscore.User) {

    var e error

    if len(pass) < rscore.PWMINLEN { return false, u }

    u.Skey = rscore.Randstr(rscore.SKEYLEN)
    u.Pass, e = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    rscore.Cherr(e)

    return true, u
}

// Changes user password
func Chpass(db *bolt.DB, settings rscore.Settings, c *rscore.Apicall,
    tu rscore.User, w http.ResponseWriter) (bool, rscore.User) {

    ok := Isadminorme(db, settings, c, tu, w)
    if !ok { return false, tu }

    ok, tu = Setpass(tu, c.Pass)
    if !ok { rscore.Sendstatus(rscore.C_USPW, "Unsafe password", w) }
    return ok, tu
}

// Changes user admin status
func Chadminstatus(db *bolt.DB, op int, umax int, c *rscore.Apicall,
    tu rscore.User, w http.ResponseWriter) (bool, rscore.User) {

    var ok bool

    ok, _ = Userv(db, w, umax, c, rscore.ALEV_ADMIN)
    if !ok { return false, tu }

    switch {
    case op == rscore.CU_MKADM:
        tu.Alev = rscore.ALEV_ADMIN

    case op == rscore.CU_RMADM:
        tu.Alev = rscore.ALEV_CONTRIB

    default:
        return false, tu
    }

    return true, tu
}

// Removes user account from db
func Rmuser(db *bolt.DB, settings rscore.Settings, c *rscore.Apicall,
    tu rscore.User, w http.ResponseWriter) (bool, rscore.Settings) {

    ok := Isadminorme(db, settings, c, tu, w)

    if ok && settings.Umax > 1 {
        e := rsdb.Rmkv(db, []byte(tu.Name), rscore.UBUC)
        rscore.Cherr(e)
        settings = rsdb.Rmufrindex(db, tu.Name, settings)
        return ok, settings
    }

    return ok, settings
}

// Creates login object
func Getloginobj(u rscore.User) rscore.Login {

    ur := rscore.Login{}
    ur.Name = u.Name
    ur.Skey = u.Skey
    ur.Alev = u.Alev

    return ur
}

// Wrapper for sending user object to frontend
func Senduser(u rscore.User, r *http.Request, w http.ResponseWriter,
    settings rscore.Settings) {

    li := Getloginobj(u)
    ml, e := json.Marshal(li)
    rscore.Cherr(e)
    rscore.Addlog(rscore.L_RESP, ml, settings.Llev, u, r)

    enc := json.NewEncoder(w)
    enc.Encode(li)
}

