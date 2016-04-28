#!/bin/bash

case $1 in
start)
    mkdir -p /home/forge/rfid_data
    ./small_tcp_server &>/dev/null &
    echo $! > /tmp/small_tcp_server.pid
;;
stop)
    pkill small_tcp_server
;;
build)
    go build small_tcp_server.go
;;
esac
