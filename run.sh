#!/bin/bash

store_path="/home/forge/rfid_data"

case $1 in
start)
    mkdir -p $store_path
    supervisorctl reread
    supervisorctl update
;;
stop)
    supervisorctl stop small_tcp_server
;;
build)
    go build small_tcp_server.go
;;
esac
