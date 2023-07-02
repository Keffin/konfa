package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Keffin/konfa/client"
	"github.com/spf13/cobra"
)

var namespaceCmd = &cobra.Command{
	Use:   `namespace`,
	Short: "Set or get the namespace for your kubernetes client.",
	Long: `For allowing Konfa to operate, you will have to setup a namespace on which it will run its command against.
	Setting up no namespace, will for now not allow it to run.
	
Usage:
	konfa namespace get
	konfa namespace set <name>
	
	`,
	Run: fetchAndUpdateNamespace,
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}

func fetchAndUpdateNamespace(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a namespace.")
		return
	}

	// Run the update
	if args[0] == "get" {
		data, err := ioutil.ReadFile(namespaceFile)
		if err != nil {
			log.Printf("Error opening the namespace config: %v \n", err)
			return
		}
		ns := string(data)
		log.Printf("Current namespace is: %s \n", ns)
	}

	if len(args) == 2 && args[0] == "set" {
		ns := args[1]
		err := ioutil.WriteFile(namespaceFile, []byte(ns), 0644)
		if err != nil {
			log.Printf("Error writing namespace to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Setting namespace to: %v \n", ns)
		konfaClient = client.New(ns, *kubeClient)
	}
}
