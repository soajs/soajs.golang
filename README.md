# soajs.golang

[![Build Status](https://travis-ci.org/soajs/soajs.golang.svg?branch=master)](https://travis-ci.org/soajs/soajs.golang)
[![Coverage Status](https://coveralls.io/repos/github/soajs/soajs.golang/badge.svg?branch=master)](https://coveralls.io/github/soajs/soajs.golang?branch=master)
[![GolangCI](https://golangci.com/badges/github.com/soajs/soajs.golang.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/soajs/soajs.golang)](https://goreportcard.com/report/github.com/soajs/soajs.golang)

SOAJS middleware for Golang services. This package provides seamless integration between your Go REST services and the SOAJS (Service Oriented Architecture JavaScript) framework, enabling registry management, request context handling, and multi-tenancy support.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Basic Setup](#basic-setup)
  - [Using Config](#using-config)
  - [Accessing SOAJS Context](#accessing-soajs-context)
  - [Registry Methods](#registry-methods)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Development](#development)
- [CI/CD](#cicd)
- [Documentation](#documentation)
- [License](#license)

## Features

- **Registry Management**: Automatic service registry synchronization with SOAJS
- **Auto-reload**: Configurable automatic registry reloading
- **Multi-tenancy**: Built-in tenant and application context handling
- **Request Context**: Access tenant, user (URAC), device, and geo information per request
- **Database Management**: Access to core and tenant meta databases through registry
- **Resource Discovery**: Service and resource lookup capabilities
- **HTTP Middleware**: Easy integration with standard Go HTTP handlers

## Requirements

- Go 1.21 or higher
- Go modules enabled
- SOAJS infrastructure (Controller, Registry API)

## Installation

```bash
go get github.com/soajs/soajs.golang
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"

    soajsgo "github.com/soajs/soajs.golang"
)

func main() {
    ctx := context.Background()

    // Initialize registry
    registry, err := soajsgo.New(ctx, "myservice", "dev", "service", true)
    if err != nil {
        log.Fatal(err)
    }

    // Create handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Access SOAJS context data
        soaData := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)

        w.Write([]byte("Hello from " + soaData.Tenant.Code))
    })

    // Apply middleware
    http.Handle("/", registry.Middleware(handler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Examples

Complete working examples are available in the [examples](./examples) directory:

- **[basic_service.go](./examples/basic_service.go)** - Full microservice using standard `net/http`
- **[gin_service.go](./examples/gin_service.go)** - Microservice using the Gin Web Framework

Each example demonstrates:
- Service initialization and configuration
- Middleware integration
- Accessing SOAJS context data (tenant, user, device, geo)
- Database and service discovery
- Custom registry configuration
- Health checks and graceful shutdown

See the [examples README](./examples/README.md) for detailed documentation and usage instructions.

## Usage

### Basic Setup

Create a registry connection using service name, environment code, and service type:

```go
import (
    "context"
    soajsgo "github.com/soajs/soajs.golang"
)

ctx := context.Background()

// Parameters: context, serviceName, envCode, serviceType, autoReload
registry, err := soajsgo.New(ctx, "myservice", "dev", "service", true)
if err != nil {
    log.Fatal(err)
}
```

### Using Config

Initialize registry from a configuration struct:

```go
config := soajsgo.Config{
    ServiceName:    "myservice",
    ServiceGroup:   "mygroup",
    ServicePort:    8080,
    ServiceIP:      "127.0.0.1",
    Type:           "service",
    ServiceVersion: "1",
    // ... additional config fields
}

registry, err := soajsgo.NewFromConfig(ctx, config)
if err != nil {
    log.Fatal(err)
}
```

### Accessing SOAJS Context

Extract SOAJS data from the request context:

```go
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    soaData := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)

    // Access tenant information
    tenantCode := soaData.Tenant.Code
    tenantId := soaData.Tenant.ID

    // Access user information (URAC)
    if soaData.Urac != nil {
        username := soaData.Urac.Username
        userId := soaData.Urac.ID
    }

    // Access device and geo information
    device := soaData.Device
    geo := soaData.Geo

    // Access registry for databases and services
    registry := soaData.Reg
})
```

### Registry Methods

The registry provides several methods for accessing databases, services, resources, and custom configurations:

```go
// Get a database by name
db, err := registry.Database("mydb")
if err == nil {
    prefix := db.Prefix
    servers := db.Servers
}

// Get all databases
dbs, err := registry.Databases()

// Get a service by name
service, err := registry.Service("myservice")
if err == nil {
    port := service.Port
    host := service.Host
}

// Get a resource by name
resource, err := registry.Resource("myresource")

// Get a specific custom registry by name
custom, err := registry.GetCustom("myCustom")
if err == nil {
    customReg := custom.(*soajsgo.CustomRegistry)
    value := customReg.Value
    locked := customReg.Locked
}

// Get all custom registries
allCustom, err := registry.GetCustom("")
if err == nil {
    customRegistries := allCustom.(soajsgo.CustomRegistries)
    for name, customReg := range customRegistries {
        // Access custom registry data
        fmt.Println(name, customReg.Value)
    }
}

// Manually reload registry
err = registry.Reload()
```

## Configuration

The `Config` struct supports the following fields:

```go
type Config struct {
    ServiceName           string       // Service name (required)
    ServiceGroup          string       // Service group (required)
    ServicePort           int          // Service port (required)
    ServiceIP             string       // Service IP address
    Type                  string       // Service type: "service" or "daemon" (required)
    ServiceVersion        string       // Service version (required)
    SubType               string       // Service subtype
    Description           string       // Service description
    Oauth                 bool         // OAuth enabled
    Urac                  bool         // URAC enabled
    UracProfile           bool         // URAC profile enabled
    UracACL               bool         // URAC ACL enabled
    UracConfig            bool         // URAC config enabled
    UracGroupConfig       bool         // URAC group config enabled
    TenantProfile         bool         // Tenant profile enabled
    ProvisionACL          bool         // Provision ACL enabled
    ExtKeyRequired        bool         // External key required
    RequestTimeout        int          // Request timeout
    RequestTimeoutRenewal int          // Request timeout renewal
    Awareness             bool         // Awareness enabled
    Maintenance           Maintenance  // Maintenance configuration
    SwaggerFilename       string       // Swagger / OpenAPI file served on GET /mcp/tools (default /swagger.yml)
}
```

## MCP Tools (GET /mcp/tools)

The SOAJS gateway can list a service's APIs as MCP tools for AI agents, filtered by each tenant's
package and each user's ACL. A service built with this middleware describes its APIs with its
Swagger 2.0 or OpenAPI 3 document; the middleware serves it on `GET /mcp/tools` and the gateway
turns it into tools.

- `NewFromConfig` sets it up from the config: `"swaggerFilename": "/swagger.yml"` (the default),
  relative to the folder the service runs from.
- With `New`, set it yourself:

  ```go
  reg.MCP = &soajsgo.MCPTools{ServiceName: "my-service", ServiceVersion: "1", SwaggerFilename: "/api/swagger.json"}
  ```

  `SwaggerDocument` (the document content) wins over `SwaggerFilename`. Without `reg.MCP` the
  middleware leaves `/mcp/tools` to your app.
- A JSON document is sent parsed, YAML as text. No document, or an invalid JSON file: the answer has
  `"tools": []`, the service has no tools.
- Only the gateway's own call is answered. A client request to `/<service>/mcp/tools` comes through
  the gateway with the `soajsinjectobj` header and goes to your app like any other request
  (`MCPTools` used as a handler on its own answers it with 404).

## Environment Variables

- `SOAJS_ENV` (required): Environment code (e.g., "dev", "staging", "production")
- `SOAJS_REGISTRY_API` (required): the gateway's maintenance address as `hostname:port` (e.g., "controller:5000")
- `SOAJS_DEPLOY_MANUAL` (optional): "true" registers the service with the gateway; not set means "false"

Example:

```bash
export SOAJS_ENV=dev
export SOAJS_REGISTRY_API=localhost:5000
export SOAJS_DEPLOY_MANUAL=true
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
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --config .golangci.yml
```

### Code Formatting

```bash
# Format code
gofmt -w .

# Organize imports
goimports -w .
```

## CI/CD

GitHub Actions, like the other SOAJS repositories:

- **CI** (`.github/workflows/ci.yml`), on pushes to `master`, `develop` and `feature/*` and on pull
  requests: golangci-lint v1.64.8 (with Go 1.24) and `go test -race` on Go 1.21, 1.24 and the latest
  stable release.
- **Release** (`.github/workflows/release.yml`), on pushes to `release/*`: the same checks, then the
  tag `v<Version>` from `version.go` and a GitHub release. Go modules resolve versions from these
  tags, so bump `Version` before releasing.

Branches: work on `develop`, merge to `master`, then to `release/v1.x`.

## Documentation

For more information about SOAJS:
- [SOAJS Documentation](https://soajsorg.atlassian.net/wiki/spaces/SOAJ/overview)
- [SOAJS Framework](https://www.soajs.org)

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass
2. Code is properly formatted (gofmt)
3. Linting passes (golangci-lint)
4. Coverage is maintained or improved

## License

*Copyright SOAJS All Rights Reserved.*

Use of this source code is governed by an Apache license that can be found in the LICENSE file at the root of this repository.



