// For ctx subcommand handler
package main

import (
	"fmt"
	"strings"

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
	Use:           "ctx [context-name]",
	Short:         "Switch kubernetes context",
	Long:          `Switch to a different kubernetes context. If no context name provided, show interactive selection.`,
	SilenceUsage:  true, // Don't show usage when error
	SilenceErrors: true, // Don't show error
	// Run when user type 'chg-k8s-ctx ctx'
	Args: cobra.MaximumNArgs(2), // 0 or 2 arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		if debugMode {
			fmt.Println("ctx args: ", args)
		}

		path, err := getKubeconfigPath()
		if err != nil {
			return err
		}

		config, err := loadConfig(path)
		if err != nil {
			return err
		}

		// --list flag
		if listFlag {
			// fmt.Println("List all contexts")
			listContexts(config)
			return nil
		}

		// --current/-c flag
		if currentFlag {
			showCurrentContext(config)
			return nil
		}

		// rename mode
		if renameFlag {
			// direct mode. Look stupid but simple xD
			if len(args) == 2 {
				if err := renameContext(config, path, args[0], args[1]); err != nil {
					return err
				}
				// interactive mode
			} else if len(args) == 0 {
				if err := renameContext(config, path); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Unsupported arguments, only 2 arguments are supported")
			}

			return nil // without return, it will run switchContext below
			// if dont need to return if we put renameFlag condition into if-else below
			// But it is not clear i guess.
		}

		// --delete/x flag
		if deleteFlag {
			return deleteContext(config, path, args...)
		}

		// For deleteUserFlag and deleteClusterFlag
		if deleteUserFlag {
			return deleteUser(config, path, args...)
		}

		if deleteClusterFlag {
			// deleteCluster(config, path, args...)
			// LOL, we literal copy deleteUser function xD. But could use same function aswell without small change!
			// Holy fucking shit, opportunity for refactoring again!
			return deleteCluster(config, path, args...)
		}

		// --cleanup flag
		if cleanupFlag {
			return deleteOrphanData(config, path)
		}

		if len(args) == 0 {
			// Interactive mode
			return switchContext(config, path)
		} else if len(args) == 1 && args[0] == "-" {
			// Load previous context
			prevCtx, err := loadPreviousContext()
			if err != nil {
				return err
			}
			// Switch to previous context
			prevCtx = strings.TrimSpace(prevCtx) // Remove trailing \n
			// Don't forget to save previous context
			oldCtx := config.CurrentContext
			savePreviousContext(oldCtx)

			// Update struct for current context
			config.CurrentContext = prevCtx

			// Save config
			err = saveConfig(path, config)
			if err != nil {
				return err
			}

			fmt.Printf("Switched to context %s\n", green(prevCtx))

		} else if len(args) == 1 {
			// Direct switch: ctx ctx-name
			targetCtx := args[0]

			// Check if current context first.
			if targetCtx == config.CurrentContext {
				return fmt.Errorf("You are already on context %s", yellow(targetCtx))
			}

			// Replace item exists with additional fuzzy search
			// if !itemExists(config, "context", targetCtx) {
			// 	fmt.Printf("Context %s does not exist\n", red(targetCtx))
			// 	return
			// }

			// Check exact match
			if itemExists(config, "context", targetCtx) {

				// Save previous context BEFORE switch
				// Without this shit, we can't switch back to previous context
				savePreviousContext(config.CurrentContext)

				// Update struct
				config.CurrentContext = targetCtx
				// Save config
				err := saveConfig(path, config)
				if err != nil {
					return err
				}

				fmt.Printf("Switched to context %s\n", green(targetCtx))
			} else {
				// Fuzzy search handle here.
				matches := fuzzyFindContext(config, targetCtx)

				// Switch like Go idiom xD
				// Count matches from fuzzy search to handle different cases
				switch len(matches) {
				case 0:
					return fmt.Errorf("Context %s does not exist", red(targetCtx))
				case 1:
					targetCtx = matches[0] // use single match to match context
					// Because we only able to choose one context, we can just use the first match

					// Save previous context BEFORE switch
					savePreviousContext(config.CurrentContext)

					// Update struct
					config.CurrentContext = targetCtx
					fmt.Printf("Fuzzy matched: %s\n", green(targetCtx))
				default:
					return fmt.Errorf("Multiple matches for '%s' : %v", targetCtx, matches)
				}

				// Save config
				err := saveConfig(path, config)
				if err != nil {
					return err
				}

				fmt.Printf("Switched to context %s\n", green(targetCtx))
			}

		} else {
			// exception xD
			return fmt.Errorf("Too many arguments")
		}

		return nil

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
