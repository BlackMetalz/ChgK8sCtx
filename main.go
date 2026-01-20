package main

import (
	"fmt"
	"os"
	"strings"

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
	case "ns":
		fmt.Println("Change namespace")
	default:
		fmt.Println("Unknown action. Use 'Change context' or 'Change namespace'")
	}

	// List all context and append all cluster name to a slice of string
	// Init default slide of string
	var contextList []string
	// fmt.Println(strings.Repeat("=", 10) + "All context:" + strings.Repeat("=", 10))
	for _, ctx := range config.Contexts {
		// fmt.Println(abc.Name) // No need to print this shit.
		// Append to slice
		// Add a text: (current-context) next to current context to let user know this is current fucking context.
		if ctx.Name == config.CurrentContext {
			// We will need to remove this text when we write back to file in next step.
			contextList = append(contextList, ctx.Name+"(current-context)")
		} else {
			contextList = append(contextList, ctx.Name)
		}
	}

	// Ok, we have fucking promptui
	prompt := promptui.Select{
		Label: "Select k8s cluster",
		Items: contextList,
	}

	// Run fucking prompt
	_, result, err := prompt.Run()

	// Go idiom style
	switch err {
	case nil:
		// success, do nothing
	case promptui.ErrInterrupt:
		fmt.Println("Interrupted!")
		return
	case promptui.ErrEOF:
		fmt.Println("Cancelled!")
		return
	default:
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Check if user select same context
	if result == config.CurrentContext+"(current-context)" {
		fmt.Println("You already selected this context!")
		return
	}

	// if err != nil {
	// 	if err == promptui.ErrInterrupt {
	// 		// User click Ctrl+C
	// 		fmt.Println("Interrupted by Ctrl+C! Exiting...")
	// 	} else if err == promptui.ErrEOF {
	// 		// User click Ctrl+D or ESC
	// 		// Test in Mac, ESC doesn't trigger this. Only Ctrl+D works
	// 		fmt.Println("Cancelled by Ctrl+D or ESC! Exiting...")
	// 	} else {
	// 		// Other errors
	// 		fmt.Printf("Prompt failed: %v\n", err)
	// 	}
	// 	return
	// }

	fmt.Println("Your selection: ", result)

	// Write back to file section //
	// Remove (current-context) from result. StrimSuffix first, update struct later on.
	result = strings.TrimSuffix(result, "(current-context)")
	// Update struct
	config.CurrentContext = result

	// saveConfig
	err = saveConfig("./testdata/kubeconfig", config) // config is pointer already.
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
}
