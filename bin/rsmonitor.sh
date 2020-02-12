#!/bin/sh

INTERVAL=30
PIDFILE="data/rsserver.pid"
RSLOGFILE="static/log/rsserver.log"
MONLOGFILE="static/log/rsmonitor.log"

touch $PIDFILE
touch $RSLOGFILE
touch $MONLOGFILE

while true; do
    git pull
    PID=`cat $PIDFILE`
    kill -0 $PID > /dev/null
    if [ $? -eq 1 ]; then
        bin/build.sh >> $MONLOGFILE
        DATE=`date +'%a | %Y-%m-%d | %R:%S'`
        echo "$DATE: Restarting server" >> $MONLOGFILE
        nohup bin/rsserver -v >> $MONLOGFILE &
    fi
    sleep $INTERVAL
done
