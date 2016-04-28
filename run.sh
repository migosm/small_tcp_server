#!/bin/bash

case $1 in
start)
    mkdir -p /home/forge/rfid_data
    ./small_tcp_server &
;;
build)
    go build small_tcp_server.go
;;
esac
