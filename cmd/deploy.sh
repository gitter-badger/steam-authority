#!/bin/sh

cd ../

git fetch origin
git reset --hard origin/master

dep ensure
go build

curl https://api.rollbar.com/api/1/deploy/ \
  -F access_token=${STEAM_ROLLBAR_PRIVATE} \
  -F environment=${ENV} \
  -F revision=$(git log -n 1 --pretty=format:"%H") \
  -F local_username=Jleagle \
  --silent > /dev/null

/etc/init.d/steam restart
