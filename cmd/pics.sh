#!/bin/sh

echo "Stop old PICS"
pkill -f steam-pics-api

echo "Start new PICS"
nohup npm start --prefix ${STEAM_PICS_PATH} &
