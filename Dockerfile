####### Start from a golang base image ###############
FROM golang:1.13.6-buster as builder
LABEL maintainer="Motiv Labs <dev@motivsolutions.com>"
WORKDIR /app
COPY ./ ./

RUN go mod download

RUN make build

FROM ubuntu:20.04 as prod

COPY --from=builder /app/cassandra/schema.sql /usr/local/bin

COPY --from=builder /app/dist/janus /bin/janus
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

# just to have it
RUN ["/bin/janus", "--version"]
