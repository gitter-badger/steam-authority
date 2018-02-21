#!/usr/bin/env bash

cd ../

if [ "${ENV}" = "local" ]; then

    realize start --run x -pics -consumers

else

    steam-authority -pics -consumers

fi
