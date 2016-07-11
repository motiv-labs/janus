FROM golang

# Set apps home directory.
ENV APP_DIR $GOPATH/src/github.com/hellofresh/api-gateway

# Add sources.
ADD . $APP_DIR

# Define current working directory.
WORKDIR $APP_DIR

# Get gin from main repo
RUN go get github.com/codegangsta/gin

# Get godeps from main repo
RUN go get github.com/tools/godep

# Restore godep dependencies
RUN godep restore

RUN go install github.com/hellofresh/api-gateway

# Expose port 3000 to the host so we can access the gin proxy
EXPOSE 3000

# Now tell Docker what command to run when the container starts
CMD gin run
