#!/bin/sh

cd ../

# Get the latest version
git fetch origin
git reset --hard origin/master

# Build
dep ensure
go build

# Copy over crontab
cp ./crontab /etc/cron.d/steamauthority

# Restart PICS
chmod +x ./cmd/pics.sh
./cmd/pics.sh

# Tell Rollbar
curl https://api.rollbar.com/api/1/deploy/ \
  -F access_token=${STEAM_ROLLBAR_PRIVATE} \
  -F environment=${ENV} \
  -F revision=$(git log -n 1 --pretty=format:"%H") \
  -F local_username=Jleagle \
  --silent > /dev/null

# Restart web server
/etc/init.d/steam restart
