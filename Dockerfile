# Build Stage
FROM golang:1.18-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o duplizer-cli .

# Run Stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/duplizer-cli .

# Set the entrypoint to the executable, allowing additional arguments to be passed
ENTRYPOINT ["./duplizer-cli"]

# Default command
CMD ["--help"]
