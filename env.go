package soajsgo

const (
	// EnvRegistryAPIAddress is the environment variable name that contains the IP address and port of
	// the controller service that runs in the same environment. The SOAJS middleware uses this variable to fetch
	// the registry of this environment and supply it to your service.
	EnvRegistryAPIAddress = "SOAJS_REGISTRY_API"

	// EnvSoajsEnv is the environment variable name that contains the name of the environment where the service is running at.
	EnvSoajsEnv = "SOAJS_ENV"

	// EnvDeployManual is the environment variable name that indicates if the service has been deployed manually or not.
	EnvDeployManual = "SOAJS_DEPLOY_MANUAL"
)
