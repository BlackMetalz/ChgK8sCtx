package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var exportCmd = &cobra.Command{
	Use:   "export <context-name>",
	Short: "Export a context to stdout (use: ctx export dev > dev.yaml)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := getKubeconfigPath()
		config, _ := loadConfig(path)

		exported := exportContext(config, args[0])
		if exported == nil {
			return
		}

		data, _ := yaml.Marshal(exported)
		fmt.Println(string(data))
	},
}

var importCmd = &cobra.Command{
	Use:   "import <file-path>",
	Short: "Import a context from file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] = file path
		importContext(args[0])
	},
}

func init() {
	ctxCmd.AddCommand(exportCmd)
	ctxCmd.AddCommand(importCmd)
}
