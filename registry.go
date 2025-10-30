package soajsgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// New creates and initializes new registry by service name and code.
// This function starts registry auto reload every AutoReloadRegistry if turnOnAutoReload set as true. You can break
// this process using context.
// nolint: errcheck
func New(ctx context.Context, serviceName, envCode, serviceType string, turnOnAutoReload bool) (*Registry, error) {
	if serviceName == "" || envCode == "" {
		return nil, errors.New("service name and env code are required")
	}
	addr, err := registryAddress()
	if err != nil {
		return nil, fmt.Errorf("could not init registry api path: %v", err)
	}
	res, err := httpClient.Get(addr.getRegistry(serviceName, envCode, serviceType))
	if err != nil {
		return nil, fmt.Errorf("could not init registry from api gateway: %v", err)
	}
	defer res.Body.Close()
	reg, err := registryResponse(res)
	if err != nil {
		return nil, err
	}
	reg.ServiceType = serviceType
	if turnOnAutoReload {
		go reg.autoReload(ctx)
	}
	return reg, nil
}

// NewFromConfig creates and initializes new registry by the configuration.
// This function starts registry auto reload every AutoReloadRegistry. You can break this process using context.
func NewFromConfig(ctx context.Context, config Config) (*Registry, error) {
	addr, err := registryAddress()
	if err != nil {
		return nil, fmt.Errorf("could not init registry api path: %v", err)
	}
	soajsEnv := strings.ToLower(os.Getenv(EnvSoajsEnv))
	if soajsEnv == "" {
		return nil, fmt.Errorf("could not find environment variable %s", EnvSoajsEnv)
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	reg, err := New(ctx, config.ServiceName, soajsEnv, config.Type, true)
	if err != nil {
		return nil, fmt.Errorf("could not fetch registry: %v", err)
	}
	err = manualDeploy(config, addr)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

// nolint: errcheck
func manualDeploy(config Config, addr registryPath) error {
	manualDeploySrt := os.Getenv(EnvDeployManual)
	manualDeploy, err := strconv.ParseBool(manualDeploySrt)
	if err != nil {
		return fmt.Errorf("could not parse %s environment variable: %v", EnvDeployManual, err)
	}
	if manualDeploy {
		if config.ServiceIP == "" {
			config.ServiceIP = "127.0.0.1"
		}
		regConf := registerConf{
			Name:                  config.ServiceName,
			Group:                 config.ServiceGroup,
			Port:                  config.ServicePort,
			IP:                    config.ServiceIP,
			Type:                  config.Type,
			Version:               config.ServiceVersion,
			SubType:               config.SubType,
			Description:           config.Description,
			Oauth:                 config.Oauth,
			Urac:                  config.Urac,
			UracProfile:           config.UracProfile,
			UracACL:               config.UracACL,
			UracConfig:            config.UracConfig,
			UracGroupConfig:       config.UracGroupConfig,
			TenantProfile:         config.TenantProfile,
			ProvisionACL:          config.ProvisionACL,
			RequestTimeout:        config.RequestTimeout,
			RequestTimeoutRenewal: config.RequestTimeoutRenewal,
			Middleware:            true,
			ExtKeyRequired:        config.ExtKeyRequired,
			Maintenance:           config.Maintenance,
			InterConnect:          config.InterConnect,
		}
		d, err := json.Marshal(regConf)
		if err != nil {
			return fmt.Errorf("could not marshal manual deploy auto register config: %v", err)
		}
		res, err := httpClient.Post(addr.register(), "application/json", bytes.NewBuffer(d))
		if err != nil {
			return fmt.Errorf("could not call %s: %v", addr.register(), err)
		}
		defer res.Body.Close()
		_, err = registryResponse(res)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reload does the same that New does, It reloads registry from soajs.
func (reg *Registry) Reload() error {
	r, err := New(context.Background(), reg.Name, reg.Environment, reg.ServiceType, false)
	if err != nil {
		return err
	}
	// Thread-safe update of registry data
	reg.mu.Lock()
	defer reg.mu.Unlock()

	reg.TimeLoaded = r.TimeLoaded
	reg.Name = r.Name
	reg.Environment = r.Environment
	reg.ServiceType = r.ServiceType
	reg.CoreDBs = r.CoreDBs
	reg.TenantMetaDBs = r.TenantMetaDBs
	reg.ServiceConfig = r.ServiceConfig
	reg.Custom = r.Custom
	reg.Resources = r.Resources
	reg.Services = r.Services

	return nil
}

// You can run this method in go routine.
func (reg *Registry) autoReload(ctx context.Context) {
	ticker := time.NewTicker(reg.autoReloadDuration())
	for {
		select {
		case <-ticker.C:
			err := reg.Reload()
			if err == nil {
				ticker.Stop()
				ticker = time.NewTicker(reg.autoReloadDuration())
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (reg *Registry) autoReloadDuration() time.Duration {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if reg.ServiceConfig.Awareness.AutoReloadRegistry > 0 {
		duration := reg.ServiceConfig.Awareness.AutoReloadRegistry * time.Millisecond
		// Ensure minimum reload interval of 1 second to prevent performance issues
		if duration < time.Second {
			return time.Second
		}
		return duration
	}
	return time.Hour
}

// Database returns one database by name.
func (reg *Registry) Database(dbName string) (*Database, error) {
	if dbName == "" {
		return nil, errors.New("database name is required")
	}
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if db, ok := reg.CoreDBs[dbName]; ok {
		return &db, nil
	}
	if db, ok := reg.TenantMetaDBs[dbName]; ok {
		return &db, nil
	}
	return nil, errors.New("could not found database")
}

// Databases returns all databases.
func (reg *Registry) Databases() (map[string]Database, error) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	dbs := make(map[string]Database, len(reg.CoreDBs)+len(reg.TenantMetaDBs))
	for dbName := range reg.CoreDBs {
		dbs[dbName] = reg.CoreDBs[dbName]
	}
	for dbName := range reg.TenantMetaDBs {
		dbs[dbName] = reg.TenantMetaDBs[dbName]
	}
	if len(dbs) > 0 {
		return dbs, nil
	}
	return nil, errors.New("could not found databases")
}

// Resource returns one resource.
func (reg *Registry) Resource(name string) (*Resource, error) {
	if name == "" {
		return nil, errors.New("resource name is required")
	}
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	for _, resourceList := range reg.Resources {
		for resourceKey, resourceData := range resourceList {
			if resourceKey == name {
				return &resourceData, nil
			}
		}
	}
	return nil, errors.New("resource not found")
}

// Service returns one service by name.
func (reg *Registry) Service(name string) (*Service, error) {
	if name == "" {
		return nil, errors.New("service name is required")
	}
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if s, ok := reg.Services[name]; ok {
		return &s, nil
	}
	return nil, errors.New("service not found")
}

// GetCustom returns one custom registry by name. If name is empty, returns all custom registries.
func (reg *Registry) GetCustom(name string) (interface{}, error) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if name != "" {
		if custom, ok := reg.Custom[name]; ok {
			return &custom, nil
		}
		return nil, errors.New("custom registry not found")
	}
	if len(reg.Custom) > 0 {
		return reg.Custom, nil
	}
	return nil, errors.New("no custom registries found")
}
