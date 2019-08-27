FROM golang:1.12-alpine AS builder

ARG VERSION='0.0.1-docker'

WORKDIR /janus

COPY . ./

RUN apk add --update bash make git
RUN export JANUS_BUILD_ONLY_DEFAULT=1 && \
    export VERSION=$VERSION && \
    make build

# ---

FROM alpine

COPY --from=builder /janus/dist/janus /

RUN apk add --no-cache ca-certificates
RUN mkdir -p /etc/janus/apis && \
    mkdir -p /etc/janus/auth

RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081/status || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus", "start"]
