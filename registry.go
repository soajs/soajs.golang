package soajsGo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/soajs/soajs.golang/registry/structs"
)

type RegistryObj struct {
	Env         string `json:"env"`
	ServiceName string `json:"serviceName"`
}

type RegistryApiResponse struct {
	Result  bool              `json:"result"`
	Ts      int64             `json:"ts"`
	Service map[string]string `json:"service"`
	Data    structs.Registry  `json:"data"`
}

var (
	registryStruct map[string]structs.Registry
	regObj         RegistryObj
)

var autoReloadChannel = make(chan string)

/**
 * Check if the environment registry exists
 *
 */
func detectEnvRegistry(reg *RegistryObj) error {
	if reg.Env == "" || registryStruct[reg.Env].Environment == "" {
		return errors.New("environment registry not found")
	}

	return nil
}

/**
 * Get one database
 * @param  {String}     dbName
 * @return {Database}
 */
func (reg *RegistryObj) GetDatabase(dbName string) (structs.Database, error) {
	var database structs.Database

	if dbName == "" {
		return database, errors.New("database name is required")
	}

	if err := detectEnvRegistry(reg); err != nil {
		return database, err
	}

	if len(registryStruct[reg.Env].CoreDBs) > 0 && registryStruct[reg.Env].CoreDBs[dbName].Name != "" {
		database = registryStruct[reg.Env].CoreDBs[dbName]
	} else if len(registryStruct[reg.Env].TenantMetaDBs) > 0 && registryStruct[reg.Env].TenantMetaDBs[dbName].Name != "" {
		database = registryStruct[reg.Env].TenantMetaDBs[dbName]
	} else {
		return database, errors.New("database not found")
	}

	return database, nil
}

/**
 * Get all databases
 *
 * @return {Databases}
 */
func (reg *RegistryObj) GetDatabases() (structs.Databases, error) {
	var databases structs.Databases
	if err := detectEnvRegistry(reg); err != nil {
		return databases, err
	}

	if len(registryStruct[reg.Env].CoreDBs) > 0 {
		databases = registryStruct[reg.Env].CoreDBs
	}

	if len(registryStruct[reg.Env].TenantMetaDBs) > 0 {
		for dbName, dbConfig := range registryStruct[reg.Env].TenantMetaDBs {
			databases[dbName] = dbConfig
		}
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
	if err := detectEnvRegistry(reg); err != nil {
		return serviceConfig, err
	}

	serviceConfig = registryStruct[reg.Env].ServiceConfig
	return serviceConfig, nil
}

/**
 * Get custom registry
 *
 * @return {CustomRegistries}
 */
func (reg *RegistryObj) GetCustom() (structs.CustomRegistries, error) {
	var customRegistry structs.CustomRegistries
	if err := detectEnvRegistry(reg); err != nil {
		return customRegistry, err
	}

	customRegistry = registryStruct[reg.Env].Custom
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
		return resource, errors.New("resource name is required")
	}

	if err := detectEnvRegistry(reg); err != nil {
		return resource, err
	}

	if len(registryStruct[reg.Env].Resources) == 0 {
		return resource, errors.New("resource not found")
	}

	for _, resourceList := range registryStruct[reg.Env].Resources {
		for resourceKey, resourceData := range resourceList {
			if resourceKey == resourceName {
				resource = resourceData
			}
		}
	}

	if resource == (structs.Resource{}) {
		return resource, errors.New("resource not found")
	}

	return resource, nil
}

/**
 * Get all resources
 *
 * @return {Resources}
 */
func (reg *RegistryObj) GetResources() (structs.Resources, error) {
	var resources structs.Resources
	if err := detectEnvRegistry(reg); err != nil {
		return resources, err
	}

	resources = registryStruct[reg.Env].Resources
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
		return service, errors.New("service name is required")
	}

	if err := detectEnvRegistry(reg); err != nil {
		return service, err
	}

	if len(registryStruct[reg.Env].Services) == 0 || registryStruct[reg.Env].Services[serviceName].Group == "" {
		return service, errors.New("service not found")
	}

	service = registryStruct[reg.Env].Services[serviceName]
	return service, nil
}

/**
 * Get all services
 *
 * @return {Services}
 */
func (reg *RegistryObj) GetServices() (structs.Services, error) {
	var services structs.Services
	if err := detectEnvRegistry(reg); err != nil {
		return services, err
	}

	services = registryStruct[reg.Env].Services
	return services, nil
}

/**
 * Reload registry
 *
 * @return {Boolean}
 */
func (reg *RegistryObj) Reload() (bool, error) {
	if reg.Env == "" || reg.ServiceName == "" {
		return false, errors.New("cannot reload registry env and serviceName are not set")
	}

	param := map[string]string{"envCode": reg.Env, "serviceName": reg.ServiceName}
	ExecRegistry(param) //TODO check return type of execRegistry

	autoReloadChannel <- "reset"

	return true, nil
}

/**
 * Call registry api
 *
 */
func ExecRegistry(param map[string]string) (RegistryObj, error) {
	registryApi := os.Getenv("SOAJS_REGISTRY_API")

	if index := strings.Index(registryApi, ":"); index == -1 {
		return RegistryObj{}, errors.New("Invalid format for SOAJS_REGISTRY_API [hostname:port]: " + registryApi)
	}

	registryApiPort := strings.Split(registryApi, ":")[1]
	if _, err := strconv.Atoi(registryApiPort); err != nil {
		return RegistryObj{}, errors.New("Port must be an integer [" + registryApiPort + "]")
	}

	reqUrl := "http://" + registryApi + "/getRegistry?env=" + param["envCode"] + "&serviceName=" + param["serviceName"]
	httpResponse, err := http.Get(reqUrl)
	if err != nil {
		return RegistryObj{}, errors.New("unable to get registry from api gateway")
	} else {
		defer httpResponse.Body.Close()
	}
	apiResponse, _ := ioutil.ReadAll(httpResponse.Body)

	var temp RegistryApiResponse
	err = json.Unmarshal(apiResponse, &temp)

	if (err != nil && temp.Result != true) || temp.Result != true {
		return RegistryObj{}, errors.New("unable to convert registry to json from returned api gateway response")
	}

	if len(registryStruct) == 0 {
		registryStruct = make(map[string]structs.Registry)
	}

	registryStruct[temp.Data.Environment] = temp.Data

	regObj.Env = param["envCode"]
	regObj.ServiceName = param["serviceName"]
	return regObj, nil
}

func autoReload(param map[string]string) chan string {
	log.Println("auto reloading ...")
	regObj, err := ExecRegistry(param)
	if err != nil {
		log.Println(err)
	} else {
		serviceConfig, _ := regObj.GetServiceConfig()
		//TODO assertion on service config content

		interval := time.Duration(serviceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
		ticker := time.NewTicker(interval)

		go func() {
			for {

				select {
				case <-ticker.C:
					log.Println("Reloading ...")
					go ExecRegistry(param)

					serviceConfig, _ := regObj.GetServiceConfig()
					interval = time.Duration(serviceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
					ticker = time.NewTicker(interval)
				case msg := <-autoReloadChannel:
					if msg == "stop" {
						ticker.Stop()
						return
					} else if msg == "reset" {
						serviceConfig, _ := regObj.GetServiceConfig()
						interval = time.Duration(serviceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
						ticker = time.NewTicker(interval)
					}

				}

			}
		}()
	}
	return autoReloadChannel
}
