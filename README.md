![GitHub Release](https://img.shields.io/github/v/release/zapharaos/fihub-backend)
<a href="https://github.com/go-goyave/goyave/actions"><img src="https://github.com/go-goyave/goyave/workflows/CI/badge.svg" alt="Build Status"/></a>
[![codecov](https://codecov.io/gh/Zapharaos/fihub-backend/graph/badge.svg?token=BL7YP0GTK9)](https://codecov.io/gh/Zapharaos/fihub-backend)

![GitHub License](https://img.shields.io/github/license/zapharaos/fihub-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zapharaos/fihub-backend)](https://goreportcard.com/report/github.com/Zapharaos/fihub-backend)

# fihub-backend

The backend handles users' requests to list their financial transactions and provide analysis. It connects to brokers selected by the users to retrieve their assets and transactions.

## Dependencies

- Makefile is used to run commands 
```bash
sudo apt-get install make # If using WSL or any Linux distribution
choco install make # If using Powershell (not recommended, some commands may not work)
```

- Install gRPC [here](https://grpc.io/docs/languages/go/quickstart/).

- Goose is used for database migrations. Install it with the following command:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- Swagg is used to generate the swagger file. Install it with the following command:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Development

### Configuration

Open the `config/fihub-backend.toml` file and set the `APP_ENV` variable to `development`.

### Docker

This project is using Docker. Get started [here](https://www.docker.com/get-started).

#### Build

To build the project, you can use either of the following commands:
```bash
make build
make build-plain # build with verbose output progress=plain
```

#### Start

To start the project, you can use either of the following commands:
```bash
make dev
make dev-d # detached mode
make dev-b # build on top
make dev-bd # detached & build on top
```

Includes [Air](https://github.com/air-verse/air) for hot-reloading.

#### Debug

To debug the project, you can use either of the following commands:
```bash
make debug
make debug-d # detached mode
make debug-b # build on top
make debug-bd # detached & build on top
```

Includes [Delve](https://github.com/go-delve/delve) for debugging on top of Air.

When using JetBrains Goland, learn how to attach the debugger to a Go process that is running in a Docker container [here](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#attach-to-a-process-in-the-docker-container).

### Generation - gRPC

When editing the proto files, you need to regenerate the Go files. To do this, run the following command:

```bash
make proto-gen
```

### Swagger

Generate the swagger file (reused by frontend).

```bash
make swagger
```

### Migrations - PostgreSQL

To create a new migration file
```bash
goose create file_name sql
```

To apply the migrations
```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=fihub sslmode=disable" up
```

## Production

Work in progress.

### Configuration

Fill in the `config/fihub-backend.prod.toml` file by overriding variables. Don't forget to set the `APP_ENV` variable to `production` in the default .toml file.

### Run

```bash
make prod-bd # detached & build on top
```