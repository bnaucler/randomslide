# randomslide v0.2A
Generate random slides for online presentation.

## Created by
Björn W Naucler (bwn@randomslide.com)  
Morgan Andersson (ma@randomslide.com)

## Dependencies
[BoltDB](https://github.com/boltdb/bolt) & [bcrypt](https://golang.org/x/crypto/bcrypt)  
The code has been tested on Arch Linux 5.5 and FreeBSD 12, but should be fairly portable.

## Project purpose
Slideshow karaokae or other pranks. Actual real world usefullness can be questioned.

## Installation
```
go get github.com/boltdb/bolt
go get golang.org/x/crypto/bcrypt
go get github.com/bnaucler/randomslide
```

`bin/build.sh all` to build server and tools.  

## Usage
Launch with `bin/rsserver`  
You can also use `bin/rsmonitor.sh` to automatically restart the server.

Output of `bin/rsserver -h`:  
```
Usage of bin/rsserver:
  -d string
    	specify database to open (default "./data/rs.db")
  -p int
    	port number to listen (default 6291)
  -v	verbose mode
```

Server log files can be accessed at `static/logs` or in the admin interface.

## CLI tools
`dbdump` will, perhaps unsurprisingly, dump the database to console. Not recommended with large data sets, but useful when troubleshooting.  
`batchimport` can import a UTF-8-encoded text file as text objects directly into the database.

### API reference

```
Endpoint:               Variables:              Comment:
/restart                                        Graceful server shutdown.
                                                Requires launching with -x
                        wipe        string      Wipes database and images if "yes"


/register                                       Register a user account
                        user        string      Username
                        pass        string      Password


/login                                          Login to get higher access level
                        user        string      Username
                        pass        string      Password


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


/feedback*                                      Give feedback on user experience
                        fb          string      The feedback info itself

```
Endpoints marked with `*` requires the user to be logged in.

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for best practices and information on how to get involved in the project.

## License
MIT (do whatever you want)
