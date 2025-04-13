FROM golang:1.23-alpine AS build

# Copy shared files from server, used by this microservice
WORKDIR /app
COPY go.mod go.sum ./

# For this microservice
WORKDIR /app/cmd/health
COPY cmd/health/go.mod cmd/health/go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -modfile=./go.mod && go mod verify

FROM build AS development

# Install: Air = hot-reload; Delve = debugger
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Copy whole project (shared + microservices) to the build stage
WORKDIR /app
COPY . .

# Run the microservice with Air for hot-reloading
WORKDIR /app/cmd/health
COPY config/ ./config
CMD ["air", "-c", ".air.linux.toml"]

FROM build AS build-production

# Copy whole project (shared + microservices) to the build stage
WORKDIR /app
COPY . .

# Build the Go microservice
WORKDIR /app/cmd/health
RUN go build -v -o health

FROM scratch AS production
# Start a new lightweight stage from scratch

WORKDIR /

# Copy the binary from the build stage
COPY --from=build-production /app/cmd/health/health /health

# Copy config files from the build stage
COPY --from=build-production /app/config /config

# Run the compiled binary
CMD ["/health"]