package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
	Use:   "ns [namespace]",
	Short: "Switch kubernetes namespace",
	Long: `Switch to a different kubernetes namespace. 
Examples:
  # Interactive mode - show selection menu
  kctx ns
  
  # Direct mode - switch immediately  
  kctx ns kube-public
  kctx ns my-namespace`,
	// Add Example field (optional)
	Example: `  kctx ns
  kctx ns kube-public`,
	Args: cobra.MaximumNArgs(1), // Set max args to 1
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			fmt.Println("ns args: ", args)
		}

		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)

		// Option C: warn when a context name conflicts with this subcommand name.
		// e.g. if user has a context named "ns", `kctx ns` routes here instead.
		// Suggest using `kctx ctx ns` as explicit escape hatch.
		if config != nil && itemExists(config, "context", cmd.Name()) {
			fmt.Printf(yellow("⚠ Hint:")+" context '%s' exists. To switch context, use: kctx ctx %s\n\n",
				cmd.Name(), cmd.Name())
		}

		if len(args) == 0 {
			// Interactive mode
			switchNamespace(config, path)
		} else if len(args) > 0 && args[0] == "kube-system" {
			fmt.Println(red("You can't switch to kube-system namespace"))
		} else {
			// Direct mode. Set namespace without asking for selection
			targetNS := args[0]

			// Get current context entry
			entry := getCurrentContextEntry(config)
			if entry == nil {
				fmt.Println("Error getting current context entry")
				return
			}

			// Check namespace is currently selected
			if entry.Context.Namespace == targetNS {
				fmt.Println(red("Namespace is already selected"))
				return
			}

			// Update and save
			entry.Context.Namespace = targetNS
			err := saveConfig(path, config)
			if err != nil {
				fmt.Println("Error saving config: ", err)
				return
			}

			fmt.Println("Namespace switched to: ", targetNS)

		}

	},
}

// Add nsCmd to rootCmd
func init() {
	rootCmd.AddCommand(nsCmd)
}
