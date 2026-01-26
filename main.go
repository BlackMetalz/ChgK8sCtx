package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui" // Promptui for interactive prompt
	// Import this shit for unmarshal yaml
)

func main() {

	// Example:
	// args[0] = binary name (ex: "./ChgK8sCtx")
	// args[1] = first arg (ex: "ctx" or "ns")

	// Get namespace list
	// nsList, err := listNamespaces()
	// if err != nil {
	// 	fmt.Println("Error getting namespace list:", err)
	// 	return
	// } else {
	// 	fmt.Println("Namespace list:", nsList)
	// }
	// My job done here xD

	// Get kubeconfig path
	kubeconfigPath, err := getKubeconfigPath()
	if err != nil {
		fmt.Println("Error getting kubeconfig path:", err)
		return
	}

	// Call fucking function from config.go, same main package
	config, err := loadConfig(kubeconfigPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Check args
	var action string
	if len(os.Args) < 2 {
		// Ask user to select action
		prompt := promptui.Select{
			Label: "Select action",
			Items: []string{"Change current context", "Change default namespace"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Error running promptui:", err)
			return
		}

		// fmt.Println("Your selection: ", result)

		if result == "Change current context" {
			action = "ctx"
		} else if result == "Change default namespace" {
			action = "ns"
		}

	} else {
		action = os.Args[1] // Get first fucking arg.
	}

	switch action {
	case "ctx":
		// fmt.Println("Change context")
		err = switchContext(config, kubeconfigPath)
		if err != nil {
			fmt.Println("Error running switchContext:", err)
			return
		} else {
			fmt.Println("Successfully switched to context: ", config.CurrentContext)
		}
		return
	case "ns":
		// fmt.Println("Change default namespace")
		err = switchNamespace(config, kubeconfigPath)
		if err != nil {
			fmt.Println("Error running switchNamespace:", err)
			return
		} else {
			entry := getCurrentContextEntry(config) // nil check here is not necessary
			fmt.Println("Successfully switched to namespace: ", entry.Context.Namespace)
		}
		return
	default:
		fmt.Println("Unknown action. Use 'ctx' for change context or 'ns' for change namespace")
	}

}
