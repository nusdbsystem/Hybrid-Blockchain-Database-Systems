#!/usr/bin/env bash
# set -x

echo "Stop all bcdbnodes"
# pgrep -f "bcdbnode" || true
# pkill -f "bcdbnode"|| true
# kill -9 $(ps -ef|grep "bcdbnode"|grep -v "grep"|awk '{print $2}')
pbcdbnode=`ps -ef|grep "bcdbnode"|grep -v "grep"|wc -l`
echo ${pbcdbnode}
if (( ${pbcdbnode} > 0 )); then 
   kill $(ps -ef|grep "bcdbnode"|grep -v "grep"|awk '{print $2}')
fi
sleep 4
echo "##################### Stop bcdbnodes successfully! ##########################"
