# Golang Base Server

## Introduction

This repo is meant as a base starting point for a Go backend server.
It includes an opinionated project structure and basic user-management-related implementations. 

## Requirements

- Go 1.20
- A MySQL server

## Included implementations

- User CRUD operations (Some may be missing)
- Session management
- Email verification 
- Simple license management

## Dependency injection

`wire` is used for dependency injection. Check out the docs here: https://github.com/google/wire.

Providers are declared in each abstraction layer (ie. handler, service, etc.) in the files with name `wire.go`.
Add to those files if you need to declare new providers.

To use `wire` to generate code:
```shell
wire path/to/init_server/module
```

## Code structure/design

The structure of the repo roughly follows the MVC design pattern. 
To reduce unnecessary complexity, only two main abstraction layers are included (handler -> service).

### Handler

Contains only http-related logic:
- Parsing of request
- Handling of HTTP response
- Authentication
- etc.

### Service

Contains business logic. Data access layer is also placed here.
DB accesses should be done in this layer.

### Lib

`lib` should contain other functions/logic that do not belong in `handler` or `service`.

## Local development

1. Check `conf/config.yaml`. Make sure that the configurations are correct for your local system. 
(Check if the DB DSN matches your local DB's username and password)

2. Execute `go run cmd/main.go`. The server will run on port `8080`.

## Documentation

Swagger documentation can be found at `/swagger/index.html`.

To generate documentation based on comment tags, run the following command from project root:

```shell
swag init --parseDependency --dir cmd,handler
```

