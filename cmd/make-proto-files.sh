#!/bin/sh

protoc -I=../ --go_out=../ ./protos/xx.proto
