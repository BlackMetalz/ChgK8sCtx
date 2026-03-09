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

	// Allow any args so Cobra doesn't reject context names like "aws", "gke", etc.
	// Without this, Cobra treats unknown first-args as unknown subcommands and errors.
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,

	// Delegate to ctxCmd logic when no subcommand is matched
	RunE: func(cmd *cobra.Command, args []string) error {
		return ctxCmd.RunE(cmd, args)
	},
}

func init() {
	// PersistentFlags: shared across all subcommands
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode, show debug information xD")
	rootCmd.PersistentFlags().StringVar(&kubeconfigFlag, "kubeconfig", "", "Path to kubeconfig file specified")

	// Mirror all ctx flags on root so `chgctx --list`, `chgctx -c`, etc. work directly.
	// Same variables as ctxCmd so both `chgctx --list` and `chgctx ctx --list` work.
	rootCmd.Flags().BoolVar(&renameFlag, "rename", false, "Rename context")
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all contexts")
	rootCmd.Flags().BoolVarP(&currentFlag, "current", "c", false, "Show current context")
	rootCmd.Flags().BoolVarP(&deleteFlag, "delete", "x", false, "Delete context")
	rootCmd.Flags().BoolVar(&deleteUserFlag, "delete-user", false, "Delete user in kubeconfig")
	rootCmd.Flags().BoolVar(&deleteClusterFlag, "delete-cluster", false, "Delete cluster in kubeconfig")
	rootCmd.Flags().BoolVar(&cleanupFlag, "cleanup", false, "Delete orphan user/cluster in kubeconfig")
}

// This func will be called by main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
