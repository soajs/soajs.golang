# soajs.golang

[![Build Status](https://travis-ci.org/soajs/soajs.golang.svg?branch=master)](https://travis-ci.org/soajs/soajs.golang)
[![Coverage Status](https://coveralls.io/repos/github/soajs/soajs.golang/badge.svg?branch=master)](https://coveralls.io/github/soajs/soajs.golang?branch=master)
[![GolangCI](https://golangci.com/badges/github.com/soajs/soajs.golang.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/soajs/soajs.golang)](https://goreportcard.com/report/github.com/soajs/soajs.golang)

SOAJS middleware for Golang services. This middleware provides integration between your Go REST services and the SOAJS framework.

## Requirements

- Go 1.21 or higher
- Go modules enabled

## Installation

```bash
go get github.com/soajs/soajs.golang
```

## Usage

Import the middleware in your Go service and use it with your HTTP router:

```go
import (
    "github.com/soajs/soajs.golang"
)

// Initialize and use the middleware in your application
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Running Tests with Coverage

```bash
go test -v -covermode=count -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

This project uses golangci-lint for code quality checks:

```bash
golangci-lint run --config .golangci.yml
```

## CI/CD

This project uses Travis CI for continuous integration, testing against:
- Go 1.21.x
- Go 1.22.x

## Documentation

See: https://soajsorg.atlassian.net/wiki/spaces/SOAJ/overview

## License

*Copyright SOAJS All Rights Reserved.*

Use of this source code is governed by an Apache license that can be found in the LICENSE file at the root of this repository.



