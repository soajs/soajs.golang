package soajsgo

const (
	// EnvProfile is the environment variable name that contains location of the profile
	// to use so that SOAJS can connect to the core database.
	EnvProfile = "SOAJS_PROFILE"

	// EnvSRVIP is optional environment variable used to specify which IP address to use
	// if the machine has more than one active interface.
	EnvSRVIP = "SOAJS_SRVIP"

	// EnvSOLO is optional environment variable used to launch any service on top of SOAJS
	// without the need of a database.
	EnvSOLO = "SOAJS_SOLO"

	// EnvSrvAutoRegisterHost is optional environment variable used in case a service should register itself or not.
	EnvSrvAutoRegisterHost = "SOAJS_SRV_AUTOREGISTERHOST"

	// EnvDaemonGRPConf is the environment variable name that contains the name of the daemon group to use;
	// available for daemons ONLY.
	EnvDaemonGRPConf = "SOAJS_DAEMON_GRP_CONF"

	// EnvGCName is the environment variable name that contains mandatory variable if deploying a GCS service
	// and contains the name of that GCS service.
	EnvGCName = "SOAJS_GC_NAME"

	// EnvDCVersion is the environment variable name that contains mandatory variable if deploying a GCS service
	// and contains the version of that GCS service.
	EnvDCVersion = "SOAJS_GC_VERSION"

	// EnvGCMaxUploadLimit is optional variable if deploying a GCS service that specifies the maximum upload limit
	// of file sizes to accept.
	EnvGCMaxUploadLimit = "SOAJS_GC_MAX_UPLOAD_LIMIT"

	// EnvRegistryAPIAddress is the environment variable name that contains the IP address and port of
	// the controller service that runs in the same environment. The SOAJS middleware uses this variable to fetch
	// the registry of this environment and supply it to your service.
	EnvRegistryAPIAddress = "SOAJS_REGISTRY_API"

	// EnvSoajsEnv is the environment variable name that contains the name of the environment where the service is running at.
	EnvSoajsEnv = "SOAJS_ENV"

	// EnvDeployManual is the environment variable name that indicates if the service has been deployed manually or not.
	EnvDeployManual = "SOAJS_DEPLOY_MANUAL"
)
