#!/usr/bin/env bash

INTERVAL=30
SPIDFILE="data/rsserver.pid"
MPIDFILE="data/rsmonitor.pid"
RSLOGFILE="static/log/rsserver.log"
MONLOGFILE="static/log/rsmonitor.log"

startserver() {
    nohup bin/build.sh all >> /dev/null
    DATE=`date +'%Y/%m/%d %R:%S'`
    echo "$DATE: Restarting server" >> $MONLOGFILE
    nohup bin/rsserver $@ >> $MONLOGFILE &
}

cleanup() {
    rm $MPIDFILE
    DATE=`date +'%Y/%m/%d %R:%S'`
    echo "$DATE: Caught trap, cleaning up" >> $MONLOGFILE
}

trap cleanup EXIT
echo $$ > $MPIDFILE

while true; do
    git pull

    if [ ! -f $SPIDFILE ]; then startserver $@; fi

    PID=`cat $SPIDFILE`
    kill -0 $PID > /dev/null

    if [ $? -eq 1 ]; then startserver $@; fi

    sleep $INTERVAL
done
