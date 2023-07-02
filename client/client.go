package client

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	client    kubernetes.Clientset
	namespace string
}

func New(namespace string, client kubernetes.Clientset) *Client {

	return &Client{
		client:    client,
		namespace: namespace,
	}
}

// TODO
// input: ConfigMap
// output: None
// Takes a configmap, verifies whether key belongs to file or key prop, updates accordingly.
func (c *Client) UpdateConfigMap(isFileProp bool, key, val, configName string) {

	if isFileProp {
		updateConfigMapFileProperty(key, val, c.namespace, configName, c.client)
	} else {
		updateConfigMapKeyProperty(key, val, c.namespace, configName, c.client)
	}

}

func (c *Client) GetNamespace() string {
	return c.namespace
}

// Deployments use the AppsV1(), while config etc uses CoreV1()
// and for Spec content.
// Spec content that can be updated:
// Replicas
// Resources, this will be issue since you can have multiple images with diff resources, so dupe keys.
// Ports

// Allow something like
// Key: ports.ContainerPort, newContainerKey: 81
// Key: resources.Limits.memory, newContainerKey: 100Mi
func (c *Client) UpdateContainerSpecs(containerName, containerKey, newContainerKey, namespace, deployment string) {
	containers := c.GetDeploymentImages(namespace, deployment)
	d, err := c.client.AppsV1().Deployments(namespace).Get(context.Background(), deployment, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get deployment: %v", err)
	}

	var hasBeenUpdated bool

	if contains(containers, containerName) {
		for idx, container := range d.Spec.Template.Spec.Containers {
			if container.Name == containerName {
				if containerKey == "ports.ContainerPort" {
					newPortNum, err := strconv.Atoi(newContainerKey)
					if err != nil {
						log.Fatalf("Error when converting container port num to int: %v", err)
					}
					for pdx := range container.Ports {
						d.Spec.Template.Spec.Containers[idx].Ports[pdx].ContainerPort = int32(newPortNum)
					}
					hasBeenUpdated = true
				} else if containerKey == "resources.limits" || containerKey == "resources.requests" {
					r := &container.Resources
					newResourceValue := resource.MustParse(newContainerKey)
					switch containerKey {
					case "resources.limits":
						r.Limits[corev1.ResourceMemory] = newResourceValue
					case "resources.requests":
						r.Requests[corev1.ResourceMemory] = newResourceValue
					}
					hasBeenUpdated = true
				}
				break
			}
		}
	}

	if hasBeenUpdated {
		updated, err := c.client.AppsV1().Deployments(namespace).Update(context.Background(), d, metav1.UpdateOptions{})
		if err != nil {
			log.Fatalf("Error updating deployment: %v", err)
		}

		log.Printf("Deployment updated: %s", updated.Name)
	}
}

func (c *Client) UpdateReplicas(deploymentName, namespace string, replicaNum int32) {
	d, err := c.client.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})

	if err != nil {
		os.Exit(1)
	}

	d.Spec.Replicas = &replicaNum
	ud, err := c.client.AppsV1().Deployments(namespace).Update(context.Background(), d, metav1.UpdateOptions{})

	if err != nil {
		log.Fatalf("Error updating deployment %v", err)
	}

	log.Printf("Deployment %s updated: replicas = %d \n", ud.Name, *ud.Spec.Replicas)
}

func (c *Client) GetDeploymentImages(namespace, deployment string) []string {
	d, err := c.client.AppsV1().Deployments(namespace).Get(context.Background(), deployment, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to Get deployment: %v\n", err)
		os.Exit(1)
	}

	images := make([]string, 0)

	for _, v := range d.Spec.Template.Spec.Containers {
		images = append(images, v.Name)
	}

	log.Println(images)
	return images
}

func updateConfigMapFileProperty(key, newVal, namespace, cname string, client kubernetes.Clientset) {
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), cname, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting configmap: %v \n", err)
		os.Exit(1)
	}
	parents := make([]string, 0)
	for k, v := range cm.Data {
		if isFileConf(k) {
			if strings.HasSuffix(k, ".yaml") {

				var yRes interface{}
				err = yaml.Unmarshal([]byte(v), &yRes)
				if err != nil {
					log.Printf("Error parsing YAML for key %s: %v\n", k, err)
					continue
				}

				updateValueNested(yRes, key, newVal, &parents)

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
	//parents = append(parents, key)
	log.Println("Nested configmap value successfully updated.")
	log.Println(parents)
}

// We want the user to supply a key as well, so we know what we are searching for and only update that one.
func updateConfigMapKeyProperty(key, newVal, namespace, cmname string, client kubernetes.Clientset) {
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

func updateValueNested(data interface{}, key string, value interface{}, p *[]string) {
	switch d := data.(type) {
	case map[string]interface{}:
		if _, ok := d[key]; ok {
			d[key] = value
		} else {
			for k, v := range d {
				*p = append(*p, k)
				updateValueNested(v, key, value, p)
			}
		}
	case []interface{}:
		for _, v := range d {
			updateValueNested(v, key, value, p)
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

func contains(list []string, target string) bool {
	for _, value := range list {
		if value == target {
			return true
		}
	}
	return false
}
