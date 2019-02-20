package soajsgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// NewRegistry creates and initializes new registry by service name and code.
// see: https://soajsorg.atlassian.net/wiki/spaces/SOAJ/pages/61347270/Service
// nolint: errcheck
func NewRegistry(ctx context.Context, serviceName, envCode string, turnOnAutoReload bool) (*Registry, error) {
	if serviceName == "" || envCode == "" {
		return nil, errors.New("service name and env code are required")
	}
	addr, err := registryAddress()
	if err != nil {
		return nil, fmt.Errorf("could not init registry api path: %v", err)
	}
	res, err := http.Get(addr.getRegistry(serviceName, envCode))
	if err != nil {
		return nil, fmt.Errorf("could not init registry from api gateway: %v", err)
	}
	defer res.Body.Close()
	reg, err := registryResponse(res)
	if err != nil {
		return nil, err
	}
	if turnOnAutoReload {
		go reg.autoReload(ctx)
	}
	return reg, nil
}

// Reload does the same that NewRegistry does, It reloads registry from soajs.
func (reg *Registry) Reload() error {
	r, err := NewRegistry(context.Background(), reg.Name, reg.Environment, false)
	if err != nil {
		return err
	}
	// TODO: potential concurrency problems here.
	*reg = *r
	return nil
}

// You can run this method in go routine.
func (reg *Registry) autoReload(ctx context.Context) {
	ticker := time.NewTicker(reg.ServiceConfig.Awareness.AutoReloadRegistry * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := reg.Reload()
			if err == nil {
				ticker.Stop()
				ticker = time.NewTicker(reg.ServiceConfig.Awareness.AutoReloadRegistry * time.Millisecond)
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// Database returns one database by name.
func (reg Registry) Database(dbName string) (*Database, error) {
	if dbName == "" {
		return nil, errors.New("database name is required")
	}
	if db, ok := reg.CoreDBs[dbName]; ok {
		return &db, nil
	}
	if db, ok := reg.TenantMetaDBs[dbName]; ok {
		return &db, nil
	}
	return nil, errors.New("could not found database")
}

// Databases returns all databases.
func (reg Registry) Databases() (map[string]Database, error) {
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
func (reg Registry) Resource(name string) (*Resource, error) {
	if name == "" {
		return nil, errors.New("resource name is required")
	}
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
func (reg Registry) Service(name string) (*Service, error) {
	if name == "" {
		return nil, errors.New("service name is required")
	}
	if s, ok := reg.Services[name]; ok {
		return &s, nil
	}
	return nil, errors.New("service not found")
}
