// For ctx subcommand handler
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ctxCmd = &cobra.Command{
	Use:   "ctx [context-name]",
	Short: "Switch kubernetes context",
	Long:  `Switch to a different kubernetes context. If no context name provided, show interactive selection.`,

	// Run when user type 'chg-k8s-ctx ctx'
	Args: cobra.MaximumNArgs(1), // 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ctx args: ", args)
		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)
		if len(args) == 0 {
			// Interactive mode
			switchContext(config, path)
		} else {
			// Direct mode
			fmt.Println("TODO: Direct mode: ", args[0])
		}
	},
}

func init() {
	// Add ctx command to root command
	rootCmd.AddCommand(ctxCmd)
}
