#!/bin/sh

curl --basic --user ${STEAM_AUTH_USER}:${STEAM_AUTH_PASS} --silent "$STEAM_LOCAL_DOMAIN/admin/ranks"
