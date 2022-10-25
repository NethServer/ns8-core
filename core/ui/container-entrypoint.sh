#!/bin/bash

if [[ -z "$1" ]]; then
    echo "Missing parameter: append 'serve', 'build' or 'storybook'"
elif [[ "$1" = "serve" ]]; then
    yarn install && yarn serve
elif [[ "$1" == "build" ]]; then
    yarn install && yarn build
elif [[ "$1" == "storybook" ]]; then
    yarn install && yarn storybook
else
    echo "Parameter not recognized: '$1'. Only 'serve', 'build' or 'storybook' are allowed"
fi
