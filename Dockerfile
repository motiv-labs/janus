FROM alpine
ADD ci/assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus_linux-amd64 /
RUN mkdir -p /etc/janus/apis && mkdir -p /etc/janus/auth
EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus_linux-amd64"]
