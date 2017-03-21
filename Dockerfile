FROM scratch
ADD ci/assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus /
RUN mkdir -p /etc/janus/apis && mkdir -p /etc/janus/auth
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/janus"]
