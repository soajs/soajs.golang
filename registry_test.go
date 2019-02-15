package soajsgo

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistry(t *testing.T) {
	tt := []struct {
		name             string
		serviceName      string
		envCode          string
		envRegAPI        string
		expectedRegistry *Registry
		expectedError    error
	}{
		{
			name:             "empty arguments",
			serviceName:      "",
			envCode:          "",
			envRegAPI:        "",
			expectedRegistry: nil,
			expectedError:    errors.New("service name and env code are required"),
		},
		{
			name:             "empty environment",
			serviceName:      "test",
			envCode:          "test",
			envRegAPI:        "",
			expectedRegistry: nil,
			expectedError:    fmt.Errorf("could not find environment variable %s", EnvRegistryAPIAddress),
		},
		{
			name:             "bad api path",
			serviceName:      "test",
			envCode:          "test",
			envRegAPI:        "localhost",
			expectedRegistry: nil,
			expectedError:    fmt.Errorf("invalid format for %s. Got [localhost], expected [hostname:port]: ", EnvRegistryAPIAddress),
		},
		{
			name:             "bad api path port",
			serviceName:      "test",
			envCode:          "test",
			envRegAPI:        "localhost:test",
			expectedRegistry: nil,
			expectedError:    errors.New("port must be an integer, got test"),
		},
		{
			name:             "bad api call",
			serviceName:      "test",
			envCode:          "test",
			envRegAPI:        "localhost:123",
			expectedRegistry: nil,
			expectedError:    errors.New("could not get registry from api gateway: Get http://localhost:123/getRegistry?env=test&serviceName=test: dial tcp [::1]:123: connect: connection refused"),
		},
	}
	lastEnvRegAPI := os.Getenv(EnvRegistryAPIAddress)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(EnvRegistryAPIAddress, tc.envRegAPI))
			reg, err := NewRegistry(tc.serviceName, tc.envCode)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedRegistry, reg)
			assert.NoError(t, os.Setenv(EnvRegistryAPIAddress, lastEnvRegAPI))
		})
	}
}

func TestRegistry_Reload(t *testing.T) {
	reg := Registry{}
	assert.Error(t, reg.Reload())
}

func TestRegistry_Database(t *testing.T) {
	tt := []struct {
		name             string
		dbName           string
		reg              Registry
		expectedDatabase *Database
		expectedErr      error
	}{
		{
			name:             "empty db name",
			dbName:           "",
			reg:              Registry{},
			expectedDatabase: nil,
			expectedErr:      errors.New("database name is required"),
		},
		{
			name:   "core dbs",
			dbName: "core",
			reg: Registry{
				CoreDBs: map[string]Database{"core": {
					Name: "core database",
				}},
			},
			expectedDatabase: &Database{
				Name: "core database",
			},
			expectedErr: nil,
		},
		{
			name:   "meta dbs",
			dbName: "meta",
			reg: Registry{
				TenantMetaDBs: map[string]Database{"meta": {
					Name: "meta database",
				}},
			},
			expectedDatabase: &Database{
				Name: "meta database",
			},
			expectedErr: nil,
		},
		{
			name:             "no dbs",
			dbName:           "empty",
			reg:              Registry{},
			expectedDatabase: nil,
			expectedErr:      errors.New("could not found database"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, err := tc.reg.Database(tc.dbName)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedDatabase, db)
		})
	}
}

func TestRegistry_Databases(t *testing.T) {
	tt := []struct {
		name              string
		reg               Registry
		expectedDatabases map[string]Database
		expectedErr       error
	}{
		{
			name:              "not found",
			reg:               Registry{},
			expectedDatabases: nil,
			expectedErr:       errors.New("could not found databases"),
		},
		{
			name: "found",
			reg: Registry{
				CoreDBs:       map[string]Database{"core": {Name: "core"}},
				TenantMetaDBs: map[string]Database{"meta": {Name: "meta"}},
			},
			expectedDatabases: map[string]Database{"core": {Name: "core"}, "meta": {Name: "meta"}},
			expectedErr:       nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dbs, err := tc.reg.Databases()
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedDatabases, dbs)
		})
	}
}

func TestRegistry_Resource(t *testing.T) {
	tt := []struct {
		name             string
		resourceName     string
		reg              Registry
		expectedResource *Resource
		expectedErr      error
	}{
		{
			name:             "empty name",
			reg:              Registry{},
			expectedResource: nil,
			expectedErr:      errors.New("resource name is required"),
		},
		{
			name:         "found",
			resourceName: "good",
			reg: Registry{
				Resources: Resources{"0": map[string]Resource{
					"bad":  {Name: "bad"},
					"good": {Name: "good"},
				}},
			},
			expectedResource: &Resource{Name: "good"},
			expectedErr:      nil,
		},
		{
			name:         "not found",
			resourceName: "good",
			reg: Registry{
				Resources: Resources{"0": map[string]Resource{
					"bad": {Name: "bad"},
				}},
			},
			expectedResource: nil,
			expectedErr:      errors.New("resource not found"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.reg.Resource(tc.resourceName)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResource, res)
		})
	}
}

func TestRegistry_Service(t *testing.T) {
	tt := []struct {
		name            string
		serviceName     string
		reg             Registry
		expectedService *Service
		expectedErr     error
	}{
		{
			name:            "empty name",
			serviceName:     "",
			reg:             Registry{},
			expectedService: nil,
			expectedErr:     errors.New("service name is required"),
		},
		{
			name:        "found",
			serviceName: "good",
			reg: Registry{
				Services: map[string]Service{
					"bad":  {Port: 1},
					"good": {Port: 2},
				},
			},
			expectedService: &Service{Port: 2},
			expectedErr:     nil,
		},
		{
			name:        "not found",
			serviceName: "good",
			reg: Registry{
				Services: map[string]Service{
					"bad": {Port: 1},
				},
			},
			expectedService: nil,
			expectedErr:     errors.New("service not found"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.reg.Service(tc.serviceName)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedService, s)
		})
	}
}
