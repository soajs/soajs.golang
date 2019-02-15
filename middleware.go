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

const (
	// The environment variable name that contains the name of the environment where the service is running at.
	EnvEnv = "SOAJS_ENV"

	// The environment variable name that indicates if the service has been deployed manually or not.
	EnvDeployManual = "SOAJS_DEPLOY_MANUAL"

	// The SOAJS Gateway injected object attached to the header of each request between the gateway and teh service
	headerDataName = "soajsinjectobj"

	EnvRegistryAPIAddress = "SOAJS_REGISTRY_API"
)

// InitMiddleware returns http soajs middleware.
func InitMiddleware(config SOA) func(http.Handler) http.Handler {
	registryAPI := os.Getenv(EnvRegistryAPIAddress)
	soajsEnv := strings.ToLower(os.Getenv(EnvEnv))
	if soajsEnv != "" && registryAPI != "" {

		manualDeploySrt := os.Getenv(EnvDeployManual)
		manualDeploy, err := strconv.ParseBool(manualDeploySrt)
		if err != nil {
			panic(fmt.Errorf("could not parse %s envaronment variable: %v", EnvDeployManual, err))
			return SoajsMiddleware
		}
		if manualDeploy {
			var conf RegisterConf

			if config.ServiceIP == "" {
				config.ServiceIP = "127.0.0.1"
			}

			conf.Name = config.ServiceName
			conf.Type = config.Type
			conf.Mw = true
			conf.Group = config.ServiceGroup
			conf.Port = config.ServicePort
			conf.Swagger = config.Swagger
			conf.RequestTimeout = config.RequestTimeout
			conf.RequestTimeoutRenewal = config.RequestTimeoutRenewal
			conf.Version = config.ServiceVersion
			conf.ExtKeyRequired = config.ExtKeyRequired
			conf.Urac = config.Urac
			conf.UracProfile = config.UracProfile
			conf.UracACL = config.UracACL
			conf.ProvisionACL = config.ProvisionACL
			conf.Oauth = config.Oauth
			conf.IP = config.ServiceIP
			conf.Maintenance = config.Maintenance

			d, err := json.Marshal(conf)
			if err != nil {
				panic(fmt.Errorf("could not marshal config: %v", err))
				return SoajsMiddleware
			}
			reqURL := fmt.Sprintf("http://%s/register", registryAPI)
			res, err := http.Post(reqURL, "application/json", bytes.NewBuffer(d))
			if err != nil {
				panic(fmt.Errorf("could not call %s: %v", reqURL, err))
				return SoajsMiddleware
			}
			defer res.Body.Close()
			if res.StatusCode < 200 || res.StatusCode > 299 {
				b, _ := ioutil.ReadAll(res.Body)
				panic(fmt.Errorf("non 2xx status code: %d %v", res.StatusCode, b))
				return SoajsMiddleware
			}
		}
	}
	return SoajsMiddleware
}

// SoajsMiddleware the middleware that gets triggered per request
func SoajsMiddleware(next http.Handler) http.Handler {
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
			//Reg:            reg,
		}
		out.Tenant.Key.IKey = d.Key.IKey
		out.Tenant.Key.EKey = d.Key.EKey

		out.Tenant.Application = d.Application
		out.Tenant.Application.PackageACL = d.Package.ACL
		out.Tenant.Application.PackageACLAllEnv = d.Package.ACLAllEnv

		soajs := context.WithValue(r.Context(), "soajs", out)
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

func (a Host) GetHost(args ...string) string {
	var serviceName, version string
	switch len(args) {
	//controller
	case 1:
		serviceName = args[0]
		//controller, 1
	case 2:
		serviceName = args[0]
		version = args[1]
		//controller, 1, dash [dash is ignored]
	case 3:
		serviceName = args[0]
		version = args[1]
	}

	host := a.Host
	host += ":" + strconv.Itoa(a.Port) + "/"

	if serviceName != "" && strings.ToLower(serviceName) != "controller" {
		host += serviceName + "/"

		if _, err := strconv.Atoi(version); err == nil {
			host += "v" + version + "/"
		}
	}

	return host
}
