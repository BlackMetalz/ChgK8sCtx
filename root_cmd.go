package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debugMode      bool
	kubeconfigFlag string
	version        = "dev" // Injected at build time with -ldflags
)

var rootCmd = &cobra.Command{
	Use:     "chgctx",
	Version: version,
	Short:   "A kubectx/kubens clone written in Go for learning purpose",
	Long:    "ChgK8sCtx is a CLI tool to switch Kubernetes contexts and namespaces easily. My fucking pet project for learning Go",

	// Run ctxCmd when no subcommand is specified
	// args is slice of string, contains all arguments after command name
	RunE: func(cmd *cobra.Command, args []string) error {
		// Route default execution to ctxCmd
		return ctxCmd.RunE(cmd, args)
	},
}

func init() {
	// PersistentFlags returns the persistent FlagSet specifically set in the current command.
	// We need to add this flag to root command so that it can be used by all subcommands.
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode, show debug information xD")
	// Add kubeconfig flag
	rootCmd.PersistentFlags().StringVar(&kubeconfigFlag, "kubeconfig", "", "Path to kubeconfig file specified")
}

// This func will be called by main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
