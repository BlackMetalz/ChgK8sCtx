// For ctx subcommand handler
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Add rename flag
// Day5 updated: add listFlag / currentFlag
var (
	renameFlag        bool
	listFlag          bool
	currentFlag       bool
	deleteFlag        bool // Default for delete context xD
	deleteUserFlag    bool
	deleteClusterFlag bool
	cleanupFlag       bool // For clean up orphan user/cluster.
)

var ctxCmd = &cobra.Command{
	Use:   "ctx [context-name]",
	Short: "Switch kubernetes context",
	Long:  `Switch to a different kubernetes context. If no context name provided, show interactive selection.`,

	// Run when user type 'chg-k8s-ctx ctx'
	Args: cobra.MaximumNArgs(2), // 0 or 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("ctx args: ", args) // No need to debug anymoar xD
		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)

		// --list flag
		if listFlag {
			// fmt.Println("List all contexts")
			listContexts(config)
			return
		}

		// --current flag
		if currentFlag {
			showCurrentContext(config)
			return
		}

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

		// --delete/x flag
		if deleteFlag {
			deleteContext(config, path, args...) // Take variadic args
			return
		}

		// For deleteUserFlag and deleteClusterFlag
		if deleteUserFlag {
			deleteUser(config, path, args...)
			return
		}

		if deleteClusterFlag {
			// deleteCluster(config, path, args...)
			// LOL, we literal copy deleteUser function xD. But could use same function aswell without small change!
			// Holy fucking shit, opportunity for refactoring again!
			deleteCluster(config, path, args...)
			return
		}

		// --cleanup flag
		if cleanupFlag {
			deleteOrphanData(config, path)
			return
		}

		if len(args) == 0 {
			// Interactive mode
			switchContext(config, path)
		} else {
			// exception xD
			fmt.Println("This tool does not support arguments like this: ", args[0])
		}
	},
}

func init() {
	// Rename flag here
	// Param for BoolVar: pointer to variable, name of flag, default value, description
	// BoolVar defines a bool flag with specified name, default value, and usage string
	ctxCmd.Flags().BoolVar(&renameFlag, "rename", false, "Rename context")

	// Add list Flag
	// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
	ctxCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all contexts")

	// Add current Flag
	ctxCmd.Flags().BoolVarP(&currentFlag, "current", "c", false, "Show current context")

	// Add delete Flag
	ctxCmd.Flags().BoolVarP(&deleteFlag, "delete", "x", false, "Delete context")

	// Add delete user Flag
	ctxCmd.Flags().BoolVar(&deleteUserFlag, "delete-user", false, "Delete user in kubeconfig")

	// Add delete cluster Flag
	ctxCmd.Flags().BoolVar(&deleteClusterFlag, "delete-cluster", false, "Delete cluster in kubeconfig")

	// Add cleanup flag
	ctxCmd.Flags().BoolVar(&cleanupFlag, "cleanup", false, "Delete orphan user/cluster in kubeconfig")

	// Add ctx command to root command
	rootCmd.AddCommand(ctxCmd)
}
