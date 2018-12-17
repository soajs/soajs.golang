package soajsGo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type SOAJSObject struct {
	Tenant         Tenant                 `json:"tenant"`
	Urac           Urac                   `json:"urac"`
	ServicesConfig map[string]interface{} `json:"servicesConfig"`
	Device         string                 `json:"device"`
	Geo            map[string]string      `json:"geo"`
	Awareness      Awareness              `json:"awareness"`
	Controller     string                 `json:"controller"`
	Reg            RegistryObj            `json:"reg"`
}

type JSON = map[string]interface{}

var globalConfig JSON

func mapInjectedObject(r *http.Request) (SOAJSData, error) {
	soajsHeader := r.Header.Get("soajsinjectobj")

	if soajsHeader == "" {
		return SOAJSData{}, nil
	}
	var input, output SOAJSData
	if inputType := reflect.TypeOf(soajsHeader).String(); inputType == "string" {
		if jsonError := json.Unmarshal([]byte(soajsHeader), &input); jsonError != nil {
			return SOAJSData{}, errors.New("unable to parse SOAJS header object")
		}
	}

	// map information to output
	output.Tenant = input.Tenant
	output.Key = input.Key
	output.Application = input.Application
	output.Package = input.Package
	output.Device = input.Device
	output.Geo = input.Geo
	output.Urac = input.Urac
	output.Awareness = input.Awareness

	return output, nil
}

func (a Awareness) GetHost(args ...string) string {
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

func SoajsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		injectObject, err := mapInjectedObject(r)
		if err != nil {
			log.Println(err)
		} else if !injectObject.IsEmpty() {
			middlewareOutput := SOAJSObject{}
			middlewareOutput.Tenant = injectObject.Tenant

			middlewareOutput.Tenant.Key.IKey = injectObject.Key.IKey
			middlewareOutput.Tenant.Key.EKey = injectObject.Key.EKey

			middlewareOutput.Tenant.Application = injectObject.Application
			middlewareOutput.Tenant.Application.Package_acl = injectObject.Package.Acl
			middlewareOutput.Tenant.Application.Package_acl_all_env = injectObject.Package.Acl_all_env

			middlewareOutput.Urac = injectObject.Urac
			middlewareOutput.ServicesConfig = injectObject.Key.Config
			middlewareOutput.Device = injectObject.Device
			middlewareOutput.Geo = injectObject.Geo
			middlewareOutput.Awareness = injectObject.Awareness

			if os.Getenv("SOAJS_REGISTRY_API") != "" && os.Getenv("SOAJS_ENV") != "" {
				middlewareOutput.Reg = regObj
			}

			soajs := context.WithValue(r.Context(), "soajs", middlewareOutput)
			r = r.WithContext(soajs)
		}
		next.ServeHTTP(w, r)
	})
}

func InitMiddleware(config JSON) func(http.Handler) http.Handler {
	globalConfig = config

	serviceName := globalConfig["serviceName"].(string)

	registryApi := os.Getenv("SOAJS_REGISTRY_API")
	soajsEnv := os.Getenv("SOAJS_ENV")
	if soajsEnv != "" && registryApi != "" {
		params := map[string]string{"envCode": strings.ToLower(soajsEnv), "serviceName": serviceName}
		AutoReload(params)

		manualDeploy := os.Getenv("SOAJS_DEPLOY_MANUAL")
		if manualDeploy == "1" {
			mwIP := "127.0.0.1"
			if globalConfig["IP"] != nil {
				mwIP = globalConfig["IP"].(string)
			}
			reqUrl := "http://" + registryApi + "/register"
			globalConfig["ip"] = mwIP
			globalConfig["type"] = "service"
			globalConfig["mw"] = true
			globalConfig["group"] = globalConfig["serviceGroup"]
			globalConfig["name"] = globalConfig["serviceName"]
			globalConfig["port"] = globalConfig["servicePort"]
			globalConfig["version"] = globalConfig["serviceVersion"]
			bytesRepresentation, err := json.Marshal(globalConfig)
			if err != nil {
				log.Fatalln(err)
			}
			httpResponse, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
			if err != nil {
				log.Println("unable to get registry host @ gateway")
				log.Println(err)
			} else {
				defer httpResponse.Body.Close()
			}
		}
	}

	return SoajsMiddleware
}
