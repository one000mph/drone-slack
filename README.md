# drone-slack

Drone plugin for sending Slack notifications. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack/).

## NOTE:

This is a fork of the original project, customized for Talky Inc. deploy
pipeline.


## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack
docker build --rm -t plugins/slack .
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e SLACK_WEBHOOK=https://hooks.slack.com/services/... \
  -e PLUGIN_CHANNEL=foo \
  -e PLUGIN_USERNAME=drone \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=hello-world \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  -e DRONE_TAG=1.0.0 \
  plugins/slack
```
