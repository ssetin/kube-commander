package kube

import (
	"github.com/AnatolyRugalev/kube-commander/internal/cfg"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmd "k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

type KubeClient struct {
	*kubernetes.Clientset
}

var config = &struct {
	Path    string `mapstructure:"kube-config"`
	Context string `mapstructure:"kube-context"`
}{}

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")
	cfg.AddPkg(&cfg.Pkg{
		Struct: config,
		PersistentFlags: cfg.FlagsDeclaration{
			"kube-config":  {home + "/.kube/config", "Kubernetes kubeconfig path", "KUBECONFIG"},
			"kube-context": {"", "Kubernetes context to use", "KUBECONTEXT"},
		},
	})
}

func getClientConfig() (*rest.Config, error) {
	return cmd.
		NewNonInteractiveDeferredLoadingClientConfig(
			&cmd.ClientConfigLoadingRules{ExplicitPath: config.Path},
			&cmd.ConfigOverrides{CurrentContext: config.Context},
		).ClientConfig()
}

func GetClient() (*KubeClient, error) {
	c, err := getClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &KubeClient{clientSet}, nil
}
