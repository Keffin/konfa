package client

import (
	"context"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	client kubernetes.Clientset
	namespace string
}

func New(namespace string, client kubernetes.Clientset) *Client {

	return &Client{
		client: client, 
		namespace: namespace,
	}
}

// TODO
// input: ConfigMap
// output: None
// Takes a configmap, verifies whether key belongs to file or key prop, updates accordingly.
func (c *Client) UpdateConfigMap(isFileProp bool, key, val, namespace, configName string) {

	if isFileProp {
		updateConfigMapFileProperty(key, val, namespace, configName, c.client)
	} else {
		updateConfigMapKeyProperty(key, val, namespace, configName, c.client)
	}

}

// Deployments use the AppsV1(), while config etc uses CoreV1()
// For now only allow update to MetaData -> Namespace 
// and for Spec content.
// Spec content that can be updated:
// Replicas
// Resources, this will be issue since you can have multiple images with diff resources, so dupe keys.
// VolumeMounts
// Volumes
// Ports
func (c *Client) UpdateDeployment(key, val, namespace, deployment string) {
	updateDeploymentProps(key, val, namespace, deployment, c.client)
	
	// If key: MetaData.NameSpace then Update MetaData.NameSpace
	// If key: Spec.Replicas then Update Spec.Replicas.
}

func updateDeploymentProps(key, newVal, ns, deployment string, client kubernetes.Clientset) {
	d, err := client.AppsV1().Deployments(ns).Get(context.Background(), deployment, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error parsing value to int: %v\n", err)
		os.Exit(1)
	}

	if err != nil {
		log.Printf("Failed to Get deployment: %v \n", err)
		os.Exit(1)
	}

	kd := strings.Split(key, ".")
	cpuQ := resource.MustParse(newVal)

	for _, cont := range d.Spec.Template.Spec.Containers {
		if kd[2] == "requests" {
			r := &cont.Resources
			r.Requests[corev1.ResourceMemory] = cpuQ
		} else if kd[2] == "limits" {
			r := &cont.Resources
			r.Limits[corev1.ResourceMemory] = cpuQ
		}
		
	}
	
	updated, err := client.AppsV1().Deployments(ns).Update(context.Background(), d, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating deployment: %v\n", err)
		os.Exit(1)
	}
	

	log.Printf("Deployment updated: %s \n", updated.Name)
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