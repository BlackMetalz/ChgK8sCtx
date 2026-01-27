// For ctx subcommand handler
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Add rename flag
var renameFlag bool

var ctxCmd = &cobra.Command{
	Use:   "ctx [context-name]",
	Short: "Switch kubernetes context",
	Long:  `Switch to a different kubernetes context. If no context name provided, show interactive selection.`,

	// Run when user type 'chg-k8s-ctx ctx'
	Args: cobra.MaximumNArgs(2), // 0 or 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ctx args: ", args)
		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)

		// rename mode
		if renameFlag {
			if len(args) == 2 {
				renameContext(config, path, args[0], args[1]) // direct mode. Look stupid but simple xD
			} else if len(args) == 0 {
				renameContext(config, path) // interactive mode
			} else {
				fmt.Println("Unsupported arguments, only 2 arguments are supported")
			}

			return // without return, it will run switchContext below
			// if dont need to return if we put renameFlag condition into if-else below
			// But it is not clear i guess.
		}

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
	// Rename flag here
	// Param for BoolVar: pointer to variable, name of flag, default value, description
	ctxCmd.Flags().BoolVar(&renameFlag, "rename", false, "Rename context")

	// Add ctx command to root command
	rootCmd.AddCommand(ctxCmd)
}
