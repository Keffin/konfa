package cmd

import (
	"fmt"
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
		Use:   "konfa CLI",
		Short: "A CLI for incident creation.",
		Long: `
		konfa is a CLI tool created to swiftly be able to edit and re-deploy changes to your kubernetes setup. 
		The idea behind this is to work as a sort of nice-to-have tool for testing out incident scenarios, where configuration breaks.
		Educational purpose. 
		`,
	}
	kubeConfigPath string
	kubeClient     *kubernetes.Clientset
	konfaClient    *client.Client
)

// Files
var (
	namespaceFile  string
	deploymentFile string
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

	konfaDir := filepath.Join(usrHome, ".konfa")
	err = os.MkdirAll(konfaDir, 0700)
	if err != nil {
		log.Printf("Error creating directory: %v \n", err)
		os.Exit(1)
	}
	// Store namespace & deployment file under ~/.konfa/namespace, ~/.konfa/deployment
	namespaceFile = filepath.Join(konfaDir, "namespace")
	deploymentFile = filepath.Join(konfaDir, "deployment")

	// Open or create the namespace file
	nsFile, err := os.OpenFile(namespaceFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening or creating namespace file: %v\n", err)
		os.Exit(1)
	}
	defer nsFile.Close()

	// Open or create the deployment file
	deployFile, err := os.OpenFile(deploymentFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening or creating deployment file: %v\n", err)
		os.Exit(1)
	}
	defer deployFile.Close()

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of konfa.",
	Run: func(cmd *cobra.Command, args []string) {
		colorGreen := "\033[32m"
		fmt.Println(string(colorGreen), "konfa command line tool v0.1")
	},
}
