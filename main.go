package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Keffin/konfa/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		log.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}

	kubeConfigPath := filepath.Join(usrHome, ".kube", "config")
	log.Printf("Using kubeConfigPath: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		log.Printf("error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	namespace := "kevdev"
	//cmname := "myconfig"
	deploymentName := "firstdeployment"
	m := client.New(namespace, *clientset)
	//m.UpdateConfigMap(true, "nestlevel2", "struct nest bump", namespace, cmname)
	//m.UpdateConfigMap(false, "file_data", "struct key prop bump", namespace, cmname)
	//m.UpdateDeployment("containers.nginx.requests.memory", "100Mi", namespace, deploymentName)
	//m.UpdateDeployment("containers.nginx.limits.memory", "200Mi", namespace, deploymentName)
	//i := m.GetDeploymentImages(namespace, deploymentName)
	//m.UpdateContainer("nginx", "ports.ContainerPort", "70", namespace, deploymentName)
	m.UpdateContainerSpecs("nginx", "resources.requests", "300Mi", namespace, deploymentName)
	//m.UpdateContainerSpecs("nginx-2", "resources.limits", "600Mi", namespace, deploymentName)
	//deploymentName := "firstdeployment"
	replicas := int32(2)
	//UpdateReplicas(deploymentName, namespace, replicas, *clientset)
	m.UpdateReplicas(deploymentName, namespace, replicas)

}

func ListPods(namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	log.Println("Get kubernetes pods.")
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v", err)
		return nil, err
	}
	return pods, nil
}
