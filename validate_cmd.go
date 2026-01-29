package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate orphaned context/user/cluster in kubeconfig file",
	Run: func(cmd *cobra.Command, args []string) {
		// DEFINE CONFIG
		path, err := getKubeconfigPath()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config, err := loadConfig(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		errors := validateConfig(config)

		if len(errors) > 0 {
			for _, e := range errors {
				fmt.Println(e)
			}
			os.Exit(1)
		} else {
			fmt.Println("No orphaned context/user/cluster found")
		}
	},
}

func init() {
	ctxCmd.AddCommand(validateCmd)
}

func validateConfig(config *KubeConfig) []string {
	var errors []string

	// Check each context, then loop each user/cluster
	for _, ctx := range config.Contexts {
		// Check user exists
		if !itemExists(config, "user", ctx.Context.User) {
			errors = append(errors, fmt.Sprintf("Context %s references non-existent user", ctx.Context.User))
		}

		// Check cluster exists
		if !itemExists(config, "cluster", ctx.Context.Cluster) {
			errors = append(errors, fmt.Sprintf("Context %s references non-existent cluster", ctx.Context.Cluster))
		}
	}

	// Check current-context
	if !itemExists(config, "context", config.CurrentContext) {
		errors = append(errors, fmt.Sprintf("Current context %s references non-existent context", config.CurrentContext))
	}

	return errors
}
