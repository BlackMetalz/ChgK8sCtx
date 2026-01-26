package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
	Use:   "ns [namespace-name]",
	Short: "Switch kubernetes namespace",
	Long:  `Switch to a different kubernetes namespace. If no namespace name provided, show interactive selection.`,
	Args:  cobra.MaximumNArgs(1), // Set max args to 1
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ns args: ", args)

		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)

		if len(args) == 0 {
			// Interactive mode
			switchNamespace(config, path)
		} else if len(args) > 0 && args[0] == "kube-system" {
			fmt.Println("You can't switch to kube-system namespace")
		} else {
			// Direct mode
			fmt.Println("TODO: Direct mode")
		}

	},
}

// Add nsCmd to rootCmd
func init() {
	rootCmd.AddCommand(nsCmd)
}
