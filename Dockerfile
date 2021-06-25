FROM golang:1.16.5

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/pooladkhay/tinier-shortener-service

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 3011 to the outside world
EXPOSE 3011

# Run the executable
CMD ["tinier-shortener-service"]
