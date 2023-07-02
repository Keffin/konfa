package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Keffin/konfa/cmd"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cmd.Execute()
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
