#!/bin/bash

store_path="/home/forge/rfid_data"

case $1 in
start)
    mkdir -p $store_path
    ./small_tcp_server -port=8085 -path=$store_path &>>/var/log/small_tcp_server.log &
    echo $! > /tmp/small_tcp_server.pid
;;
stop)
    pkill small_tcp_server
;;
build)
    go build small_tcp_server.go
;;
esac
