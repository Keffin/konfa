package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Keffin/konfa/client"
	"github.com/spf13/cobra"
)

var configUpdate = &cobra.Command{
	Use:   "config",
	Short: "Update the configmaps key value, or fetch a configmaps key value pair.",
	Long: `For fetching and/or updating configmap values without having to run kubectl edit.
	
Usage:
	konfa config get <key>
	konfa config set <key> <value>
	`,
	Run: executeConfig,
}

func init() {
	rootCmd.AddCommand(configUpdate)
	rootCmd.PersistentFlags().Bool("is-file", false, "Add whether the entry to change is a file property or key property")
	rootCmd.PersistentFlags().String("config", "", "Add name of configmap you wish to change")
}

func executeConfig(cmd *cobra.Command, args []string) {
	if len(args) != 3 && len(args) != 2 {
		fmt.Println("Command format faulty. Please provide either get <key> or set <key> <value>")
	}

	data, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		log.Printf("Error opening the namespace config: %v \n", err)
		return
	}
	ns := string(data)
	konfaClient = client.New(ns, *kubeClient)

	// Run set scenario
	if len(args) == 3 {
		file, err := cmd.Flags().GetBool("is-file")
		if err != nil {
			fmt.Println("Flag is-file must be of type bool")
		}
		if file {
			configName, err := cmd.Flags().GetString("config")
			if err != nil {
				fmt.Println("No config name supplied, exiting.")
				return
			}
			konfaClient.UpdateConfigMap(file, args[1], args[2], configName)
		} else {
			configName, err := cmd.Flags().GetString("config")
			if err != nil {
				fmt.Println("No config name supplied, exiting.")
				return
			}

			konfaClient.UpdateConfigMap(file, args[1], args[2], configName)
		}
	} else {
		configName, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("No config name supplied, exiting.")
			return
		}
		v, err := konfaClient.GetConfigMapValue(configName, args[1])
		if err != nil {
			log.Printf("Error getting value from ConfigMap: %v\n", err)
			return
		}
		log.Printf("Value of key '%s' in ConfigMap '%s' is: %s\n", args[1], configName, v)
	}
}
