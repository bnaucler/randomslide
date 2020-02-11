#!/bin/sh

INTERVAL=30
LOGFILE="data/rsserver.log"

while true; do
    git pull
    PID=`cat data/rsserver.pid`
    ISUP=`ps -p $PID | wc -l`
    if [ $ISUP == "1" ]; then
        bin/build.sh
        DATE=`date +'%a | %Y-%m-%d | %R:%S'`
        echo "$DATE: Restarting server"
        nohup bin/rsserver -v > $LOGFILE&
    fi
    sleep $INTERVAL
done
