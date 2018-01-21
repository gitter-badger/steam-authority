#!/bin/sh

FILES=/path/to/*
for f in $PROJECT/cmd/protos/*.proto
do
  echo "Processing $f file..."

  protoc \
  --proto_path=../ \
  --go_out=../generated/ \
  $f

done
