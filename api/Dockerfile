############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

ENV GO111MODULE=on

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/nicholasjackson/cloud-pong
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/cloud-pong ./api/main.go

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/cloud-pong /go/bin/cloud-pong
# Run the hello binary.
ENTRYPOINT ["/go/bin/cloud-pong"]
