#!/bin/sh

INTERVAL=30
PIDFILE="data/rsserver.pid"
RSLOGFILE="static/log/rsserver.log"
MONLOGFILE="static/log/rsmonitor.log"
PORT=6291

startserver() {
    bin/build.sh all >> $MONLOGFILE
    DATE=`date +'%a | %Y-%m-%d | %R:%S'`
    echo "$DATE: Restarting server" >> $MONLOGFILE
    nohup bin/rsserver -v -p $1 >> $MONLOGFILE &
}

usage() {
    echo "-p <port> - listen on port (default: $PORT)"
    exit 0
}

while getopts 'hp:' flag; do
    case "${flag}" in
        h) usage ;;
        p) PORT="${OPTARG}";;
    esac
done

while true; do
    git pull

    if [ ! -f $PIDFILE ]; then startserver $PORT; fi

    PID=`cat $PIDFILE`
    kill -0 $PID > /dev/null

    if [ $? -eq 1 ]; then startserver; fi

    sleep $INTERVAL
done
