package kube

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func SetupKubernetesClient() kubernetes.Interface {
	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = getClientOutOfCluster()
	} else {
		kubeClient = getClientInCluster()
	}

	return kubeClient
}

func getClientInCluster() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Cannot get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("Cannot create new kubernetes client from config: %v", err)
	}

	return clientset
}
