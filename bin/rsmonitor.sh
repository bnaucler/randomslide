#!/bin/sh

INTERVAL=30

while true; do
    git pull
    PID=`cat data/rsserver.pid`
    ISUP=`ps -p $PID | wc -l`
    if [ $ISUP == "1" ]; then
        bin/build.sh
        DATE=`date +'%a | %Y-%m-%d | %R:%S'`
        echo "$DATE: Restarting server"
        nohup bin/rsserver -v &
    fi
    sleep $INTERVAL
done
