package structs

type ServiceConfig struct {
	Awareness Awareness `json:"awareness"`
	Agent     Agent     `json:"agent"`
	Key       Key       `json:"key"`
	Logger    Logger    `json:"logger"`
	Ports     Ports     `json:"ports"`
	Cookie    Cookie    `json:"cookie"`
	Session   Session   `json:"session"`
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
	Src       bool      `json:"src"`
	Level     string    `json:"level"`
	Formatter Formatter `json:"formatter"`
}

type Formatter struct {
	LevelInString bool   `json:"levelInString"`
	OutputMode    string `json:"outputMode"`
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
	Name              string        `json:"name"`
	Secret            string        `json:"secret"`
	Cookie            SessionCookie `json:"cookie"`
	Resave            bool          `json:"resave"`
	SaveUninitialized bool          `json:"saveUninitialized"`
	Rolling           bool          `json:"rolling"`
	Unset             string        `json:"unset"`
}

type SessionCookie struct {
	Path     string      `json:"path"`
	HTTPOnly bool        `json:"httpOnly"`
	Secure   bool        `json:"secure"`
	MaxAge   interface{} `json:"maxAge"`
}
