package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   `deployment`,
	Short: "Set or get the deployment for your kubernetes client.",
	Long: `For allowing Konfa to operate, you will have to setup a deployment on which it will run its deployment specific command against.
	Setting up no deployment, will for now cause the program to not be able to run certain deployment edits.
	
Usage:
	konfa deployment get
	konfa deployment set <name>
	
	`,
	Run: fetchAndUpdateDeployment,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
}

func fetchAndUpdateDeployment(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a deployment.")
		return
	}

	if args[0] == "get" {
		data, err := ioutil.ReadFile(deploymentFile)
		if err != nil {
			log.Printf("Error opening the deployment config: %v \n", err)
			return
		}
		ns := string(data)
		log.Printf("Current deployment is: %s \n", ns)
	}

	if len(args) == 2 && args[0] == "set" {
		ns := args[1]
		err := ioutil.WriteFile(deploymentFile, []byte(ns), 0644)
		if err != nil {
			log.Printf("Error writing deployment to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Setting deployment to: %v \n", ns)
	}
}
