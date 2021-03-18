FROM ubuntu:20.04

COPY schema.sql /usr/local/bin

ADD janus /bin/janus
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
