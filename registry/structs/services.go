package structs

type Services map[string]Service

type Service struct {
    Group                 string                        `json:"group"`
    Port                  int                           `json:"port"`
    RequestTimeout        int                           `json:"requestTimeout"`
    RequestTimeoutRenewal int                           `json:"requestTimeoutRenewal"`
    MaxPoolSize           int                           `json:"maxPoolSize"`
    Authorization         bool                          `json:"authorization"`
    Version               int                           `json:"version"`
    ExtKeyRequired        bool                          `json:"extKeyRequired"`
    Versions              map[string]ServiceVersion     `json:"versions"`
}

type ServiceVersion struct {
    ExtKeyRequired         bool                         `json:"extKeyRequired"`
	Urac                   bool                         `json:"urac"`
	UracProfile            bool                         `json:"urac_Profile"`
	UracACL                bool                         `json:"urac_ACL"`
	ProvisionACL           bool                         `json:"provision_ACL"`
	Oauth                  bool                         `json:"oauth"`
    Apis                   []ServiceVersionApis         `json:"apis"`
}

type ServiceVersionApis struct {
    L                      string                       `json:"l"`
    V                      string                       `json:"v"`
    M                      string                       `json:"m"`
    Group                  string                       `json:"group"`
    GroupMain              bool                         `json:"groupMain,omitempty"`
}
