#!/bin/sh

echo "Building rsserver"
go build -o bin/rsserver rsserver/rsserver.go

if [ "$1" = "all" ]; then
    echo "Building tools"
    go build -o bin/dbdump tools/dbdump.go
fi
