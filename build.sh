#!/bin/bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack \
&& docker build --rm -t one000mph/drone-slack:$1 .
