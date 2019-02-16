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

	// HeaderInfo represents header info structure.
	HeaderInfo struct {
		Tenant      Tenant            `json:"tenant"`
		Key         Key               `json:"key"`
		Application Application       `json:"application"`
		Package     Package           `json:"package"`
		Device      string            `json:"device"`
		Geo         map[string]string `json:"geo"`
		Urac        Urac              `json:"urac"`
		Awareness   Host              `json:"awareness"`
		Param       Param             `json:"param"`
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
	// Package represents the Productization that the tenant making the call to the API is using.
	Package struct {
		ACL       interface{} `json:"acl"`
		ACLAllEnv interface{} `json:"acl_all_env"`
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
	// Param represents Urac params.
	Param struct {
		UracProfile bool `json:"urac_Profile"`
		UracACL     bool `json:"urac_ACL"`
	}
	// Host represents host information.
	Host struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	// RegistryAPIResponse represents registry API response from soajs gateway
	RegistryAPIResponse struct {
		Result  bool  `json:"result"`
		Ts      int64 `json:"ts"`
		Service struct {
			ServiceName string `json:"service"`
			Type        string `json:"type"`
			Route       string `json:"route"`
		} `json:"service"`
		Registry Registry `json:"data"`
	}
	// RegisterAPIResponse represents registry API response from soajs gateway
	RegisterAPIResponse struct {
		Result  bool  `json:"result"`
		Ts      int64 `json:"ts"`
		Service struct {
			ServiceName string `json:"service"`
			Type        string `json:"type"`
			Route       string `json:"route"`
		} `json:"service"`
	}

	// Registry represents registry structure.
	Registry struct {
		Url         string `json:"url"`
		TimeLoaded  int64  `json:"timeLoaded"`
		Name        string `json:"name"`
		Environment string `json:"environment"`

		CoreDBs       Databases `json:"coreDB"`
		TenantMetaDBs Databases `json:"tenantMetaDB"`

		ServiceConfig ServiceConfig    `json:"serviceConfig"`
		Custom        CustomRegistries `json:"custom"`
		Resources     Resources        `json:"resources"`
		Services      Services         `json:"services"`
	}

	Databases        map[string]Database
	CustomRegistries map[string]CustomRegistry
	Resources        map[string]map[string]Resource
	Services         map[string]Service

	// Database represents a Database structure with configuration fields.
	Database struct {
		Name             string           `json:"name"`
		Prefix           string           `json:"prefix"`
		Cluster          string           `json:"cluster"`
		Server           []Host           `json:"servers"`
		Credentials      Credentials      `json:"credentials"`
		Streaming        interface{}      `json:"streaming"`
		RegistryLocation RegistryLocation `json:"registryLocation"`
		URLParam         interface{}      `json:"URLParam"`
		ExtraParam       interface{}      `json:"extraParam"`

		// NOTE: session specific entries
		Store       interface{} `json:"store,omitempty"`
		Collection  string      `json:"collection,omitempty"`
		Stringify   bool        `json:"stringify,omitempty"`
		ExpireAfter int         `json:"expireAfter,omitempty"`
	}
	// Credentials contains username and password.
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// RegistryLocation represents database location.
	RegistryLocation struct {
		L1  string `json:"l1"`
		L2  string `json:"l2"`
		Env string `json:"env"`
	}
	// ServiceConfig represents service config.
	ServiceConfig struct {
		Awareness ServiceConfigIntervals `json:"awareness"`
		Agent     Agent                  `json:"agent"`
		Key       ServiceKey             `json:"key"`
		Logger    Logger                 `json:"logger"`
		Port      ServicePort            `json:"ports"`
		Cookie    Cookie                 `json:"cookie"`
		Session   Session                `json:"session"`
	}
	// ServiceConfigIntervals represents amount of time in milliseconds.
	ServiceConfigIntervals struct {
		CacheTTL            int  `json:"cacheTTL"`
		HealthCheckInterval int  `json:"healthCheckInterval"`
		AutoReloadRegistry  int  `json:"autoRelaodRegistry"`
		MaxLogCount         int  `json:"maxLogCount"`
		AutoRegisterService bool `json:"autoRegisterService"`
	}
	// Agent contains topology direction.
	Agent struct {
		TopologyDir string `json:"topologyDir"`
	}
	// ServiceKey represents which algorithm should bcrypt use and secret phrase to encrypt tenant key security accordingly.
	ServiceKey struct {
		Algorithm string `json:"algorithm"`
		Password  string `json:"password"`
	}
	// Logger represents custom logger object.
	Logger struct {
		Src       bool      `json:"src"`
		Level     string    `json:"level"`
		Formatter Formatter `json:"formatter"`
	}
	// Formatter contains specific configuration for logging formatting.
	Formatter struct {
		LevelInString bool   `json:"levelInString"`
		OutputMode    string `json:"outputMode"`
	}
	// ServicePort provides default ports.
	ServicePort struct {
		Controller     int `json:"controller"`
		MaintenanceInc int `json:"maintenanceInc"`
		RandomInc      int `json:"randomInc"`
	}
	// Cookie contents the cookie secret phrase, used to encrypt cookie values, minimum 5 characters.
	Cookie struct {
		Secret string `json:"secret"`
	}
	// Session is session data from the multi-tenant session.
	Session struct {
		Name              string        `json:"name"`
		Secret            string        `json:"secret"`
		Cookie            SessionCookie `json:"cookie"`
		Resave            bool          `json:"resave"`
		SaveUninitialized bool          `json:"saveUninitialized"`
		Rolling           bool          `json:"rolling"`
		Unset             string        `json:"unset"`
	}
	// SessionCookie represents path where cookies should be created and other options related to cookies.
	SessionCookie struct {
		Path     string      `json:"path"`
		HTTPOnly bool        `json:"httpOnly"`
		Secure   bool        `json:"secure"`
		MaxAge   interface{} `json:"maxAge"`
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
