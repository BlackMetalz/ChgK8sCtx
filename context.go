package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func switchContext(config *KubeConfig, kubeconfigPath string) error {
	// Move code from main.go to here
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

	if err != nil {
		// Use Helper func to handle Prompt Errors
		return handlePromptError(err)
	}

	// Check if user select same context
	if result == config.CurrentContext+"(current-context)" {
		fmt.Println("You already selected this context!")
		return nil
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
	err = saveConfig(kubeconfigPath, config) // config is pointer already.
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}

func switchNamespace(config *KubeConfig, kubeconfigPath string) error {

	// List all namespace for selection, not manually
	nsList, err := listNamespaces()
	if err != nil {
		fmt.Println("Error getting namespace list:", err)
		return err
	}

	var currentNS string
	// Lets replace this with helper func
	// for i, ctx := range config.Contexts {
	// 	if ctx.Name == config.CurrentContext { // Compare NAME
	// 		currentNS = config.Contexts[i].Context.Namespace
	// 		// config.Contexts[i].Context.Namespace = newNS // Update STRUCT
	// 		break
	// 	}
	// }

	// Use helper func to get current context entry
	entry := getCurrentContextEntry(config)
	if entry == nil {
		fmt.Println("Error getting current context entry")
		return nil
	}
	currentNS = entry.Context.Namespace

	// we had nsList from listNamespaces() func
	var displayList []string // Make fucking slice for display and highlight current NS
	for _, ns := range nsList {
		// "" treat as default namespace, bro!
		if ns == currentNS || (currentNS == "" && ns == "default") {
			displayList = append(displayList, ns+"(current-namespace)")
		} else {
			displayList = append(displayList, ns)
		}
	}

	// Ask user to enter namespace
	prompt := promptui.Select{
		Label: "Select namespace",
		Items: displayList,
	}

	_, newNS, err := prompt.Run()
	if err != nil {
		// Use Helper func to handle Prompt Errors
		return handlePromptError(err)
	}

	// Cleanup before save
	newNS = strings.TrimSuffix(newNS, "(current-namespace)")

	// Check if user select same namespace
	if newNS == currentNS {
		fmt.Println("You already selected this namespace!")
		return nil
	}

	// Update struct
	entry.Context.Namespace = newNS

	// saveConfig
	err = saveConfig(kubeconfigPath, config) // config is pointer already.
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}
