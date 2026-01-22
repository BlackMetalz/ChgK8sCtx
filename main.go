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

	// Call fucking function from config.go, same main package
	config, err := loadConfig("./testdata/kubeconfig")
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
			Items: []string{"Change context", "Change namespace"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Error running promptui:", err)
			return
		}

		// fmt.Println("Your selection: ", result)

		if result == "Change context" {
			action = "ctx"
		} else if result == "Change namespace" {
			action = "ns"
		}

	} else {
		action = os.Args[1] // Get first fucking arg.
	}

	switch action {
	case "ctx":
		fmt.Println("Change context")
		err = switchContext(config)
		if err != nil {
			fmt.Println("Error running switchContext:", err)
			return
		}
		return
	case "ns":
		fmt.Println("Change namespace")
		err = switchNamespace(config)
		if err != nil {
			fmt.Println("Error running switchNamespace:", err)
			return
		}
		return
	default:
		fmt.Println("Unknown action. Use 'ctx' for change context or 'ns' for change namespace")
	}

}
