package structs

type ServiceConfig struct {
  Awareness       Awareness       `json:"awareness"`
  Agent           Agent           `json:"agent"`
  Key             Key             `json:"key"`
  Logger          Logger          `json:"logger"`
  Cors            Cors            `json:"cors"`
  Oauth           Oauth           `json:"oauth"`
  Ports           Ports           `json:"ports"`
  Cookie          Cookie          `json:"cookie"`
  Session         Session         `json:"session"`
}

type Awareness struct {
  CacheTTL            int  `json:"cacheTTL"`
  HealthCheckInterval int  `json:"healthCheckInterval"`
  AutoReloadRegistry  int  `json:"autoRelaodRegistry"`
  MaxLogCount         int  `json:"maxLogCount"`
  AutoRegisterService bool `json:"autoRegisterService"`
}

type Agent struct {
  TopologyDir string `json:"topologyDir"`
}

type Key struct {
  Algorithm string `json:"algorithm"`
  Password  string `json:"password"`
}

type Logger struct {
  Src       bool   `json:"src"`
  Level     string `json:"level"`
  Formatter Formatter `json:"formatter"`
}

type Formatter struct {
  LevelInString bool   `json:"levelInString"`
  OutputMode    string `json:"outputMode"`
}

type Cors struct {
  Enabled     bool   `json:"enabled"`
  Origin      string `json:"origin"`
  Credentials string `json:"credentials"`
  Methods     string `json:"methods"`
  Headers     string `json:"headers"`
  Maxage      int    `json:"maxage"`
}

type Oauth struct {
  Grants               []string `json:"grants"`
  Debug                bool     `json:"debug"`
  AccessTokenLifetime  int      `json:"accessTokenLifetime"`
  RefreshTokenLifetime int      `json:"refreshTokenLifetime"`
}

type Ports struct {
  Controller     int `json:"controller"`
  MaintenanceInc int `json:"maintenanceInc"`
  RandomInc      int `json:"randomInc"`
}

type Cookie struct {
  Secret string `json:"secret"`
}

type Session struct {
  Name   string `json:"name"`
  Secret string `json:"secret"`
  Cookie SessionCookie `json:"cookie"`
  Resave            bool   `json:"resave"`
  SaveUninitialized bool   `json:"saveUninitialized"`
  Rolling           bool   `json:"rolling"`
  Unset             string `json:"unset"`
}

type SessionCookie struct {
  Path     string      `json:"path"`
  HTTPOnly bool        `json:"httpOnly"`
  Secure   bool        `json:"secure"`
  MaxAge   interface{} `json:"maxAge"`
}
