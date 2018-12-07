package soajsGo

import (
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
  Tenant                    Tenant                    `json:"tenant"`
  Urac                      Urac                      `json:"urac"`
  ServicesConfig            map[string]interface{}    `json:"servicesConfig"`
  Device                    string                    `json:"device"`
  Geo                       map[string]string         `json:"geo"`
  Awareness                 Awareness                 `json:"awareness"`
  Controller                string                    `json:"controller"`
  Reg                       RegistryObj               `json:"reg"`
}

type JSON = map[string]interface{}
var globalConfig JSON

func mapInjectedObject(r *http.Request) (SOAJSData, error) {
    soajsHeader := r.Header.Get("soajsinjectobj")

    var input, output SOAJSData
    if inputType := reflect.TypeOf(soajsHeader).String(); inputType == "string" {
        if jsonError := json.Unmarshal([]byte(soajsHeader), &input); jsonError != nil {
            return SOAJSData{}, errors.New("Unable to parse SOAJS header object")
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

    if serviceName != "" && strings.ToLower(serviceName) != "controller" {
        host += ":" + strconv.Itoa(a.Port) + "/"
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
        } else {
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

func InitMiddleware(config JSON) (func(http.Handler) http.Handler) {
    globalConfig = config

    serviceName := globalConfig["serviceName"].(string)

    registryApi := os.Getenv("SOAJS_REGISTRY_API")
    soajsEnv := os.Getenv("SOAJS_ENV")
    if soajsEnv != "" && registryApi != "" {
        params := map[string]string{"envCode": strings.ToLower(soajsEnv), "serviceName": serviceName}
        AutoReload(params)
    }

    return SoajsMiddleware
}
