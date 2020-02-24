#!/bin/sh

TOOLS=("dbdump")

echo "Building rsserver"
go build -o bin/rsserver rsserver/rsserver.go

if [ "$1" = "all" ]; then
    for FN in $TOOLS; do
        echo "Building $FN"
        go build -o bin/$FN tools/$FN.go
    done
fi
