FROM alpine AS builder

ARG JANUS_VERSION

RUN echo $JANUS_VERSION

RUN apk update \
    && apk add --virtual .build-deps wget tar ca-certificates \
	&& apk add libgcc openssl \
    && wget -O janus.tar.gz https://github.com/hellofresh/janus/releases/download/${JANUS_VERSION}/janus_linux-amd64.tar.gz \
    && tar -xzf janus.tar.gz -C /tmp

# ---
FROM alpine

ADD assets/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp/janus_linux-amd64 /

RUN mkdir -p /etc/janus/apis && \
    mkdir -p /etc/janus/auth
    
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8081/status || exit 1

EXPOSE 8080 8081 8443 8444
ENTRYPOINT ["/janus_linux-amd64"]
