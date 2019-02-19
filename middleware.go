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
)

type (
	key int
)

const (
	// headerDataName is the SOAJS Gateway injected object attached to the header of each request
	// between the gateway and tech service.
	headerDataName = "soajsinjectobj"
	// SoajsKey use this key to init soajs data from context.
	SoajsKey = key(1)
)

// InitMiddleware returns http soajs middleware with registry inside.
// This function starts registry auto reload every AutoReloadRegistry. You can break this process using context.
// nolint: errcheck
func InitMiddleware(ctx context.Context, config Config) (func(http.Handler) http.Handler, error) {
	addr, err := registryAddress()
	if err != nil {
		return nil, fmt.Errorf("could not init registry api path: %v", err)
	}
	soajsEnv := strings.ToLower(os.Getenv(EnvEnv))
	if soajsEnv == "" {
		return nil, fmt.Errorf("could not find environment variable %s", EnvEnv)
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	reg, err := NewRegistry(ctx, config.ServiceName, soajsEnv, true)
	if err != nil {
		return nil, fmt.Errorf("could not fetch registry: %v", err)
	}

	manualDeploySrt := os.Getenv(EnvDeployManual)
	manualDeploy, err := strconv.ParseBool(manualDeploySrt)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s envaronment variable: %v", EnvDeployManual, err)
	}
	if manualDeploy {
		if config.ServiceIP == "" {
			config.ServiceIP = "127.0.0.1"
		}
		regConf := registerConf{
			Name:                  config.ServiceName,
			Type:                  config.Type,
			Middleware:            true,
			Group:                 config.ServiceGroup,
			Port:                  config.ServicePort,
			Swagger:               config.Swagger,
			RequestTimeout:        config.RequestTimeout,
			RequestTimeoutRenewal: config.RequestTimeoutRenewal,
			Version:               config.ServiceVersion,
			ExtKeyRequired:        config.ExtKeyRequired,
			Urac:                  config.Urac,
			UracProfile:           config.UracProfile,
			UracACL:               config.UracACL,
			ProvisionACL:          config.ProvisionACL,
			Oauth:                 config.Oauth,
			IP:                    config.ServiceIP,
			Maintenance:           config.Maintenance,
		}
		d, err := json.Marshal(regConf)
		if err != nil {
			return nil, fmt.Errorf("could not marshal manual deploy auto register config: %v", err)
		}
		res, err := http.Post(addr.register(), "application/json", bytes.NewBuffer(d))
		if err != nil {
			return nil, fmt.Errorf("could not call %s: %v", addr.register(), err)
		}
		defer res.Body.Close()
		_, err = registryResponse(res)
		if err != nil {
			return nil, err
		}
	}

	return reg.Middleware, nil
}

// Middleware is http middleware that gets triggered per request.
func (reg Registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := headerData(r)
		if err != nil {
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

func headerData(r *http.Request) (*headerInfo, error) {
	headerData := r.Header.Get(headerDataName)
	if headerData == "" {
		return nil, nil
	}
	d := new(headerInfo)
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
		host = fmt.Sprintf("%s%s/", host, serviceName)
		if _, err := strconv.Atoi(version); err == nil {
			host = fmt.Sprintf("%sv%s/", host, version)
		}
	}
	return host
}
