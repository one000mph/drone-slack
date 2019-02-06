FROM plugins/base:multiarch

LABEL maintainer="Andyet OPS <ops@andyet.com>" \
  org.label-schema.name="Drone Slack" \
  org.label-schema.schema-version="2.1.0"

ADD release/linux/amd64/drone-slack /bin/
ENTRYPOINT ["/bin/drone-slack"]
