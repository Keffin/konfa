package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	//namespace := "kevdev"
	//deploymentName := "firstdeployment"
	//replicas := int32(2)
	//UpdateReplicas(deploymentName, namespace, replicas, *clientset)

}

// TODO
func UpdateConfigMap() {

}

func UpdateReplicas(deploymentName, namespace string, replicaNum int32, client kubernetes.Clientset) {
	d, err := client.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})	

	if err != nil {
		os.Exit(1)
	}

	d.Spec.Replicas = &replicaNum
	ud, err := client.AppsV1().Deployments(namespace).Update( context.Background(),d, metav1.UpdateOptions{})
	
	if err != nil {
		log.Fatalf("Error updating deployment %v", err)
	}
	
	log.Printf("Deployment %s updated: replicas = %d \n", ud.Name, *ud.Spec.Replicas)
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

func FetchMemory(pods *v1.PodList) {
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			memoryLimit := container.Resources.Limits.Memory()
			memoryRequest := container.Resources.Requests.Memory()

			if strings.Contains(pod.Name, "firstdeployment") {
				log.Printf("Pod name: %v\n", pod.Name)
				log.Printf("Container name: %v\n", container.Name)
				log.Printf("Memory limit: %v\n", memoryLimit)
				log.Printf("Memory request: %v\n", memoryRequest)
				log.Println()
			}
		}

	}
}
