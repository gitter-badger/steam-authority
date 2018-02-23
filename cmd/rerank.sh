#!/bin/sh

echo $(curl --basic --user $STEAM_AUTH_USER:$STEAM_AUTH_PASS --show-error --silent "http://$STEAM_LOCAL_DOMAIN/admin/rerank")
