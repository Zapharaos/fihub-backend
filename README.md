# fihub-backend

## Development

### Docker

This project is using Docker. Get started [here](https://www.docker.com/get-started).

### Start

```bash
docker-compose up
```

Includes [Air](https://github.com/air-verse/air) for hot-reloading.

### Debug
```bash
docker-compose -f docker-compose.yml -f docker-compose.debug.yml up
```

Includes [Delve](https://github.com/go-delve/delve) for debugging on top of Air.

When using JetBrains Goland, learn how to attach the debugger to a Go process that is running in a Docker container [here](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#attach-to-a-process-in-the-docker-container).


## Production

Work in progress.

## Dependencies

Install goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Install swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Goose

To create a new migration file
```bash
goose create file_name sql
```

To apply the migrations
```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=fihub sslmode=disable" up
```


## Swagger

Generate the swagger file.

```bash
swag init -ot yaml
```

Display swagger-ui with the generated swagger file [at this local url](http://localhost:80/)

```bash
# Bash
docker run --rm -d -p 80:8080 -e SWAGGER_JSON=/tmp/swagger.yaml -v `pwd`/docs:/tmp swaggerapi/swagger-ui
# Powershell
docker run --rm -d -p 80:8080 -e SWAGGER_JSON=/tmp/swagger.yaml -v "$(pwd)/docs:/tmp" swaggerapi/swagger-ui
```

Generate the angular typescript client.

```bash
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/docs/swagger.yaml -g typescript-angular -o /local/docs/angular

# Powershell
docker run --rm -v ${PWD}:/local -v ${PWD}\..\fihub-ui\src\app\core\api:/local2 openapitools/openapi-generator-cli generate -i /local/docs/swagger.yaml -g typescript-angular -o /local2 
# Bash
docker run --rm -v ${PWD}/GolandProjects/caisse-back:/local -v ${PWD}/PhpstormProjects/caisse-front/src/app/core/api:/local2 openapitools/openapi-generator-cli generate -i /local/docs/swagger.yaml -g typescript-angular -o /local2
```