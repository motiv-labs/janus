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

# Expose port 3000 to the host so we can access the gin proxy
EXPOSE 3000

# Now tell Docker what command to run when the container starts
CMD gin run
