#!/bin/bash
# # !/usr/bin/expect

set -e

echo "go build file ..."
GOOS=linux GOARCH=amd64 go build -o ./temp/run

echo "scp run file to remote ..."
scp ./temp/run root@121.40.139.31:/home/qianuuu/scmj/scmj


expect -c '
spawn ssh root@121.40.139.31
expect "]*" 
send "cd /home/qianuuu/scmj\n"
send "./scmj \n"
interact
'
# interact
