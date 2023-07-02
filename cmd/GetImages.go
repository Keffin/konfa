package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/Keffin/konfa/client"
	"github.com/spf13/cobra"
)

var getImage = &cobra.Command{
	Use:   "get",
	Short: "Get the deployments for setup namespace",
	Run: func(cmd *cobra.Command, args []string) {
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
		konfaClient.GetDeploymentImages(namespace, deployment)
	},
}

func init() {
	rootCmd.AddCommand(getImage)
}
