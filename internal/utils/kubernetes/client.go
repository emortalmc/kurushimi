package kubernetes

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kurushimi/internal/utils"
	"time"
)

func CreateClients() (kubeClient *kubernetes.Clientset, agonesClient *versioned.Clientset) {
	kubeConfig, err := createKubernetesConfig()
	if err != nil {
		panic(err)
	}

	kubeClient = kubernetes.NewForConfigOrDie(kubeConfig)

	// AgonesClient contains the Agones client for creating GameServerAllocation objects
	agonesClient = versioned.NewForConfigOrDie(kubeConfig)

	return kubeClient, agonesClient
}

func createKubernetesConfig() (*rest.Config, error) {
	config, err := utils.CreateKubernetesConfig()
	if err != nil {
		return nil, err
	}

	config.Timeout = time.Second * 5
	return config, nil
}
