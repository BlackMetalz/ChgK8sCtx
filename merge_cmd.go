package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	outputFlag string // --output flag
)

var mergeCmd = &cobra.Command{
	Use:   "merge <source1> <source2>",
	Short: "Merge 2 kubeconfig file",
	Args:  cobra.ExactArgs(2), // exactly 2 args
	Run: func(cmd *cobra.Command, args []string) {

		// Define 2 source files
		file1, file2 := args[0], args[1]

		// Load 2 source files
		config1, err := loadConfig(file1)
		if err != nil {
			fmt.Println("Error loading source1:", err)
			return
		}

		config2, err := loadConfig(file2)
		if err != nil {
			fmt.Println("Error loading source2:", err)
			return
		}

		// Merge configs
		mergeConfigs(config2, config1) // merge config2 into config1

		// Save to output file
		err = writeConfig(outputFlag, config1) // save merged config to output file
		if err != nil {
			fmt.Println("Error saving output:", err)
			return
		}

		fmt.Printf("Merged successfully to: %s\n", outputFlag)
	},
}

func init() {
	mergeCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file after merged")
	mergeCmd.MarkFlagRequired("output") // make it required

	// This register sub command to ctx command
	ctxCmd.AddCommand(mergeCmd) // because this is sub command of ctx
}

func mergeConfigs(source, target *KubeConfig) {
	// Merge cluster from source --> target
	for _, c := range source.Clusters {
		if !itemExists(target, "cluster", c.Name) {
			target.Clusters = append(target.Clusters, c)
		} else {
			fmt.Printf("Warning: Cluster %s already exists, skipping\n", c.Name)
		}
	}

	// Merge user from source --> target
	for _, u := range source.Users {
		if !itemExists(target, "user", u.Name) {
			target.Users = append(target.Users, u)
		} else {
			fmt.Printf("Warning: User %s already exists, skipping\n", u.Name)
		}
	}

	// Merge context from source --> target
	for _, ctx := range source.Contexts {
		if !itemExists(target, "context", ctx.Name) {
			target.Contexts = append(target.Contexts, ctx)
		} else {
			fmt.Printf("Warning: Context %s already exists, skipping\n", ctx.Name)
		}
	}
}
