package soajsGo

import "reflect"

type Tenant struct {
	Id          string      `json:"id"`
	Code        string      `json:"code"`
	Locked      bool        `json:"locked"`
	Key         Key         `json:"key"`
	Roaming     interface{} `json:"roaming,omitempty"` //TODO implement struct
	Application Application `json:"application,omitempty"`
}

type KeyData struct {
	Config map[string]interface{} `json:"config"`
	IKey   string                 `json:"iKey"`
	EKey   string                 `json:"eKey"`
}

type Key struct {
	IKey string `json:"iKey"`
	EKey string `json:"eKey"`
}

type Application struct {
	Product             string                 `json:"product"`
	Package             string                 `json:"package"`
	AppId               string                 `json:"appId"`
	Acl                 map[string]interface{} `json:"acl"`
	Acl_all_env         map[string]interface{} `json:"acl_all_env"`
	Package_acl         map[string]interface{} `json:"package_acl"`
	Package_acl_all_env map[string]interface{} `json:"package_acl_all_env"`
}

type Package struct {
	Acl         map[string]interface{} `json:"acl"`
	Acl_all_env map[string]interface{} `json:"acl_all_env"`
}

type Awareness struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Urac struct {
	Id        string   `json:"_id"`
	Username  string   `json:"username"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Groups    []string `json:"groups"`
	Tenant    struct {
		ID   string `json:"id"`
		Code string `json:"code"`
	} `json:"tenant"`
	Profile   interface{} `json:"profile"`
	Acl       interface{} `json:"acl"`
	AclAllEnv interface{} `json:"acl_AllEnv"`
}

type Param struct {
	Urac_Profile bool `json:"urac_profile"`
	Urac_ACL     bool `json:"urac_ACL"`
}

type SOAJSData struct {
	Tenant      Tenant            `json:"tenant"`
	Key         KeyData           `json:"key"`
	Application Application       `json:"application"`
	Package     Package           `json:"package"`
	Device      string            `json:"device"`
	Geo         map[string]string `json:"geo"`
	Urac        Urac              `json:"urac"`
	Awareness   Awareness         `json:"awareness"`
	Param       Param             `json:"param"`
}

func (m SOAJSData) IsEmpty() bool {
	return reflect.DeepEqual(SOAJSData{}, m)
}
