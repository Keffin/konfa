package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configUpdate = &cobra.Command{
	Use:   "config <key> <value>",
	Short: "Update the configmaps key value",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Please provide both key and value, example: config <key> <value>")
			return
		}
		fmt.Println(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(configUpdate)
}
