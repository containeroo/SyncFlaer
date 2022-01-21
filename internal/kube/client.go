package kube

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// CreateKubernetesClient returns a k8s clientset
func CreateKubernetesClient() kubernetes.Interface {
	var kubeClient kubernetes.Interface
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = getClientOutOfCluster()
	} else {
		kubeClient = getClientInCluster(config)
	}

	return kubeClient
}

func getClientInCluster(config *rest.Config) kubernetes.Interface {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kube client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Cannot get kube config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("Cannot create new kube client from config: %v", err)
	}

	return clientset
}
