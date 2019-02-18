package soajsgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// EnvProfile is the environment variable name that contains location of the profile
	// to use so that SOAJS can connect to the core database.
	EnvProfile = "SOAJS_PROFILE"

	// EnvSRVIP is optional environment variable used to specify which IP address to use
	// if the machine has more than one active interface.
	EnvSRVIP = "SOAJS_SRVIP"

	// EnvSOLO is optional environment variable used to launch any service on top of SOAJS
	// without the need of a database.
	EnvSOLO = "SOAJS_SOLO"

	// EnvSrvAutoRegisterHost is optional environment variable used in case a service should register itself or not.
	EnvSrvAutoRegisterHost = "SOAJS_SRV_AUTOREGISTERHOST"

	// EnvDaemonGRPConf is the environment variable name that contains the name of the daemon group to use;
	// available for daemons ONLY.
	EnvDaemonGRPConf = "SOAJS_DAEMON_GRP_CONF"

	// EnvGCName is the environment variable name that contains mandatory variable if deploying a GCS service
	// and contains the name of that GCS service.
	EnvGCName = "SOAJS_GC_NAME"

	// EnvDCVersion is the environment variable name that contains mandatory variable if deploying a GCS service
	// and contains the version of that GCS service.
	EnvDCVersion = "SOAJS_GC_VERSION"

	// EnvGCMaxUploadLimit is optional variable if deploying a GCS service that specifies the maximum upload limit
	// of file sizes to accept.
	EnvGCMaxUploadLimit = "SOAJS_GC_MAX_UPLOAD_LIMIT"

	// EnvRegistryAPIAddress is the environment variable name that contains the IP address and port of
	// the controller service that runs in the same environment. The SOAJS middleware uses this variable to fetch
	// the registry of this environment and supply it to your service.
	EnvRegistryAPIAddress = "SOAJS_REGISTRY_API"
)

// NewRegistry creates and initialises new registry by service name and code.
// see: https://soajsorg.atlassian.net/wiki/spaces/SOAJ/pages/61347270/Service
// nolint: errcheck
func NewRegistry(serviceName, envCode string) (*Registry, error) {
	if serviceName == "" || envCode == "" {
		return nil, errors.New("service name and env code are required")
	}
	addr := os.Getenv(EnvRegistryAPIAddress)
	if addr == "" {
		return nil, fmt.Errorf("could not find environment variable %s", EnvRegistryAPIAddress)
	}
	if index := strings.Index(addr, ":"); index == -1 {
		return nil, fmt.Errorf("invalid format for %s. Got [%s], expected [hostname:port]: ", EnvRegistryAPIAddress, addr)
	}
	port := strings.Split(addr, ":")[1]
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("port must be an integer, got %s", port)
	}

	reqURL := fmt.Sprintf("http://%s/getRegistry?env=%s&serviceName=%s", addr, envCode, serviceName)
	res, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("could not get registry from api gateway: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("non 2xx status code: %d %s", res.StatusCode, b)
	}
	var temp RegistryAPIResponse
	err = json.NewDecoder(res.Body).Decode(&temp)
	if err != nil {
		return nil, fmt.Errorf("unable to convert registry response: %v", err)
	}
	if len(temp.Errors.Details) > 0 {
		return nil, fmt.Errorf("bad response: [%d] %s", temp.Errors.Details[0].Code, temp.Errors.Details[0].Message)
	}
	return &temp.Registry, nil
}

// Reload does the same that NewRegistry does, It reloads registry from soajs.
func (reg *Registry) Reload() error {
	r, err := NewRegistry(reg.Name, reg.Environment)
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
			if err != nil {
				// TODO: it is not correct to print some logs in library.
				log.Printf("could not reload registry data: %v", err)
			} else {
				ticker = time.NewTicker(reg.ServiceConfig.Awareness.AutoReloadRegistry * time.Millisecond)
			}
		case <-ctx.Done():
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
