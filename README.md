# randomslide v0.1A
Generate random slides for online presentation.

## Created by
Björn W Naucler (mail@bnaucler.se)  
Morgan Andersson (info@ameste.se)

## Dependencies
Go and [BoltDB](https://github.com/boltdb/bolt) - the code has been tested on Arch Linux 5.5 and FreeBSD 12, but should be fairly portable.

## Project purpose
Slideshow karaokae or other pranks. Actual real world usefullness can be questioned.

## Usage
Build the server with `bin/build.sh all` to build server and tools.  
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

### API reference

```
Endpoint:               Variables:              Comment:
/restart                                        Graceful server shutdown.
                                                Requires VOLATILEMODE true
                        wipe        string      Wipes database and images if "yes"


/gettags                <null>                  Retrieves list of all tags
                                                in database


/getdeck                                        Request for slide deck
                        id          int         Retrieve saved deck with id#
                        amount      int         # of slides requested
                        lang        string      language code 'en', 'sv' etc
                        tags        string      tags on which to base deck


/addtext                                        Adds new text to database
                        tags        string      Which tags to associate text with
                        ttext       string      Title text
                        btext       string      Body text

```

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for best practices and information on how to get involved in the project.

## License
MIT (do whatever you want)
