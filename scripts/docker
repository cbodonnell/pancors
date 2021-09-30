#!/bin/sh
docker run -d --env-file="$ENV.env" --network auth_net -p 127.0.0.1:8086:8080 --name pancors --restart unless-stopped cheebz/pancors
