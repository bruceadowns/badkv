# Use goreman to run `go get github.com/mattn/goreman`
badkv1: ./badkv --name badkv1 --listen localhost:10001 --cluster badkv1=localhost:10001,badkv2=localhost:10002,badkv3=localhost:10003
badkv2: ./badkv --name badkv2 --listen localhost:10002 --cluster badkv1=localhost:10001,badkv2=localhost:10002,badkv3=localhost:10003
badkv3: ./badkv --name badkv3 --listen localhost:10003 --cluster badkv1=localhost:10001,badkv2=localhost:10002,badkv3=localhost:10003
