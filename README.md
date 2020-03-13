# randomslide v0.3A
Generate random slides for slideshow karaokae or other pranks. Actual real world usefulness can be questioned.

## Created by
Björn W Nauclér (bwn@randomslide.com)  
Morgan Andersson (ma@randomslide.com)

## Dependencies
[BoltDB](https://github.com/boltdb/bolt), [bcrypt](https://golang.org/x/crypto/bcrypt) & [nfnt/resize](https://github.com/nfnt/resize)  
The code has been tested on Arch Linux 5.5 and FreeBSD 12, but should be fairly portable.

## Installation
To build server and tools:  
```
go get github.com/bnaucler/randomslide
cd $GOPATH/src/github.com/bnaucler/randomslide
bin/build.sh all
```

## Usage
Launch with `bin/rsserver`  
Or use `bin/rsmonitor.sh` to automatically restart the server after `/restart`.

Output of `bin/rsserver -h`:  
```
Usage of bin/rsserver:
  -d string
    	specify database to open (default "./data/rs.db")
  -p int
    	port number to listen (default 6291)
  -v	increase log level
```

The first user who registers an account will automatically be provided with admin rights.  
Server log files can be accessed at `static/logs` or in the admin interface.

## CLI tools
There are a few CLI helper tools bundled with randomslide:  
`dbdump` dumps the database to console. Not recommended with large data sets.  
`batchimport` can import a UTF-8-encoded text file or image directory directly into the database.  
`imgclass` iterates through a directory and checks for classes and final dimensions of images.

### API reference

```
Endpoint:               Variables:              Comment:
/restart*                                       Graceful server shutdown.
                                                Requires admin rights
                        wipe        string      Wipes database and images if "yes"


/register                                       Register a user account
                        user        string      Username
                        pass        string      Password
                        email       string      Email address


/login                                          Login to get higher access level
                        user        string      Username
                        pass        string      Password


/getusers               <null>                  Retrieves list of all users
                                                in database


/gettags                <null>                  Retrieves list of all tags
                                                in database


/getdeck                                        Request for slide deck
                        id          int         Retrieve saved deck with id#
                        amount      int         # of slides requested
                        lang        string      language code 'en', 'sv' etc
                        tags        string      tags on which to base deck


/addtext*                                       Adds new text to database
                        tags        string      Which tags to associate text with
                        ttext       string      Title text
                        btext       string      Body text


/addimg*                                        Adds new images to the database
                        file        file        The image file itself
                        tags        string      Which tags to associate the image with


/chuser*                                        Change user settings
                                                Some ops requires admin rights
                        tuser       string      User to edit
                        pass        string      New password (if applicable)
                        op          int         Operation:
                                                    0: Make admin
                                                    1: Remove admin rights
                                                    2: Change password
                                                    3: Remove user account

/feedback*                                      Give feedback on user experience
                        fb          string      The feedback info itself

```
Endpoints marked with `*` requires the user to be logged in, authenticated by `user` & `skey` being included with the request.

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for best practices and information on how to get involved in the project.

## License
MIT (do whatever you want)
