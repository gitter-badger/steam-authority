#!/bin/sh

cd ../

# Get the latest version
echo "### Pulling"
git fetch origin
git reset --hard origin/master

# Build
echo "### Building"
dep ensure
go build

# Copy over crontab
cp ./crontab /etc/cron.d/steamauthority

# Tell Rollbar
echo "### Rollbar"
curl https://api.rollbar.com/api/1/deploy/ \
  -F access_token=${STEAM_ROLLBAR_PRIVATE} \
  -F environment=${ENV} \
  -F revision=$(git log -n 1 --pretty=format:"%H") \
  -F local_username=Jleagle \
  --silent > /dev/null

# Restart web server & PICS
echo "### Restart"
/etc/init.d/steam restart
