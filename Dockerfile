FROM golang:1.23-alpine AS build

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

FROM build AS development

# Install : Air  = hot-reload ; Delve = debugger
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

# Run the application with Air for hot-reloading
CMD ["air", "-c", ".air.toml"]

FROM build AS build-production

# Copy whole project (shared + microservices) to the build stage
WORKDIR /app
COPY . .

# Build the Go app
RUN go build -v -o server

FROM scratch AS production
# Start a new ligthwheight stage from scratch

WORKDIR /

# Copy the binary from the build stage
COPY --from=build-production /app/server /server

# Copy config files from the build stage
COPY --from=build-production /app/config /config

# Run the compiled binary
CMD ["/server"]
