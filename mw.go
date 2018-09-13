package soajsGo

import (
  "log"
  "context"
  "reflect"
  "strings"
  "strconv"
  "net/http"
  "encoding/json"
)

func mapInjectedObject(r *http.Request) SOAJSData {
  soajsHeader := r.Header.Get("soajsinjectobj")

  var input, output SOAJSData
  if inputType := reflect.TypeOf(soajsHeader).String(); inputType == "string" {
    json.Unmarshal([]byte(soajsHeader), &input)
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

  return output
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
    host += ":" + a.Port + "/"
    host += serviceName + "/"

    if _, err := strconv.Atoi(version); err == nil {
      host += "v" + version + "/"
    }
  }

  return host
}

func init() {
  var temp = map[string]string{"envCode": "dashboard", "serviceName": "soajs.urac"}
  go ExecRegistry(temp)
}

func SoajsMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("SOAJS Middleware triggered")

    injectObject := mapInjectedObject(r)

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

    soajs := context.WithValue(r.Context(), "soajs", middlewareOutput)
    r = r.WithContext(soajs)

    next.ServeHTTP(w, r)
  })
}
