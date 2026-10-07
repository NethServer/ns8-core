#!/bin/bash

if [[ -z "$1" ]]; then
    echo "Missing parameter: append 'serve' or 'build'"
elif [[ "$1" = "serve" ]]; then
    yarn install && yarn serve
elif [[ "$1" == "build" ]]; then
    yarn install && yarn build
else
    echo "Parameter not recognized: '$1'. Only 'serve' or 'build' are allowed"
fi
