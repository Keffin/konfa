package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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
	cmname := "myconfig"
	//UpdateConfigMapKeyProperty("", "testingNewBump", namespace, cmname, *clientset)
	UpdateConfigMap(true, "nestlevel2", "dyn nest bump", namespace, cmname, *clientset)
	UpdateConfigMap(false, "file_data", "dyn bump", namespace, cmname, *clientset)
	//deploymentName := "firstdeployment"
	//replicas := int32(2)
	//UpdateReplicas(deploymentName, namespace, replicas, *clientset)

}

// TODO
// input: ConfigMap
// output: None
// Takes a configmap, verifies whether key belongs to file or key prop, updates accordingly.
func UpdateConfigMap(isFileProp bool, key, val, namespace, configName string, client kubernetes.Clientset) {

	if isFileProp {
		UpdateConfigMapFileProperty(key, val, namespace, configName, client)
	} else {
		UpdateConfigMapKeyProperty(key, val, namespace, configName, client)
	}

}

func UpdateConfigMapFileProperty(key, newVal, namespace, cname string, client kubernetes.Clientset) {
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), cname, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting configmap: %v \n", err)
		os.Exit(1)
	}
	for k, v := range cm.Data {
		if isFileConf(k) {
			if strings.HasSuffix(k, ".yaml") {

				var yRes interface{}
				err = yaml.Unmarshal([]byte(v), &yRes)
				if err != nil {
					log.Printf("Error parsing YAML for key %s: %v\n", k, err)
					continue
				}

				updateValueNested(yRes, key, newVal)

				updatedData, err := yaml.Marshal(yRes)

				if err != nil {
					log.Printf("Error marshaling YAML for key %s: %v\n", k, err)
					continue
				}

				cm.Data[k] = string(updatedData)

			}
		}
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating ConfigMap: %v\n", err)
		os.Exit(1)
	}

	log.Println("Nested configmap value successfully updated.")
}

func updateValueNested(data interface{}, key string, value interface{}) {
	switch d := data.(type) {
	case map[string]interface{}:
		if _, ok := d[key]; ok {
			d[key] = value
		} else {
			for _, v := range d {
				updateValueNested(v, key, value)
			}
		}
	case []interface{}:
		for _, v := range d {
			updateValueNested(v, key, value)
		}
	}
}

// We want the user to supply a key as well, so we know what we are searching for and only update that one.
func UpdateConfigMapKeyProperty(key, newVal, namespace, cmname string, client kubernetes.Clientset) {
	log.Println("Updating key property in configmap")
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), cmname, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting configmap: %v \n", err)
		os.Exit(1)
	}
	for k, _ := range cm.Data {
		// Match that the key is also correct
		if !isFileConf(k) {
			cm.Data[k] = newVal
			_, err := client.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Error updating configmap: %v\n", err)
				os.Exit(1)
			}
			log.Println("Successfully updated regular key prop in configmap")
		} else {
			continue
		}
	}
}

// Currently only yaml, yml supported. Properties and JSON is TODO.
func isFileConf(conf string) bool {
	//`^.*\.(yaml|yml|properties|txt|json)$`
	fileExtensions := []string{".yaml", ".yml", ".txt", ".properties", ".json"}

	for _, ext := range fileExtensions {
		if strings.HasSuffix(conf, ext) {
			return true
		}
	}
	return false
}

func UpdateReplicas(deploymentName, namespace string, replicaNum int32, client kubernetes.Clientset) {
	d, err := client.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})

	if err != nil {
		os.Exit(1)
	}

	d.Spec.Replicas = &replicaNum
	ud, err := client.AppsV1().Deployments(namespace).Update(context.Background(), d, metav1.UpdateOptions{})

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
