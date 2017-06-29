FROM alpine

ADD assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus_linux-amd64 /

RUN mkdir -p /etc/janus/apis && \
    mkdir -p /etc/janus/auth
    
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081/status || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus_linux-amd64"]
