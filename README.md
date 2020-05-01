
![randomslide logo](visuals/randomslide-08.png)


# randomslide v0.4B
Generate random slides for slideshow karaokae or other pranks. Actual real world usefulness can be questioned.

The service is currently in early beta stage and can be accessed [here](https://randomslide.com). For feedback - please use the builtin feedback function whenever possible.

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
bin/emailset
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

For API reference, take a look at [APIREF.md](APIREF.md).

## Themes
Place your css theme files in `static/css/themes`. Filenames starting with `_` will not be indexed. To create new theme files based on an empty template, there's a fairly well-commented `_empty.css` to use as a base.

## CLI tools
There are a few CLI helper tools bundled with randomslide:  
`emailset` configures the SMTP settings (will be removed once better solution is deployed)  
`dbdump` dumps the database to console. Not recommended with large data sets.  
`batchimport` can import a UTF-8-encoded text file or image directory directly into the database.  
`imgclass` iterates through a directory and checks for classes and final dimensions of images.

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for best practices and information on how to get involved in the project.

## License
MIT (do whatever you want)
