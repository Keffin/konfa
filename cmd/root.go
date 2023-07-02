package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Keffin/konfa/client"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Konfa CLI",
		Short: "A CLI for incident creation.",
		Long: `
		Konfa is a CLI tool created to swiftly be able to edit and re-deploy changes to your kubernetes setup. 
		The idea behind this is to work as a sort of nice-to-have tool for testing out incident scenarios, where configuration breaks.
		Educational purpose. 
		`,
	}
	kubeConfigPath string
	kubeClient     *kubernetes.Clientset
	konfaClient    *client.Client
	namespaceFile  string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	usrHome, err := os.UserHomeDir()
	if err != nil {
		log.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}

	kubeConfigPath = filepath.Join(usrHome, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	kubeClient, err = kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		log.Printf("error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(getImage)

	// Store namespace file under ~/.konfa/namespace
	namespaceFile = filepath.Join(usrHome, ".konfa", "namespace")
	err = os.MkdirAll(filepath.Dir(namespaceFile), 0700)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Konfa.",
	Run: func(cmd *cobra.Command, args []string) {
		colorGreen := "\033[32m"
		fmt.Println(string(colorGreen), "Konfa command line tool v0.1")
	},
}

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Set the namespace for your kubernetes client.",
	Long: `For allowing Konfa to operate, you will have to setup a namespace on which it will run its command against.
	Setting up no namespace, will for now not allow it to run.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a namespace.")
			return
		}
		namespace := args[0]
		err := ioutil.WriteFile(namespaceFile, []byte(namespace), 0644)

		if err != nil {
			log.Printf("Error writing namespace to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Setting namespace to: %v \n", namespace)
		konfaClient = client.New(namespace, *kubeClient)
	},
}
