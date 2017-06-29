#!/bin/sh

# cli options
# name - node name
# listen - listening address and port
# cluster - tuple of cluster peers; first peer is leader
# debug - t/f

$PWD/badkv \
  --name badkv1 \
  --listen localhost:10001 \
  --cluster badkv1=localhost:10001,badkv2=localhost:10002,badkv3=localhost:10003 \
  --debug

# todo
# working directory
# cluster name
# discovery
# cert store

# debug
go run main.go
go run main.go \
  --name badkv1 \
  --listen localhost:10001
go run main.go \
  --name badkv1 \
  --listen localhost:10001 \
  --cluster badkv1=localhost:10001,badkv2=localhost:10002,badkv3=localhost:10003

# test load; curl loop, apache bench
while true; do curl http://localhost:10001/api/v1/keys/mytenant/myname; done
curl -d "myvalue1" http://localhost:10001/api/v1/keys/mytenant/myname
ab -n 10000 -c 10 http://127.0.0.1:10001/api/v1/keys/mytenant/myname
ab -n 100000 -c 100 -k http://127.0.0.1:10001/api/v1/keys/mytenant/myname

# test GET
curl -v http://localhost:10001/api/v1/keys/mytenant/mykey

# test PUT
curl -v -d "bar1" http://localhost:10001/api/v1/keys/mytenant/foo1
curl http://localhost:10001/api/v1/keys/mytenant/foo1

# test DELETE
curl http://localhost:10001/api/v1/keys/mytenant/foo2
curl -d "bar1" http://localhost:10001/api/v1/keys/mytenant/foo2
curl http://localhost:10001/api/v1/keys/mytenant/foo2
curl -v -X DELETE http://localhost:10001/api/v1/keys/mytenant/foo2
curl http://localhost:10001/api/v1/keys/mytenant/foo2

# test cluster GET
curl -v http://localhost:10001/api/v1/keys/mytenant/mykey
curl http://localhost:10002/api/v1/keys/mytenant/mykey
curl http://localhost:10003/api/v1/keys/mytenant/mykey

# test cluster PUT
curl -v -d "bar1" http://localhost:10001/api/v1/keys/mytenant/foo1
curl -d "bar2" http://localhost:10002/api/v1/keys/mytenant/foo1
curl -d "bar3" http://localhost:10003/api/v1/keys/mytenant/foo1
curl http://localhost:10001/api/v1/keys/mytenant/foo1
curl http://localhost:10002/api/v1/keys/mytenant/foo1
curl http://localhost:10003/api/v1/keys/mytenant/foo1

# test cluster DELETE
curl http://localhost:10001/api/v1/keys/mytenant/foo2
curl -d "bar2" http://localhost:10001/api/v1/keys/mytenant/foo2
curl http://localhost:10002/api/v1/keys/mytenant/foo2
curl -v -X DELETE http://localhost:10002/api/v1/keys/mytenant/foo2
curl http://localhost:10001/api/v1/keys/mytenant/foo2
curl http://localhost:10002/api/v1/keys/mytenant/foo2
curl http://localhost:10003/api/v1/keys/mytenant/foo2
