FROM golang AS builder
ENV JANUS_BUILD_ONLY_DEFAULT=1
ENV JANUS_SRC="/go/src/github.com/hellofresh/janus"

ADD . ${JANUS_SRC}

RUN cd ${JANUS_SRC} && \
    make

# ---
FROM alpine

ADD assets/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/hellofresh/janus/dist/janus /

RUN mkdir -p /etc/janus/apis && \
    mkdir -p /etc/janus/auth
    
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081/status || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus"]
