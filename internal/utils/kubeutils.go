package utils

import (
	"flag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/utils/env"
	"path/filepath"
)

func IsInCluster() bool {
	return env.GetString("KUBERNETES_SERVICE_HOST", "") != ""
}

func CreateKubernetesConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	if IsInCluster() {
		config, err = rest.InClusterConfig()
	} else {
		var kubeConfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}

	if err != nil {
		return nil, err
	}
	return config, nil
}
