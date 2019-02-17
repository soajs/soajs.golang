package soajsgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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
func InitMiddleware(config SOA) (func(http.Handler) http.Handler, error) {
	registryAPI := os.Getenv(EnvRegistryAPIAddress)
	soajsEnv := strings.ToLower(os.Getenv(EnvEnv))

	if soajsEnv == "" {
		return nil, fmt.Errorf("could not find environment variable %s", EnvEnv)
	}
	if registryAPI == "" {
		return nil, fmt.Errorf("could not find environment variable %s", EnvRegistryAPIAddress)
	}
	if index := strings.Index(registryAPI, ":"); index == -1 {
		return nil, fmt.Errorf("invalid format for %s. Got [%s], expected [hostname:port]: ", EnvRegistryAPIAddress, registryAPI)
	}
	port := strings.Split(registryAPI, ":")[1]
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("port must be an integer, got %s", port)
	}

	var validVersion = regexp.MustCompile(`[0-9]+(.[0-9]+)?`)

	if config.Type == "" {
		return nil, fmt.Errorf("could not find [type] in your config, type is <required>")
	}
	if config.ServiceName == "" {
		return nil, fmt.Errorf("could not find [ServiceName] in your config, ServiceName is <required>")
	}
	if !(config.ServicePort > 0) {
		return nil, fmt.Errorf("could not find [ServicePort] in your config, ServicePort is <required>")
	}
	if config.ServiceVersion == "" {
		return nil, fmt.Errorf("could not find [ServiceVersion] in your config, ServiceVersion is <required>")
	}
	if !validVersion.MatchString(config.ServiceVersion) {
		return nil, fmt.Errorf("error with [ServiceVersion] in your config, ServiceVersion syntax is [[0-9]+(.[0-9]+)?]")
	}
	//TODO: we should add more assurance for config HERE

	reqURL := fmt.Sprintf("http://%s/getRegistry?env=%s&serviceName=%s", registryAPI, soajsEnv, config.ServiceName)
	reg, err := newRegistry(reqURL, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching registry: %v", err)
	}

	manualDeploySrt := os.Getenv(EnvDeployManual)
	manualDeploy, err := strconv.ParseBool(manualDeploySrt)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s environment variable: %v", EnvDeployManual, err)
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
			return nil, fmt.Errorf("could not marshal manual deploy auto register config: %v", err)
		}
		reqURL := fmt.Sprintf("http://%s/register", registryAPI)
		res, err := http.Post(reqURL, "application/json", bytes.NewBuffer(d))
		if err != nil {
			return nil, fmt.Errorf("could not call %s: %v", reqURL, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			b, _ := ioutil.ReadAll(res.Body)
			return nil, fmt.Errorf("non 2xx status code: %d %v", res.StatusCode, b)
		}
		var temp RegisterAPIResponse
		err = json.NewDecoder(res.Body).Decode(&temp)
		if err != nil || !temp.Result {
			return nil, fmt.Errorf("unable to register service at gateway: %v", err)
		}
	}

	return reg.soajsMiddleware, nil
}

// SoajsMiddleware the middleware that gets triggered per request
func (reg Registry) soajsMiddleware(next http.Handler) http.Handler {
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
