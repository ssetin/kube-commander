package commander

import "k8s.io/client-go/rest"

type Config interface {
	ClientConfig() (*rest.Config, error)
	Context() string
	Kubeconfig() string
	Namespace() string
}

type ConfigAccessor func() Config
