#!/bin/bash

protoc -I=. -I=$GOROOT/src --gogoslick_out=plugins=grpc:. messages.proto
#protoc -I=./ -I=$GOROOT/src --go_out=$GOROOT/src ./messages.proto
#protoc -I=./ -I=$GOPATH/src --go-grpc_out=$GOPATH/src ./messages.proto
#protoc -I=. -I=$GOROOT/src --gograin_out=. messages.proto
