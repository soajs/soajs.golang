package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	soajsgo "github.com/soajs/soajs.golang"
)

// Gin framework example with SOAJS middleware
//
// This example demonstrates how to integrate SOAJS middleware with Gin framework.
//
// To run this example:
// 1. Install Gin: go get -u github.com/gin-gonic/gin
// 2. Set environment variables:
//    export SOAJS_ENV=dev
//    export SOAJS_REGISTRY_API=http://localhost:5000
//    export SOAJS_DEPLOY_MANUAL=true
// 3. Run: go run gin_service.go

func main() {
	ctx := context.Background()

	// Initialize registry manager
	config := soajsgo.Config{
		ServiceName:    "my-gin-service",
		ServiceGroup:   "my-group",
		ServicePort:    8080,
		ServiceIP:      "127.0.0.1",
		Type:           "service",
		ServiceVersion: "1",
	}

	registry, err := soajsgo.NewFromConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Add Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add SOAJS middleware wrapper for Gin
	router.Use(soajsMiddleware(registry))

	// Define routes
	router.GET("/", rootHandler)
	router.GET("/tenant-info", tenantInfoHandler)
	router.GET("/database-info", makeDatabaseInfoHandler(registry))
	router.GET("/services", makeServicesHandler(registry))
	router.GET("/custom-config", makeCustomConfigHandler(registry))
	router.GET("/health", healthHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting Gin server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// soajsMiddleware wraps the SOAJS middleware for Gin
func soajsMiddleware(registry *soajsgo.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a wrapped handler that calls the next middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Update the Gin context with the new request context
			c.Request = r
			c.Next()
		})

		// Apply SOAJS middleware
		soajsHandler := registry.Middleware(handler)
		soajsHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// rootHandler handles the root endpoint
func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "SOAJS Gin Example",
		"version": "1.0.0",
	})
}

// tenantInfoHandler demonstrates accessing SOAJS context data
func tenantInfoHandler(c *gin.Context) {
	// Get SOAJS context from request
	soaData := c.Request.Context().Value(soajsgo.SoajsKey)
	if soaData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No SOAJS context available",
		})
		return
	}

	context := soaData.(soajsgo.ContextData)

	response := gin.H{
		"tenant_id":   context.Tenant.ID,
		"tenant_code": context.Tenant.Code,
		"device":      context.Device,
		"geo":         context.Geo,
	}

	if context.Reg != nil {
		response["environment"] = context.Reg.Environment
	}

	if context.Urac.ID != "" {
		response["user"] = gin.H{
			"id":       context.Urac.ID,
			"username": context.Urac.Username,
			"email":    context.Urac.Email,
		}
	}

	c.JSON(http.StatusOK, response)
}

// makeDatabaseInfoHandler demonstrates accessing database configuration
func makeDatabaseInfoHandler(registry *soajsgo.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get database from registry
		db, err := registry.Database("main")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		servers := make([]gin.H, 0, len(db.Server))
		for _, server := range db.Server {
			servers = append(servers, gin.H{
				"host": server.Host,
				"port": server.Port,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"database": db.Name,
			"cluster":  db.Cluster,
			"prefix":   db.Prefix,
			"servers":  servers,
		})
	}
}

// makeServicesHandler demonstrates accessing service information
func makeServicesHandler(registry *soajsgo.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		services := make(map[string]interface{})

		for name, service := range registry.Services {
			services[name] = gin.H{
				"group": service.Group,
				"port":  service.Port,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    len(services),
			"services": services,
		})
	}
}

// makeCustomConfigHandler demonstrates accessing custom registry data
func makeCustomConfigHandler(registry *soajsgo.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		custom, err := registry.GetCustom(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if name != "" {
			// Return specific custom registry
			customReg := custom.(*soajsgo.CustomRegistry)
			c.JSON(http.StatusOK, gin.H{
				"name":   name,
				"custom": gin.H{
					"id":      customReg.ID,
					"name":    customReg.Name,
					"locked":  customReg.Locked,
					"plugged": customReg.Plugged,
					"shared":  customReg.Shared,
					"value":   customReg.Value,
					"author":  customReg.Author,
				},
			})
		} else {
			// Return all custom registries
			customRegistries := custom.(soajsgo.CustomRegistries)
			c.JSON(http.StatusOK, gin.H{
				"count":   len(customRegistries),
				"customs": customRegistries,
			})
		}
	}
}

// healthHandler handles health check requests
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
