#!/usr/bin/bash

echo "protoc --proto_path=proto --go_out=proto/ proto/order.proto proto/trade.proto"
protoc --proto_path=proto --go_out=proto/ proto/order.proto proto/trade.proto

echo "cp proto/*.pb.go apollo/proto"
cp proto/*.pb.go apollo/proto

echo "cp proto/*.pb.go matcher/proto"
cp proto/*.pb.go matcher/proto
