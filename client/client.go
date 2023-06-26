package client

import (
	"context"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
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

func updateConfigMapFileProperty(key, newVal, namespace, cname string, client kubernetes.Clientset) {
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