#!/bin/bash

PATH="$PATH:$(go env GOPATH)/bin"
export PATH

# gen clean up
rm -rf tmp/gen/*
mkdir -p tmp/gen/go

# language specific clean up. Normally deleting old protos
rm -rf pkg/pb/*
rm -rf gen/java/src/main/java/*

# Make Golang protos

find ./proto -name '*.proto' -print0 | while IFS= read -r -d '' file
do
  echo "Generating protos for $file"
  protoc -I proto --go_out=tmp/gen/go --go-grpc_out=tmp/gen/go "$file"
#  protoc -I proto --java_out=gen/java/src/main/java "$file"
done

# Move Golang protos to the right place
mv tmp/gen/go/github.com/emortalmc/kurushimi/pkg/pb/* pkg/pb

# Clean up
rm -rf tmp

# Java steps

rm -rf gen/java/src/main/proto
mkdir -p gen/java/src/main/proto
cp -r proto/* gen/java/src/main/proto

cd gen/java || exit 1
./gradlew clean generateProto
