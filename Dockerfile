FROM golang:1.14-alpine AS builder

ARG VERSION='0.0.1-docker'
ENV VERSION $VERSION

WORKDIR /janus

COPY . ./

RUN apk add --update bash make git
RUN export JANUS_BUILD_ONLY_DEFAULT=1 make build

# ---

FROM ubuntu:20.04

COPY --from=builder /janus/dist/janus /bin/janus
RUN chmod a+x /bin/janus && \
    mkdir -p /etc/janus/apis && \
    mkdir -p /etc/janus/auth

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
  && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081/status || exit 1

# Use nobody user + group
USER 65534:65534

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/bin/janus", "start"]
