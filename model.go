package soajsgo

type (
	// SOA is the input config object needed by the middleware.
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

	// registerConf represents the config object to send to soajs gateway as post data.
	registerConf struct {
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

	// RegistryAPIResponse represents registry API response from soajs gateway.
	registryAPIResponse struct {
		Result  bool  `json:"result"`
		Ts      int64 `json:"ts"`
		Service struct {
			ServiceName string `json:"service"`
			Type        string `json:"type"`
			Route       string `json:"route"`
		} `json:"service"`
		Registry Registry `json:"data"`
	}
	// RegisterAPIResponse represents registry API response from soajs gateway.
	registerAPIResponse struct {
		Result  bool  `json:"result"`
		Ts      int64 `json:"ts"`
		Service struct {
			ServiceName string `json:"service"`
			Type        string `json:"type"`
			Route       string `json:"route"`
		} `json:"service"`
	}

	// HeaderInfo represents header info structure.
	headerInfo struct {
		Tenant      Tenant      `json:"tenant"`
		Key         Key         `json:"key"`
		Application Application `json:"application"`
		Package     struct {
			ACL       interface{} `json:"acl"`
			ACLAllEnv interface{} `json:"acl_all_env"`
		} `json:"package"`
		Device    string            `json:"device"`
		Geo       map[string]string `json:"geo"`
		Urac      Urac              `json:"urac"`
		Awareness Host              `json:"awareness"`
		Param     struct {
			UracProfile bool `json:"urac_Profile"`
			UracACL     bool `json:"urac_ACL"`
		} `json:"param"`
	}

	// ContextData represents http context data information.
	ContextData struct {
		Tenant         Tenant            `json:"tenant"`
		Urac           Urac              `json:"urac"`
		ServicesConfig interface{}       `json:"servicesConfig"`
		Device         string            `json:"device"`
		Geo            map[string]string `json:"geo"`
		Awareness      Host              `json:"awareness"`
		Reg            Registry          `json:"reg"`
	}

	// Tenant contains the tenant information.
	Tenant struct {
		ID      string      `json:"id"`
		Code    string      `json:"code"`
		Locked  bool        `json:"locked"`
		Roaming interface{} `json:"roaming,omitempty"` // TODO: implement struct

		Key         Key         `json:"key"`
		Application Application `json:"application,omitempty"`
	}
	// Key represents the key that is making the call to the API.
	Key struct {
		Config interface{} `json:"config"`
		IKey   string      `json:"iKey"`
		EKey   string      `json:"eKey"`
	}
	// Application represents the product that is making the call to the API.
	Application struct {
		Product   string      `json:"product"`
		Package   string      `json:"package"`
		AppID     string      `json:"appId"`
		ACL       interface{} `json:"acl"`
		ACLAllEnv interface{} `json:"acl_all_env"`

		PackageACL       interface{} `json:"package_acl"`
		PackageACLAllEnv interface{} `json:"package_acl_all_env"`
	}

	// Urac is the logged in user record in case urac is set to true.
	Urac struct {
		ID          string      `json:"_id"`
		Username    string      `json:"username"`
		FirstName   string      `json:"firstName"`
		LastName    string      `json:"lastName"`
		Email       string      `json:"email"`
		Groups      []string    `json:"groups"`
		SocialLogin interface{} `json:"socialLogin"`
		Tenant      struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"tenant"`
		Profile   interface{} `json:"profile"`
		ACL       interface{} `json:"acl"`
		ACLAllEnv interface{} `json:"acl_all_env"`
	}

	// Host represents host information.
	Host struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	// Registry represents registry structure.
	Registry struct {
		url         string `json:"url"`
		timeLoaded  int64  `json:"timeLoaded"`
		Name        string `json:"name"`
		Environment string `json:"environment"`

		CoreDBs       Databases `json:"coreDB"`
		TenantMetaDBs Databases `json:"tenantMetaDB"`

		ServiceConfig ServiceConfig    `json:"serviceConfig"`
		Custom        CustomRegistries `json:"custom"`
		Resources     Resources        `json:"resources"`
		Services      Services         `json:"services"`
	}

	// ServiceConfig represents service config.
	ServiceConfig struct {
		Awareness struct {
			CacheTTL            int  `json:"cacheTTL"`
			HealthCheckInterval int  `json:"healthCheckInterval"`
			AutoReloadRegistry  int  `json:"autoRelaodRegistry"`
			MaxLogCount         int  `json:"maxLogCount"`
			AutoRegisterService bool `json:"autoRegisterService"`
		} `json:"awareness"`
		Agent struct {
			TopologyDir string `json:"topologyDir"`
		} `json:"agent"`
		Key struct {
			Algorithm string `json:"algorithm"`
			Password  string `json:"password"`
		} `json:"key"`
		Logger struct {
			Src       bool   `json:"src"`
			Level     string `json:"level"`
			Formatter struct {
				LevelInString bool   `json:"levelInString"`
				OutputMode    string `json:"outputMode"`
			} `json:"formatter"`
		} `json:"logger"`
		Port struct {
			Controller     int `json:"controller"`
			MaintenanceInc int `json:"maintenanceInc"`
			RandomInc      int `json:"randomInc"`
		} `json:"ports"`
		Cookie struct {
			Secret string `json:"secret"`
		} `json:"cookie"`
		Session struct {
			Name   string `json:"name"`
			Secret string `json:"secret"`
			Cookie struct {
				Path     string      `json:"path"`
				HTTPOnly bool        `json:"httpOnly"`
				Secure   bool        `json:"secure"`
				MaxAge   interface{} `json:"maxAge"`
			} `json:"cookie"`
			Resave            bool   `json:"resave"`
			SaveUninitialized bool   `json:"saveUninitialized"`
			Rolling           bool   `json:"rolling"`
			Unset             string `json:"unset"`
		} `json:"session"`
	}

	Databases        map[string]Database
	CustomRegistries map[string]CustomRegistry
	Resources        map[string]map[string]Resource
	Services         map[string]Service

	// Database represents a Database structure with configuration fields.
	Database struct {
		Name    string `json:"name"`
		Prefix  string `json:"prefix"`
		Cluster string `json:"cluster"`
		Server  []struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"servers"`
		Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"credentials"`
		Streaming        interface{} `json:"streaming"`
		RegistryLocation struct {
			L1  string `json:"l1"`
			L2  string `json:"l2"`
			Env string `json:"env"`
		} `json:"registryLocation"`
		URLParam   interface{} `json:"URLParam"`
		ExtraParam interface{} `json:"extraParam"`

		// NOTE: session specific entries
		Store       interface{} `json:"store,omitempty"`
		Collection  string      `json:"collection,omitempty"`
		Stringify   bool        `json:"stringify,omitempty"`
		ExpireAfter int         `json:"expireAfter,omitempty"`
	}

	// CustomRegistry represents custom registry information.
	CustomRegistry struct {
		ID      string      `json:"_id"`
		Name    string      `json:"name"`
		Locked  bool        `json:"locked"`
		Plugged bool        `json:"plugged"`
		Shared  bool        `json:"shared"`
		Value   interface{} `json:"value"`
		Created string      `json:"created"`
		Author  string      `json:"author"`
	}

	// Resource represents resource structure.
	Resource struct {
		ID       string      `json:"_id"`
		Name     string      `json:"name"`
		Type     string      `json:"type"`
		Category string      `json:"category"`
		Created  string      `json:"created"`
		Author   string      `json:"author"`
		Locked   bool        `json:"locked"`
		Plugged  bool        `json:"plugged"`
		Shared   bool        `json:"shared"`
		Config   interface{} `json:"config"`
	}

	// Service represents service structure.
	Service struct {
		Group                 string `json:"group"`
		Port                  int    `json:"port"`
		RequestTimeout        int    `json:"requestTimeout"`
		RequestTimeoutRenewal int    `json:"requestTimeoutRenewal"`
		MaxPoolSize           int    `json:"maxPoolSize"`
		Version               string `json:"version"`
		Authorization         bool   `json:"authorization"`
		ExtKeyRequired        bool   `json:"extKeyRequired"`
	}
)
