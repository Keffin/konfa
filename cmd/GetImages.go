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

		namespace := string(data)
		konfaClient = client.New(namespace, *kubeClient)
		konfaClient.GetDeploymentImages(namespace, "firstdeployment")
	},
}

func init() {
	rootCmd.AddCommand(getImage)
}