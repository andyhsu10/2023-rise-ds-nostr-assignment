# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.20.4-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /build
RUN apk add --no-cache --update alpine-sdk git make
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -v ./cmd/relay

######## Start a new stage from scratch #######
FROM alpine:3.14

WORKDIR /app/
COPY --from=builder /build/relay ./main

# Add bash.
RUN apk add --no-cache \
    bash \
    ca-certificates
EXPOSE 8080

# Command to run the executable
CMD ./main
