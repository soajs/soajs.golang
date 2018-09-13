package structs

type Deployer struct {
	Type     string `json:"type"`
	Selected string `json:"selected"`
	Manual   Manual `json:"manual"`
	Container Container `json:"container"`
}

type Manual struct {
  Nodes string `json:"nodes"`
}

type Container struct {
  Docker Docker `json:"docker"`
  Kubernetes Kubernetes `json:"kubernetes"`
}

type Docker struct {
  Local DockerLocal  `json:"local"`
  Remote DockerRemote `json:"remote"`
}

type DockerLocal struct {
  Nodes      string `json:"nodes"`
  SocketPath string `json:"socketPath"`
}

type DockerRemote struct {
  APIPort     int    `json:"apiPort"`
  Nodes       string `json:"nodes"`
  APIProtocol string `json:"apiProtocol"`
  Auth        ContainerAuth `json:"auth"`
}

type Kubernetes struct {
  Local KubernetesLocal `json:"local"`
  Remote KubernetesRemote `json:"remote"`
}

type KubernetesLocal struct {
  Nodes     string `json:"nodes"`
  Namespace KubernetesNamespace `json:"namespace"`
  Auth ContainerAuth `json:"auth"`
}

type KubernetesRemote struct {
  Nodes     string `json:"nodes"`
  Namespace KubernetesNamespace `json:"namespace"`
  Auth ContainerAuth `json:"auth"`
}

type ContainerAuth struct {
  Token string `json:"token"`
}

type KubernetesNamespace struct {
  Default string `json:"defautl"`
  PerService bool `json:"perService"`
}
