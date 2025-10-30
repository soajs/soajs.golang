# SOAJS Go Middleware Examples

This directory contains example microservices demonstrating how to use the SOAJS Go middleware with different frameworks.

## Examples

- **basic_service.go** - Basic HTTP service using standard `net/http`
- **gin_service.go** - Service using the Gin Web Framework

## Prerequisites

1. **Go 1.21 or higher** installed
2. **SOAJS infrastructure** (Controller, Registry API) running
3. **Required Go packages**:
   - For Gin example: `go get -u github.com/gin-gonic/gin`

## Environment Variables

Before running any example, set the following environment variables:

```bash
# Environment code (dev, staging, production, etc.)
export SOAJS_ENV=dev

# Registry API endpoint
export SOAJS_REGISTRY_API=http://localhost:5000

# Manual deployment flag
export SOAJS_DEPLOY_MANUAL=true
```

## Running the Examples

### Basic HTTP Service

The basic service uses standard Go `net/http` library:

```bash
# Navigate to examples directory
cd examples

# Run the service
go run basic_service.go
```

The service will start on port 8080.

### Gin Service

The Gin example uses the popular Gin web framework:

```bash
# Install Gin (if not already installed)
go get -u github.com/gin-gonic/gin

# Run the service
go run gin_service.go
```

The service will start on port 8080.

## Available Endpoints

All examples expose the following endpoints:

### `GET /`
Root endpoint returning service information.

**Response:**
```json
{
  "message": "SOAJS Go Example",
  "version": "1.0.0"
}
```

### `GET /tenant-info`
Returns tenant information from SOAJS context.

**Headers Required:**
- `soajsinjectobj` - SOAJS gateway injected object (automatically added by SOAJS Gateway)

**Response:**
```json
{
  "tenant_id": "5551aca9e179c39b760f7a1a",
  "tenant_code": "DBTN",
  "environment": "dev",
  "device": "iPhone",
  "geo": {
    "country": "US"
  },
  "user": {
    "id": "12345",
    "username": "john.doe"
  }
}
```

### `GET /database-info`
Returns database configuration from registry.

**Response:**
```json
{
  "database": "main",
  "cluster": "cluster1",
  "servers": [
    {
      "host": "localhost",
      "port": 27017
    }
  ]
}
```

### `GET /services`
Lists all available services from the registry.

**Response:**
```json
{
  "services": {
    "controller": {
      "group": "core",
      "port": 4000,
      "host": "localhost"
    },
    "my-service": {
      "group": "my-group",
      "port": 8080,
      "host": "localhost"
    }
  }
}
```

### `GET /custom-config?name=<name>`
Returns custom registry configuration.

**Query Parameters:**
- `name` (optional) - Name of specific custom registry to retrieve. If omitted, returns all custom registries.

**Response (specific custom):**
```json
{
  "name": "myCustom",
  "custom": {
    "id": "abc123",
    "name": "myCustom",
    "locked": true,
    "plugged": false,
    "shared": true,
    "value": {
      "key": "value"
    },
    "author": "admin"
  }
}
```

**Response (all customs):**
```json
{
  "count": 2,
  "customs": {
    "custom1": {
      "id": "abc123",
      "name": "custom1",
      "value": "value1"
    },
    "custom2": {
      "id": "def456",
      "name": "custom2",
      "value": "value2"
    }
  }
}
```

### `GET /health`
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Testing the Examples

### Using curl

1. Start the service:
```bash
go run basic_service.go
```

2. Test the root endpoint:
```bash
curl http://localhost:8080/
```

3. Test health endpoint:
```bash
curl http://localhost:8080/health
```

4. Test tenant info (requires SOAJS gateway):
```bash
curl http://localhost:8080/tenant-info \
  -H "soajsinjectobj: {\"tenant\":{\"id\":\"123\",\"code\":\"TEST\"},\"device\":\"curl\"}"
```

5. Test custom config:
```bash
# Get all custom registries
curl http://localhost:8080/custom-config

# Get specific custom registry
curl http://localhost:8080/custom-config?name=myCustom
```

### Behind SOAJS Gateway

When running behind the SOAJS Gateway, the gateway automatically injects the `soajsinjectobj` header with tenant, user, and other context information.

```bash
# Call through gateway
curl http://gateway-host:gateway-port/my-go-service/tenant-info
```

## Service Configuration

Both examples use the same service configuration pattern:

```go
config := soajsgo.Config{
    ServiceName:    "my-go-service",      // Service name
    ServiceGroup:   "my-group",           // Service group
    ServicePort:    8080,                 // Service port
    ServiceIP:      "127.0.0.1",          // Service IP
    Type:           "service",            // Type: "service" or "daemon"
    ServiceVersion: "1",                  // Service version
}

registry, err := soajsgo.NewFromConfig(ctx, config)
if err != nil {
    log.Fatal(err)
}
```

## Middleware Integration

### Standard HTTP

```go
// Create your HTTP handler
mux := http.NewServeMux()
mux.HandleFunc("/", handler)

// Wrap with SOAJS middleware
wrappedHandler := registry.Middleware(mux)

// Start server
http.ListenAndServe(":8080", wrappedHandler)
```

### Gin Framework

```go
// Create Gin router
router := gin.New()

// Add SOAJS middleware
router.Use(soajsMiddleware(registry))

// Add routes
router.GET("/", handler)
```

## Accessing SOAJS Context

In your handlers, access the SOAJS context data:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Get SOAJS context
    soaData := r.Context().Value(soajsgo.SoajsKey)
    if soaData == nil {
        // No context available
        return
    }

    context := soaData.(soajsgo.ContextData)

    // Access tenant information
    tenantCode := context.Tenant.Code
    tenantID := context.Tenant.ID

    // Access user information (if authenticated)
    if context.Urac != nil {
        username := context.Urac.Username
    }

    // Access registry
    registry := context.Reg
    db, _ := registry.Database("mydb")
}
```

## Production Deployment

For production deployments:

1. Set `SOAJS_DEPLOY_MANUAL=false` to use automatic service discovery
2. Configure proper logging and monitoring
3. Use environment-specific configuration
4. Implement graceful shutdown (both examples include this)
5. Add proper error handling and recovery
6. Configure timeouts appropriately

## Troubleshooting

### "Failed to initialize registry"
- Ensure SOAJS_ENV is set correctly
- Verify SOAJS_REGISTRY_API points to your controller
- Check network connectivity to the registry API

### "No SOAJS context available"
- Ensure the request includes the `soajsinjectobj` header
- Verify the middleware is properly wrapped around your handlers
- When testing directly (not through gateway), manually add the header

### Registry data not updating
- Check the auto-reload configuration
- Verify network connectivity to the registry API
- Check registry API logs for errors

## Additional Resources

- [SOAJS Documentation](https://soajsorg.atlassian.net/wiki/spaces/SOAJ/overview)
- [SOAJS Framework](https://www.soajs.org)
- [Go Middleware Documentation](../README.md)

## Contributing

Contributions are welcome! If you have an example using a different framework or pattern, please submit a pull request.
