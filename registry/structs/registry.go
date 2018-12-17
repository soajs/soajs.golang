package structs

type Registry struct {
	TimeLoaded  int64  `json:"timeLoaded"`
	Name        string `json:"name"`
	Environment string `json:"environment"`

	CoreDBs       map[string]Database `json:"coreDB"`
	TenantMetaDBs map[string]Database `json:"tenantMetaDB"`

	ServiceConfig ServiceConfig    `json:"serviceConfig"`
	Custom        CustomRegistries `json:"custom"`
	Resources     Resources        `json:"resources"`
	Services      Services         `json:"services"`
}
