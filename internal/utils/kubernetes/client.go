package kubernetes

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kurushimi/internal/utils"
)

var (
	KubeClient kubernetes.Interface

	// AgonesClient contains the Agones client for creating GameServerAllocation objects
	AgonesClient versioned.Interface
)

func Init() {
	kubeConfig, err := createKubernetesConfig()
	if err != nil {
		panic(err)
	}

	KubeClient = kubernetes.NewForConfigOrDie(kubeConfig)

	// AgonesClient contains the Agones client for creating GameServerAllocation objects
	AgonesClient = versioned.NewForConfigOrDie(kubeConfig)
}

func createKubernetesConfig() (*rest.Config, error) {
	kConfig, err := utils.CreateKubernetesConfig()
	if err != nil {
		return nil, err
	}
	return kConfig, nil
}
