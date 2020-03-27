#!/bin/bash

trap "rm server; kill 0" EXIT

go build -o server
./server -server=8000 &
./server -server=8001 &
./server -server=8002 &
./server -server=8003 -client=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:8080/api?key=Tom" &
curl "http://localhost:8080/api?key=Tom" &
curl "http://localhost:8080/api?key=Tom" &

wait