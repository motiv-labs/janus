FROM golang

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
EXPOSE 8080

CMD janus

