package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Keffin/konfa/client"
	"github.com/spf13/cobra"
)

var containerUpdate = &cobra.Command{
	Use:   "container",
	Short: "Update the containers values",
	Long: `container is used to update certain container values in your deployment. For now it only supports updating values such as 
	ports.containerPorts, resources.Limits, resources.Requests, must follow this format for now.
	
Usage:
	konfa container <containername> set <key> <value> where key = {ports.containerPort, resources.limits, resources.requests}`,
	Run: executeContainer,
}

func init() {
	rootCmd.AddCommand(containerUpdate)
}

func executeContainer(cmd *cobra.Command, args []string) {
	if len(args) != 4 {
		fmt.Println("Command format faulty. Please provide key and new value. E.g container <containerName> set <key> <value>")
		return
	}

	data, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		log.Printf("Error reading namespace file: %v\n", err)
		os.Exit(1)
	}
	d, err := ioutil.ReadFile(deploymentFile)
	if err != nil {
		log.Printf("Error reading deployment file: %v\n", err)
		os.Exit(1)
	}

	namespace := string(data)
	deployment := string(d)
	konfaClient = client.New(namespace, *kubeClient)
	konfaClient.UpdateContainerSpecs(args[0], args[2], args[3], namespace, deployment)
}
