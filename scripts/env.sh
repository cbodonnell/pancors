#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <env_name>"
    exit 1
fi

ENV=$1

export $(grep -v '^#' $ENV.env | xargs)

