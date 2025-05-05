FROM golang:1.23-alpine AS development

# Install: Air = hot-reload; Delve = debugger
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Copy project files
WORKDIR /app
COPY . .

# Run the microservice with Air for hot-reloading
CMD ["air", "-c", "cmd/transaction/.air.toml"]

FROM golang:1.23-alpine AS build-production

# Copy whole project (shared + microservices) to the build stage
WORKDIR /app
COPY . .

# Build the Go microservice
WORKDIR /app/cmd/transaction
RUN go build -v -o transaction

FROM scratch AS production
# Start a new lightweight stage from scratch

WORKDIR /

# Copy the binary from the build stage
COPY --from=build-production /app/cmd/transaction/transaction /transaction

# Copy config files from the build stage
COPY --from=build-production /app/config /config

# Run the compiled binary
CMD ["/transaction"]