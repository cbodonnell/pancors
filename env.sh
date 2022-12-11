#!/bin/bash

if [ "$1" == "help" ]; then
    echo "Usage: $0 <env_name>"
    exit 0
fi

ENV=$1

export $(grep -v '^#' $ENV.env | xargs)

