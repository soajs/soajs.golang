package soajsgo

import (
	"errors"
	"fmt"
	"regexp"
)

type (
	// Config represent service configuration from json file.
	Config struct {
		ServiceName           string       `json:"name"`
		ServiceGroup          string       `json:"group"`
		ServicePort           int          `json:"port"`
		ServiceIP             string       `json:"IP"`
		Type                  string       `json:"type"`
		ServiceVersion        string       `json:"version"`
		SubType               string       `json:"subType"`
		Description           string       `json:"description"`
		Oauth                 bool         `json:"oauth"`
		Urac                  bool         `json:"urac"`
		UracProfile           bool         `json:"urac_Profile"`
		UracACL               bool         `json:"urac_ACL"`
		UracConfig            bool         `json:"urac_Config"`
		UracGroupConfig       bool         `json:"urac_GroupConfig"`
		TenantProfile         bool         `json:"tenant_Profile"`
		ProvisionACL          bool         `json:"provision_ACL"`
		ExtKeyRequired        bool         `json:"extKeyRequired"`
		RequestTimeout        int          `json:"requestTimeout"`
		RequestTimeoutRenewal int          `json:"requestTimeoutRenewal"`
		Maintenance           maintenance  `json:"maintenance"`
		InterConnect          interconnect `json:"interConnect"`
		Prerequisites         struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"prerequisites"`
	}
)

// Validate validates soajs config.
func (c *Config) Validate() error {

	var validator = regexp.MustCompile(`^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`)
	var versionRegexp = regexp.MustCompile(`[0-9]+(.[0-9]+)?`)

	if c.Type == "" {
		return errors.New("could not find [Type] in your config, type is <required>")
	}

	if c.ServiceName == "" {
		return errors.New("could not find [ServiceName] in your config, name is <required>")
	}
	if !validator.MatchString(c.ServiceName) {
		return fmt.Errorf("error with [ServiceName] in your config, name syntax is [%s]", validator)
	}

	if c.ServicePort == 0 {
		return errors.New("could not find [ServicePort] in your config, port is <required>")
	}

	if c.ServiceVersion == "" {
		return errors.New("could not find [ServiceVersion] in your config, version is <required>")
	}
	if !versionRegexp.MatchString(c.ServiceVersion) {
		return fmt.Errorf("error with [ServiceVersion] in your config, version syntax is [%s]", versionRegexp)
	}

	if c.Maintenance.Readiness == "" {
		return errors.New("could not find [Readiness] in your config, maintenance.readiness is <required>")
	}
	if c.Maintenance.Port.Type == "" {
		return errors.New("could not find [Maintenance Port Type] in your config, maintenance.port.type is <required>")
	}

	if c.ServiceGroup == "" {
		return errors.New("could not find [ServiceGroup] in your config, group is <required>")
	}
	if !validator.MatchString(c.ServiceGroup) {
		return fmt.Errorf("error with [ServiceGroup] in your config, group syntax is [%s]", validator)
	}
	return nil
}
