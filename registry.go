package soajsGo

import (
  "os"
  "log"
  "time"
  "errors"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "github.com/soajs/soajs.golang/registry/structs"
)

type RegistryObj struct {
  Env                       string                    `json:"env"`
  ServiceName               string                    `json:"serviceName"`
}

type RegistryApiResponse struct {
  Result                    bool                      `json:"result"`
  Ts                        int64                     `json:"ts"`
  Service                   map[string]string         `json:"service"`
  Data                      structs.Registry          `json:"data"`
}

var registry_struct map[string]structs.Registry
var regObj RegistryObj

// func (reg *RegistryObj) GetDatabase(dbName string) (nterface) {
//   return db
// }

/**
 * Get all databases
 *
 * @return {Databases}
 */
func (reg *RegistryObj) GetDatabases() (structs.Databases, error) {
  var databases structs.Databases
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return databases, errors.New("Environment registry not found")
  }

  if len(registry_struct[reg.Env].CoreDBs) > 0 {
    databases.CoreDBs = registry_struct[reg.Env].CoreDBs
  }

  if len(registry_struct[reg.Env].TenantMetaDBs) > 0 {
    databases.TenantMetaDBs = registry_struct[reg.Env].TenantMetaDBs
  }

  return databases, nil
}

/**
 * Get service config
 *
 * @return {ServiceConfig}
 */
func (reg *RegistryObj) GetServiceConfig() (structs.ServiceConfig, error) {
  var serviceConfig structs.ServiceConfig
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return serviceConfig, errors.New("Environment registry not found")
  }

  serviceConfig = registry_struct[reg.Env].ServiceConfig
  return serviceConfig, nil
}

/**
 * Get depeloyer
 *
 * @return {Deployer}
 */
func (reg *RegistryObj) GetDeployer() (structs.Deployer, error) {
  var deployer structs.Deployer
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return deployer, errors.New("Environment registry not found")
  }

  deployer = registry_struct[reg.Env].Deployer
  return deployer, nil
}

/**
 * Get custom registry
 *
 * @return {CustomRegistry}
 */
func (reg *RegistryObj) GetCustom() (structs.CustomRegistry, error) {
  var customRegistry structs.CustomRegistry
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return customRegistry, errors.New("Environment registry not found")
  }

  customRegistry = registry_struct[reg.Env].Custom
  return customRegistry, nil
}

/**
 * Get one resource
 * @param  {String}     resourceName
 * @return {Resource}
 */
func (reg *RegistryObj) GetResource(resourceName string) (structs.Resource, error) {
  var resource structs.Resource

  if resourceName == "" {
    return resource, errors.New("Resource name is required")
  }

  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return resource, errors.New("Environment registry not found")
  }

  if len(registry_struct[reg.Env].Resources) == 0 || registry_struct[reg.Env].Resources[resourceName].Id == "" {
    return resource, errors.New("Resource not found")
  }

  resource = registry_struct[reg.Env].Resources[resourceName]
  return resource, nil
}

/**
 * Get all resources
 *
 * @return {Resources}
 */
func (reg *RegistryObj) GetResources() (structs.Resources, error) {
  var resources structs.Resources
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return resources, errors.New("Environment registry not found")
  }

  resources = registry_struct[reg.Env].Resources
  return resources, nil
}

/**
 * Get one service
 * @param  {String}     serviceName
 * @return {Service}
 */
func (reg *RegistryObj) GetService(serviceName string) (structs.Service, error) {
  var service structs.Service

  if serviceName == "" {
    return service, errors.New("Service name is required")
  }

  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return service, errors.New("Environment registry not found")
  }

  if len(registry_struct[reg.Env].Services) == 0 || registry_struct[reg.Env].Services[serviceName].Group == "" {
    return service, errors.New("Service not found")
  }

  service = registry_struct[reg.Env].Services[serviceName]
  return service, nil
}

/**
 * Get all services
 *
 * @return {Services}
 */
func (reg *RegistryObj) GetServices() (structs.Services, error) {
  var services structs.Services
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return services, errors.New("Environment registry not found")
  }

  services = registry_struct[reg.Env].Services
  return services, nil
}

/**
 * Get one daemon
 * @param  {String}     daemonName
 * @return {Daemon}
 */
func (reg *RegistryObj) GetDaemon(daemonName string) (structs.Daemon, error) {
  var daemon structs.Daemon

  if daemonName == "" {
    return daemon, errors.New("Daemon name is required")
  }

  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return daemon, errors.New("Environment registry not found")
  }

  if len(registry_struct[reg.Env].Daemons) == 0 || registry_struct[reg.Env].Daemons[daemonName].Id == "" {
    return daemon, errors.New("Daemon not found")
  }

  daemon = registry_struct[reg.Env].Daemons[daemonName]
  return daemon, nil
}

/**
 * Get all daemons
 *
 * @return {Daemons}
 */
func (reg *RegistryObj) GetDaemons() (structs.Daemons, error) {
  var daemons structs.Daemons
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return daemons, errors.New("Environment registry not found")
  }

  daemons = registry_struct[reg.Env].Daemons
  return daemons, nil
}

/**
 * Reload registry
 *
 * @return {Boolean}
 */
func (reg *RegistryObj) Reload() (bool, error) {
  if reg.Env == "" || registry_struct[reg.Env].Environment == "" {
    return false, errors.New("Cannot reload registry. Env and ServiceName are not set.")
  }

  param := map[string]string{"envCode": reg.Env, "serviceName": reg.ServiceName}
  go ExecRegistry(param) //TODO check return type of ExecRegistry

  return true, nil
}

/**
 * Call registry api
 *
 */
func ExecRegistry(param map[string]string) {
  registryApi := os.Getenv("SOAJS_REGISTRY_API")

  if index := strings.Index(registryApi, ":"); index == -1 {
    panic("Invalid format for SOAJS_REGISTRY_API [hostname:port]: " + registryApi)
  }

  registryApiPort := strings.Split(registryApi, ":")[1]
  if _, err := strconv.Atoi(registryApiPort); err != nil {
    panic("Port must be an integer [" + registryApiPort + "]" )
  }

  reqUrl := "http://" + registryApi + "/getRegistry?env=" + param["envCode"] + "&serviceName=" + param["serviceName"]
  httpResponse, err := http.Get(reqUrl)
  if(err != nil) {
    panic(err)
  }

  apiResponse, _ := ioutil.ReadAll(httpResponse.Body)

  var temp RegistryApiResponse
  json.Unmarshal(apiResponse, &temp)

  if temp.Result != true {
    panic(temp)
  }

  if len(registry_struct) == 0 {
    registry_struct = make(map[string]structs.Registry)
  }

  registry_struct[temp.Data.Environment] = temp.Data

  regObj.Env = param["envCode"];
  regObj.ServiceName = param["serviceName"];

  serviceConfig, _ := regObj.GetServiceConfig()
  log.Println(serviceConfig)
  //TODO assertion on service config content

  time.Sleep(time.Duration(serviceConfig.Awareness.AutoReloadRegistry) * time.Millisecond)
  go ExecRegistry(param)
}
