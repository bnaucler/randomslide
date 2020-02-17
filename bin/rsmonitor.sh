#!/bin/sh

INTERVAL=30
PIDFILE="data/rsserver.pid"
RSLOGFILE="static/log/rsserver.log"
MONLOGFILE="static/log/rsmonitor.log"

startserver() {
    bin/build.sh all >> $MONLOGFILE
    DATE=`date +'%a | %Y-%m-%d | %R:%S'`
    echo "$DATE: Restarting server" >> $MONLOGFILE
    nohup bin/rsserver -v >> $MONLOGFILE &
}

while true; do
    git pull

    if [ ! -f $PIDFILE ]; then startserver; fi

    PID=`cat $PIDFILE`
    kill -0 $PID > /dev/null

    if [ $? -eq 1 ]; then startserver; fi

    sleep $INTERVAL
done
