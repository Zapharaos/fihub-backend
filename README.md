![GitHub Release](https://img.shields.io/github/v/release/zapharaos/fihub-backend)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/zapharaos/fihub-backend/golang.yml)
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

- Gomock is used to generate mocks. Install it with the following command:
```bash
go install go.uber.org/mock/mockgen@latest
```

## Development

### Running with Docker

This project is using Docker. Get started [here](https://www.docker.com/get-started).

#### Configuration

- Create a `.env` file to override any variables from `config/fihub-backend.toml` that you need.
- Please note that each variable you wish to override must have the `FIHUB_` prefix.
- Don't forget to override the `FIHUB_APP_ENV` variable to `development`.

#### Build

To build the project, you can use either of the following commands:
```bash
make build
make build-plain # build with verbose output progress=plain
```

#### Databases

If you only need to run the databases, you can use the following command:
```bash
make db
```

#### Start

To start the whole project, you can use either of the following commands:
```bash
make dev
make dev-d # detached mode
make dev-b # build on top
make dev-bd # detached & build on top
```

Includes [Air](https://github.com/air-verse/air) for hot-reloading.

### Running without Docker

#### Databases

We recommend you to look into the [Running with Docker - Databases](#databases) section above to run only the required databases for debugging.

#### IDE - GoLand (recommended)

We recommend you to use GoLand as it is more convenient, especially for debugging. See [Run/debug configuration](https://www.jetbrains.com/help/go/run-debug-configuration.html).

- Start by creating a new `Go build` configuration.
- Set the `Package path` to `github.com/Zapharaos/fihub-backend/cmd/api`.
  - If you want to run another service, simply change `api` to the service you want to run.
- Enable `Run after build`.
- Set the `Working directory` to the root of the project `fihub-backend`.
- Add `Environment variables` to override any config variables that you need.

Note regarding environment variables:
- Any service might require different environment variable, so please check the `config/fihub-<service>.toml` file to see which one you need to set.
- Please note that each variable you wish to override must have the `FIHUB_` prefix.
- Don't forget to override the `FIHUB_APP_ENV` variable to `development`.

#### Command line (not recommended)

You can run the `api` with the following command:
```bash
FIHUB_APP_ENV=development go run ./cmd/api
````

- If you want to run another service, simply change `api` to the service you want to run.
- Any service might require different environment variable, so please check the `config/fihub-<service>.toml` file to see which one you need to set.
- Please note that each variable you wish to override must have the `FIHUB_` prefix.
- Don't forget to override the `FIHUB_APP_ENV` variable to `development`.

### Generation

#### Protobuf & gRPC

When editing the proto files, you need to regenerate the Go files. To do this, run the following command:

```bash
make proto-gen
```

#### Mocks

When editing the code (proto, repository, ...), you may need to regenerate the mocks. To do this, run the following command:
```bash
make mocks
```

#### Swagger

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

- Create a `.env` file to override any variables from `config/fihub-backend.toml` that you need.
- Please note that each variable you wish to override must have the `FIHUB_` prefix.
- Don't forget to comment any override of the `FIHUB_APP_ENV` variable.

### Run

```bash
make prod-bd # detached & build on top
```