# Build stage: build the Go application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache make gcc musl-dev git

# Set the working directory inside the container
WORKDIR /lc

# Copy the entire source code, including submodules
COPY . .

# Initialize and update submodules
RUN git submodule update --init --recursive

# Download Go module dependencies
RUN go mod download

# Run the make command to build the application
RUN make build-light

# Final stage: use the ipfs/kubo image and add the Go application
FROM ipfs/kubo:v0.29.0

# Copy the built Go application from the build stage
COPY --from=builder /lc/bin/light-client /usr/local/bin/light-client

# Expose the necessary IPFS ports
EXPOSE 4001 5001 8080

# Set default environment variables (can be overridden)
ENV CLIENT_ID="default-client-id"

# Override the entrypoint to use a shell
ENTRYPOINT ["/bin/sh", "-c", "ipfs daemon --enable-gc & exec /usr/local/bin/light-client --loglevel debug --rpc-url wss://moonbase-alpha.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e --contract 0x916B54696A70588a716F899bE1e8f2A5fFd5f135 --topic-id DAS-TO-BQ --gcp-creds-file /lc/test/data/gcp-credentials.json --client-id $CLIENT_ID"]