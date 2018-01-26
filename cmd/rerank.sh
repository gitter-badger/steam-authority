#!/bin/sh

echo $(curl --basic --user $STEAM_ADMIN_USER:$STEAM_ADMIN_PASS --show-error --silent "http://localhost:8085/admin/rerank")
