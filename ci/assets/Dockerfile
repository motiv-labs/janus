FROM golang:1.7-alpine

# Install make and curl
RUN apk --update add make curl git

# Set apps home directory.
ENV APP_DIR /go/src/github.com/hellofresh/janus

# Creates the application directory
RUN mkdir -p $APP_DIR

# Add sources.
COPY . $APP_DIR

# Define current working directory.
WORKDIR $APP_DIR

# Build the go binary
RUN make

# Clean apk cache
RUN rm -rf /var/cache/apk/*

EXPOSE 8080
CMD janus
