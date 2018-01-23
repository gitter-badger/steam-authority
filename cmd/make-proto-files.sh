#!/bin/sh

for f in protos/*.proto
do
  echo "Processing $f file..."

  protoc \
  --proto_path=protos \
  --go_out=../generated/ \
  $f

done
