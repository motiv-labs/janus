FROM alpine

ADD assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus_linux-amd64 /

RUN mkdir -p /etc/janus/apis && mkdir -p /etc/janus/auth

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081 || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus_linux-amd64"]
