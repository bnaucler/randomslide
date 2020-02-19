# randomslide v0.1A
Generate random slides for online presentation.

## Created by
Björn W Naucler (mail@bnaucler.se)  
Morgan Andersson (info@ameste.se)

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
/restart                <null>                  Graceful server shutdown.
                                                Requires VOLATILEMODE true

/gettags                <null>                  Retrieves list of all tags
                                                in database

/getdeck                                        Request for slide deck
                        amount                  # of slides requested
                        lang                    language code 'en', 'sv' etc
                        tags                    tags on which to base deck


/addtext                                        Adds new text to database
                        tags                    Which tags to associate text wit
                        ttext                   Title text
                        btext                   Body text

```

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for best practices and information on how to get involved in the project.

## License
MIT (do whatever you want)
