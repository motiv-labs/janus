FROM alpine

ADD assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus_linux-amd64 /
ADD dist/healthchecker-linux-amd64 /

RUN mkdir -p /etc/janus/apis && mkdir -p /etc/janus/auth

HEALTHCHECK --interval=5s --timeout=5s CMD ["./healthchecker-linux-amd64", "-port=8080"] || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus_linux-amd64"]
