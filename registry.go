package soajsgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RegistryObj struct {
	Env         string `json:"env"`
	ServiceName string `json:"serviceName"`
}

var autoReloadChannel = make(chan string)

func newRegistry(reqURL string, turnOnAutoReload bool) (*Registry, error) {
	var res, err = http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("could not get registry from soajs gateway: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("non 2xx status code: %d %v", res.StatusCode, b)
	}
	var temp RegistryAPIResponse
	err = json.NewDecoder(res.Body).Decode(&temp)
	if err != nil || !temp.Result {
		return nil, fmt.Errorf("unable to convert registry response: %v", err)
	}

	temp.Registry.url = reqURL

	if turnOnAutoReload {
		go autoReload(&temp.Registry)
	}

	return &temp.Registry, nil
}

func autoReload(reg *Registry) chan string {
	interval := time.Duration(reg.ServiceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
	ticker := time.NewTicker(interval)

	go func() {
		for {

			select {
			case <-ticker.C:
				temp, err := newRegistry(reg.url, false)
				if err == nil {
					*reg = *temp
					interval = time.Duration(reg.ServiceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
					ticker = time.NewTicker(interval)
				}
			case msg := <-autoReloadChannel:
				if msg == "stop" {
					ticker.Stop()
					return
				} else if msg == "reset" {
					interval = time.Duration(reg.ServiceConfig.Awareness.AutoReloadRegistry) * time.Millisecond
					ticker = time.NewTicker(interval)
				}

			}

		}
	}()

	return autoReloadChannel
}

/**
 * Reload registry
 *
 * @return {Boolean}
 */
func (reg *Registry) Reload() (bool, error) {

	temp, err := newRegistry(reg.url, false)
	if err == nil {
		*reg = *temp
		autoReloadChannel <- "reset"
		return true, nil
	}
	return false, err
}

/**
 * Get one service
 * @param  {String}     name
 * @return {Service}
 */
func (reg Registry) GetService(name string) (*Service, error) {
	var s Service
	if name == "" {
		return nil, errors.New("service name is required")
	}

	if len(reg.Services) == 0 || reg.Services[name].Group == "" {
		return nil, errors.New("service not found")
	}

	s = reg.Services[name]
	return &s, nil
}

/**
 * Get one resource
 * @param  {String}     name
 * @return {Resource}
 */
func (reg Registry) GetResource(name string) (*Resource, error) {
	if name == "" {
		return nil, errors.New("resource name is required")
	}

	if len(reg.Resources) == 0 {
		return nil, errors.New("resource not found")
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

/**
 * Get one database
 * @param  {String}     name
 * @return {Database}
 */
func (reg Registry) GetDatabase(name string) (*Database, error) {
	if name == "" {
		return nil, errors.New("database name is required")
	}
	var d Database

	if len(reg.CoreDBs) > 0 && reg.CoreDBs[name].Name != "" {
		d = reg.CoreDBs[name]
	} else if len(reg.TenantMetaDBs) > 0 && reg.TenantMetaDBs[name].Name != "" {
		d = reg.TenantMetaDBs[name]
	} else {
		return nil, errors.New("database not found")
	}

	return &d, nil
}
