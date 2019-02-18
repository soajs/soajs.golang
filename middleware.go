package soajsgo

import (
	"bytes"
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
)

type (
	key int
)

const (
	// EnvEnv is the environment variable name that contains the service uses and registers itself under it upon start.
	EnvEnv = "SOAJS_ENV"
	// EnvDeployManual is the environment variable name that contains boolean parameter if it deploys manual.
	EnvDeployManual = "SOAJS_DEPLOY_MANUAL"

	headerDataName = "soajsinjectobj"

	// SoajsKey use this key to get soajs data from context.
	SoajsKey = key(1)
)

// InitMiddleware returns http middleware with registry inside.
// This function starts registry auto reload every AutoReloadRegistry. You can break this process using context.
// nolint: errcheck
func InitMiddleware(ctx context.Context, config Config) (func(http.Handler) http.Handler, error) {
	registryAPI := os.Getenv(EnvRegistryAPIAddress)
	soajsEnv := strings.ToLower(os.Getenv(EnvEnv))
	if soajsEnv == "" || registryAPI == "" {
		reg := Registry{}
		return reg.Middleware, nil
	}
	reg, err := NewRegistry(config.ServiceName, soajsEnv)
	if err != nil {
		return nil, fmt.Errorf("could not create new registry: %v", err)
	}
	go reg.autoReload(ctx)

	manualDeploySrt := os.Getenv(EnvDeployManual)
	manualDeploy, err := strconv.ParseBool(manualDeploySrt)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s envaronment variable: %v", EnvDeployManual, err)
	}
	if manualDeploy {
		if config.ServiceIP == "" {
			config.ServiceIP = "127.0.0.1"
		}
		config.Type = "service"
		config.Middleware = true
		config.BodyParser = true
		config.MethodOverride = true
		config.CookieParser = true
		config.Logger = true
		// TODO: do we really need to set and keep it?
		config.Group = config.ServiceGroup
		config.Name = config.ServiceName
		config.Port = config.ServicePort
		config.Version = config.ServiceVersion
		d, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("could not marshal config: %v", err)
		}
		reqURL := fmt.Sprintf("http://%s/register", registryAPI)
		res, err := http.Post(reqURL, "application/json", bytes.NewBuffer(d))
		if err != nil {
			return nil, fmt.Errorf("could not call %s: %v", reqURL, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			b, _ := ioutil.ReadAll(res.Body)
			return nil, fmt.Errorf("non 2xx status code: %d %s", res.StatusCode, b)
		}
	}

	return reg.Middleware, nil
}

// Middleware is http middleware.
func (reg Registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := headerData(r)
		if err != nil {
			log.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		if d == nil {
			next.ServeHTTP(w, r)
			return
		}
		out := ContextData{
			Tenant:         d.Tenant,
			Urac:           d.Urac,
			ServicesConfig: d.Key.Config,
			Device:         d.Device,
			Geo:            d.Geo,
			Awareness:      d.Awareness,
			Reg:            reg,
		}
		out.Tenant.Key.IKey = d.Key.IKey
		out.Tenant.Key.EKey = d.Key.EKey

		out.Tenant.Application = d.Application
		out.Tenant.Application.PackageACL = d.Package.ACL
		out.Tenant.Application.PackageACLAllEnv = d.Package.ACLAllEnv

		soajs := context.WithValue(r.Context(), SoajsKey, out)
		next.ServeHTTP(w, r.WithContext(soajs))
	})
}

func headerData(r *http.Request) (*HeaderInfo, error) {
	headerData := r.Header.Get(headerDataName)
	if headerData == "" {
		return nil, nil
	}
	d := new(HeaderInfo)
	if err := json.Unmarshal([]byte(headerData), d); err != nil {
		return nil, errors.New("unable to parse SOAJS header")
	}
	return d, nil
}

// Path returns compiled service path.
func (a Host) Path(args ...string) string {
	var serviceName, version string
	switch len(args) {
	// controller
	case 1:
		serviceName = args[0]
		// controller, 1
	case 2:
		serviceName = args[0]
		version = args[1]
		// controller, 1, dash [dash is ignored]
	case 3:
		serviceName = args[0]
		version = args[1]
	}
	host := fmt.Sprintf("%s:%d/", a.Host, a.Port)
	if strings.EqualFold(serviceName, "controller") {
		host += serviceName + "/"
		if _, err := strconv.Atoi(version); err == nil {
			host += "v" + version + "/"
		}
	}
	return host
}
