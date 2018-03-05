#!/usr/bin/env bash

if [ "${ENV}" == "local" ]
then

    bash ./cmd/pics.sh >> /dev/null
    realize start

else

    steam-authority --pics --consumers

fi
