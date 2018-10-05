package structs

type Registry struct {
  TimeLoaded            int64                       `json:"timeLoaded"`
  Name                  string                      `json:"name"`
  Environment           string                      `json:"environment"`
  ProfileOnly           bool                        `json:"profileOnly"`

  Domain                string                      `json:"domain"`
  ApiPrefix             string                      `json:"apiPrefix"`
  SitePrefix            string                      `json:"sitePrefix"`
  Protocol              string                      `json:"protocol"`
  Port                  int                         `json:"port"`

  CoreDBs               map[string]Database         `json:"coreDB"`
  TenantMetaDBs         map[string]Database         `json:"tenantMetaDB"`

  ServiceConfig         ServiceConfig               `json:"serviceConfig"`
  Deployer              Deployer                    `json:"deployer"`
  Custom                CustomRegistries            `json:"custom"`
  Resources             Resources                   `json:"resources"`
  Services              Services                    `json:"services"`
  Daemons               Daemons                     `json:daemons`
}
