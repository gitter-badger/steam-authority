#!/usr/bin/env bash

if [ "${ENV}" == "local" ]
then

    realize start

else

    steam-authority --pics --consumers

fi
