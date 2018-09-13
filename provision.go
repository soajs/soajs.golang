package soajsGo

type Tenant struct {
  Id                        string                    `json:"id"`
  Code                      string                    `json:"code"`
  Key                       Key                       `json:"key"`
  Application               Application               `json:"application"`
}

type KeyData struct {
  Config                    map[string]string         `json:"config"`
  IKey                      string                    `json:"iKey"`
  EKey                      string                    `json:"eKey"`
}

type Key struct {
  IKey                      string                    `json:"iKey"`
  EKey                      string                    `json:"eKey"`
}

type Application struct {
  Product                   string                    `json:"product"`
  Package                   string                    `json:"package"`
  AppId                     string                    `json:"appId"`
  Acl                       map[string]string         `json:"acl"`
  Acl_all_env               map[string]string         `json:"acl_all_env"`
  Package_acl               map[string]string         `json:"package_acl"`
  Package_acl_all_env       map[string]string         `json:"package_acl_all_env"`
}

type Package struct {
  Acl                       map[string]string         `json:"acl"`
  Acl_all_env               map[string]string         `json:"acl_all_env"`
}

type Awareness struct {
  Host                      string                    `json:"host"`
  Port                      string                    `json:"port"`
}

type SOAJSData struct {
  Tenant                    Tenant                    `json:"tenant"`
  Key                       KeyData                   `json:"key"`
  Application               Application               `json:"application"`
  Package                   Package                   `json:"package"`
  Device                    string                    `json:"device"`
  Geo                       map[string]string         `json:"geo"`
  Urac                      map[string]string         `json:"urac"`
  Awareness                 Awareness                 `json:"awareness"`
}

type SOAJSObject struct {
  Tenant                    Tenant                    `json:"tenant"`
  Urac                      map[string]string         `json:"urac"`
  ServicesConfig            map[string]string         `json:"servicesConfig"`
  Device                    string                    `json:"device"`
  Geo                       map[string]string         `json:"geo"`
  Awareness                 Awareness                 `json:"awareness"`
  Controller                string                    `json:"controller"`
}
