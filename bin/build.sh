#!/usr/bin/env bash

TOOLS=("dbdump batchimport imgclass emailset")
TPATH="tools"

echo "Building rsserver"
go build -o bin/rsserver rsserver.go

if [ "$1" = "all" ]; then
    for FN in $TOOLS; do
        echo "Building $FN"
        go build -o bin/$FN $TPATH/$FN.go
    done
fi
