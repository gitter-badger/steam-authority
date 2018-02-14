#!/bin/sh

cd ../

git fetch origin
git reset --hard origin/master

# Update datatore index

dep ensure
go build

/etc/init.d/steam restart
