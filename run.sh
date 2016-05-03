#!/bin/bash

store_path="/home/forge/rfid_data"

case $1 in
start)
    mkdir -p $store_path
    supervisorctl reread
    supervisorctl update
    supervisorctl start small_tcp_server
;;
stop)
    supervisorctl stop small_tcp_server
;;
build)
    mkdir -p /opt/small_tcp_server
    go build small_tcp_server.go
    cp small_tcp_server /opt/small_tcp_server/
;;
esac
