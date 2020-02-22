#!/usr/bin/env bash

SPIDFILE="data/rsserver.pid"
MPIDFILE="data/rsmonitor.pid"
SERVER="localhost"
PORT=6291

usage() {
    echo "$0 usage:"
    echo "-h: this message"
    echo "-f <arg>: specify server pidfile Location (default: $SPIDFILE)"
    echo "-m: SIGQUIT monitor"
    echo "-r: remove server pidfile"
    echo "-k: SIGQUIT server"
    echo "-x: SIGKILL all servers"
    echo "-p <arg>: set port (default: $PORT)"
    exit 0
}

skill() {
    PID=`cat $1`
    kill -9 $PID
}

while getopts 'hfmrkp:xs' flag; do
    case "${flag}" in
        h) usage ;;
        f) SPIDFILE="${OPTARG}" ;;
        r) rm $SPIDFILE $MPIDFILE;;
        k) skill $SPIDFILE ;;
        m) skill $MPIDFILE ;;
        p) PORT="${OPTARG}" ;;
        s) SERVER="${OPTARG}" ;;
        x) killall -9 rsserver ;;
        *) exit 1 ;;
    esac
done

# Default to VOLATILE soft kill
curl "$SERVER:$PORT/restart"
