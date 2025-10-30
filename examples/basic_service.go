package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	soajsgo "github.com/soajs/soajs.golang"
)

// Basic HTTP example with SOAJS middleware
//
// This example demonstrates how to integrate SOAJS middleware with standard net/http.

func main() {
	ctx := context.Background()

	// Initialize registry manager
	// For manual deployment (SOAJS_DEPLOY_MANUAL=true), use NewFromConfig
	config := soajsgo.Config{
		ServiceName:    "my-go-service",
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

	// Create HTTP handlers
	mux := http.NewServeMux()

	// Root endpoint
	mux.HandleFunc("/", rootHandler)

	// Tenant info endpoint
	mux.HandleFunc("/tenant-info", tenantInfoHandler)

	// Database info endpoint
	mux.HandleFunc("/database-info", makeDatabaseInfoHandler(registry))

	// Services listing endpoint
	mux.HandleFunc("/services", makeServicesHandler(registry))

	// Custom config endpoint
	mux.HandleFunc("/custom-config", makeCustomConfigHandler(registry))

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)

	// Wrap mux with SOAJS middleware
	handler := registry.Middleware(mux)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServicePort),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %d", config.ServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// rootHandler handles the root endpoint
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "SOAJS Go Example",
		"version": "1.0.0",
	}
	writeJSON(w, http.StatusOK, response)
}

// tenantInfoHandler demonstrates accessing SOAJS context data
func tenantInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get SOAJS context from request
	soaData := r.Context().Value(soajsgo.SoajsKey)
	if soaData == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "No SOAJS context available",
		})
		return
	}

	context := soaData.(soajsgo.ContextData)

	response := map[string]interface{}{
		"tenant_id":   context.Tenant.ID,
		"tenant_code": context.Tenant.Code,
		"device":      context.Device,
		"geo":         context.Geo,
	}

	if context.Reg != nil {
		response["environment"] = context.Reg.Environment
	}

	if context.Urac.ID != "" {
		response["user"] = map[string]interface{}{
			"id":       context.Urac.ID,
			"username": context.Urac.Username,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

// makeDatabaseInfoHandler demonstrates accessing database configuration
func makeDatabaseInfoHandler(registry *soajsgo.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get database from registry
		db, err := registry.Database("main")
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		servers := make([]map[string]interface{}, 0, len(db.Server))
		for _, server := range db.Server {
			servers = append(servers, map[string]interface{}{
				"host": server.Host,
				"port": server.Port,
			})
		}

		response := map[string]interface{}{
			"database": db.Name,
			"cluster":  db.Cluster,
			"servers":  servers,
		}

		writeJSON(w, http.StatusOK, response)
	}
}

// makeServicesHandler demonstrates accessing service information
func makeServicesHandler(registry *soajsgo.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services := make(map[string]interface{})

		for name, service := range registry.Services {
			services[name] = map[string]interface{}{
				"group": service.Group,
				"port":  service.Port,
			}
		}

		response := map[string]interface{}{
			"services": services,
		}

		writeJSON(w, http.StatusOK, response)
	}
}

// makeCustomConfigHandler demonstrates accessing custom registry data
func makeCustomConfigHandler(registry *soajsgo.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		custom, err := registry.GetCustom(name)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		var response map[string]interface{}

		if name != "" {
			// Return specific custom registry
			customReg := custom.(*soajsgo.CustomRegistry)
			response = map[string]interface{}{
				"name":   name,
				"custom": customReg,
			}
		} else {
			// Return all custom registries
			customRegistries := custom.(soajsgo.CustomRegistries)
			response = map[string]interface{}{
				"count":   len(customRegistries),
				"customs": customRegistries,
			}
		}

		writeJSON(w, http.StatusOK, response)
	}
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
	}
	writeJSON(w, http.StatusOK, response)
}

// writeJSON writes JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
	}
}
