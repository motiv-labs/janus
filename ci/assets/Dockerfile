FROM golang:1.8-alpine

# Set apps home directory.
ENV APP_DIR /go/src/github.com/hellofresh/janus

# Creates the application directory
RUN mkdir -p $APP_DIR

# Add sources.
COPY . $APP_DIR

# Define current working directory.
WORKDIR $APP_DIR

# Build the go binary
RUN set -ex \
    && apk add --no-cache --update --virtual .build-deps make curl git \
    && make \
    # Clean up
    && apk del .build-deps \
    && rm -rf /var/cache/apk/* \
    && rm -rf /go/src/* \
    && rm -rf /go/pkg/*

EXPOSE 8080
CMD janus
