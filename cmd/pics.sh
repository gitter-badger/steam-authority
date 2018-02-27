#!/bin/sh

# Stop old script
pkill -f steam-pics-api

# Start a new one
cd ${STEAM_PICS_PATH}
nohup npm start &
