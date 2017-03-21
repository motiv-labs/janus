FROM scratch
ADD ci/assets/ca-certificates.crt /etc/ssl/certs/
ADD dist/janus /
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/janus"]
