package soajsgo

type (
	SOA struct {
		Type          string `json:"type"`
		Prerequisites struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"prerequisites"`
		ServiceName           string `json:"serviceName"`
		ServiceIP             string `json:"serviceIP"`
		ServiceGroup          string `json:"serviceGroup"`
		ServicePort           int    `json:"servicePort"`
		Swagger               bool   `json:"swagger"`
		RequestTimeout        int    `json:"requestTimeout"`
		RequestTimeoutRenewal int    `json:"requestTimeoutRenewal"`
		ServiceVersion        string `json:"serviceVersion"`
		ExtKeyRequired        bool   `json:"extKeyRequired"`
		Urac                  bool   `json:"urac"`
		UracProfile           bool   `json:"urac_Profile"`
		UracACL               bool   `json:"urac_ACL"`
		ProvisionACL          bool   `json:"provision_ACL"`
		Oauth                 bool   `json:"oauth"`
		Maintenance           struct {
			Port struct {
				Type string `json:"type"`
			} `json:"port"`
			Readiness string `json:"readiness"`
			Commands  []struct {
				Label string `json:"label"`
				Path  string `json:"path"`
				Icon  string `json:"icon"`
			} `json:"commands"`
		} `json:"maintenance"`
	}

	RegisterConf struct {
		Name                  string `json:"name"`
		Type                  string `json:"type"`
		Mw                    bool   `json:"mw"`
		Group                 string `json:"group"`
		Port                  int    `json:"port"`
		Swagger               bool   `json:"swagger"`
		RequestTimeout        int    `json:"requestTimeout"`
		RequestTimeoutRenewal int    `json:"requestTimeoutRenewal"`
		Version               string `json:"version"`
		ExtKeyRequired        bool   `json:"extKeyRequired"`
		Urac                  bool   `json:"urac"`
		UracProfile           bool   `json:"urac_Profile"`
		UracACL               bool   `json:"urac_ACL"`
		ProvisionACL          bool   `json:"provision_ACL"`
		Oauth                 bool   `json:"oauth"`
		IP                    string `json:"ip"`
		Maintenance           struct {
			Port struct {
				Type string `json:"type"`
			} `json:"port"`
			Readiness string `json:"readiness"`
			Commands  []struct {
				Label string `json:"label"`
				Path  string `json:"path"`
				Icon  string `json:"icon"`
			} `json:"commands"`
		} `json:"maintenance"`
	}
)
