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
	oldContext := config.CurrentContext
	config.CurrentContext = result

	// Save previous context
	err = savePreviousContext(oldContext)
	if err != nil {
		fmt.Println("Error saving previous context")
		return err
	}

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
		Label: "Select default namespace for current context " + green(config.CurrentContext),
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

// Handler for rename context
func renameContext(config *KubeConfig, kubeconfigPath string, directArgs ...string) error {

	// directArgs is slice of string and it takes more than 1 args, but we only handle 2 for now. haha
	if len(directArgs) > 2 {
		fmt.Println("Unsupported arguments, only 2 arguments are supported")
		return nil
	}

	// So we gonna refactor this function to handle both interactive and direct mode
	// Validation part for both modes need to be shared for "DRY"

	// define
	var oldName, newName string
	var err error

	if len(directArgs) == 2 {
		// Direct mode here
		oldName = directArgs[0]
		newName = directArgs[1]
	} else {
		// Interactive mode here
		// Build list context for selection, switchContext like.

		var contextList []string
		for _, ctx := range config.Contexts {
			if ctx.Name == config.CurrentContext {
				contextList = append(contextList, ctx.Name+"(current-context)")
			} else {
				contextList = append(contextList, ctx.Name)
			}
		}

		// Prompt for context need to rename
		selectPrompt := promptui.Select{
			Label: "Select context to rename please",
			Items: contextList,
		}

		_, selectedContext, err := selectPrompt.Run()
		if err != nil {
			return handlePromptError(err)
		}

		// Clean suffix if available
		oldName = strings.TrimSuffix(selectedContext, "(current-context)")

		// Prompt for fucking new name
		inputPrompt := promptui.Prompt{
			Label: "Enter for new context name",
			// Default: oldName, // Pre-fill with old name
			// We should not pick default value for input xD. It is just a distraction.
		}

		// Need to update for err, because we are in scope of else block.
		// So we can not define like: newName, err := inputPrompt.Run()
		newName, err = inputPrompt.Run()
		if err != nil {
			return handlePromptError(err)
		}
	}

	// Simple validation for direct mode
	if len(directArgs) == 2 {
		found := false // Flag here
		// Loop throught contexts...
		for _, ctx := range config.Contexts {
			if ctx.Name == oldName {
				found = true
				break // Found, no need to continue iteration!
			}
		}

		// Yeah, if context not found, we can not rename it, right?
		// What a simple validation bro!
		if !found {
			fmt.Printf("Context %s not found bro!\n", red(oldName))
			return nil
		}
	}

	// Validation newName to prevent duplicate with other context
	// Still need to loop through all contexts
	for _, ctx := range config.Contexts {
		if ctx.Name == newName && ctx.Name != oldName {
			fmt.Printf("Context %s already exists bro!\n", red(newName))
			return nil
		}
		// Check if user enter same context name
		if ctx.Name == newName {
			fmt.Println("New name cannot be same as old name bro!")
			return nil
		}
	}

	if strings.TrimSpace(newName) == "" {
		fmt.Println("New name cannot be empty bro!")
		return nil
	}

	// Update config
	// Loop through all context
	for i := range config.Contexts {
		// If found old name
		if config.Contexts[i].Name == oldName {
			// Update it with new name
			config.Contexts[i].Name = newName

			// We would add break here for performance
			// I mean we don't need to continue loop after we found it, right?
			break
		}
	}

	// If we rename current-context, update it aswell
	if config.CurrentContext == oldName {
		config.CurrentContext = newName
	}

	// Save it
	err = saveConfig(kubeconfigPath, config)
	if err != nil {
		return err
	}

	fmt.Printf("Renamed context: %s to %s\n", oldName, newName)
	return nil
}

// List context for flag -c and -l
func listContexts(config *KubeConfig) {
	for _, ctx := range config.Contexts {
		// Show current context with different color? xD
		if ctx.Name == config.CurrentContext {
			fmt.Println(green(ctx.Name))
			// With printf, no space LOL
		} else {
			fmt.Println(ctx.Name)
		}
	}
}

// Show some useful information about current context xD
func showCurrentContext(config *KubeConfig) {
	// Print current context name
	fmt.Println(strings.Repeat("=", 10) + "Current context" + strings.Repeat("=", 10))
	fmt.Printf("Context: %s\n", green(config.CurrentContext))

	// Print user and cluster for current context
	for _, ctx := range config.Contexts {
		if ctx.Name == config.CurrentContext {
			fmt.Printf("Cluster: %s\n", green(ctx.Context.Cluster))
			fmt.Printf("User: %s\n", green(ctx.Context.User))
			if ctx.Context.Namespace != "" {
				fmt.Printf("Namespace: %s\n", green(ctx.Context.Namespace))
			} else {
				fmt.Printf("Namespace: %s\n", green("(default)"))
			}
			break // quit iteration after found xD
		}
	}
	fmt.Println(strings.Repeat("=", 35))
}

// func for delete context
func deleteContext(config *KubeConfig, kubeconfigPath string, directArgs ...string) error {
	// Support delete only one context at a time

	// Check if only 1 context left --> not allow to delete
	if len(config.Contexts) == 1 {
		fmt.Println(red("Cannot delete the only 1 context left!"))
		return nil
	}

	// Get context to delete (interactive or direct mode)
	var targetContext string
	if len(directArgs) == 1 {
		targetContext = directArgs[0]

		// Validate again, like rename context. Haizzz
		// This could be refactored to a function later.
		found := false // Fucking flag var.
		for _, ctx := range config.Contexts {
			if ctx.Name == targetContext {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Context %s not found bro!\n", red(targetContext))
			return nil
		}
	} else {
		// Interactive selection goes here.
		// Reuse code from switchContext/renameContext. Hmm another refactoring opportunity =.=!
		// fmt.Println("Interactive selection goes here.")
		var contextList []string
		for _, ctx := range config.Contexts {
			if ctx.Name == config.CurrentContext {
				contextList = append(contextList, ctx.Name+"(current-context)")
			} else {
				contextList = append(contextList, ctx.Name)
			}
		}

		prompt := promptui.Select{
			Label: "Select context to delete please",
			Items: contextList,
		}

		_, selectedContext, err := prompt.Run()
		if err != nil {
			return handlePromptError(err)
		}

		targetContext = strings.TrimSuffix(selectedContext, "(current-context)")
	}

	// Remove from context slices
	// Modify in memory only, persist in saveConfig() func
	// Use generic func to remove context by name
	// config.Contexts = removeContextByName(config.Contexts, targetContext) // Before
	config.Contexts = deleteEntryByName(config.Contexts, targetContext, func(ctx Context) string {
		return ctx.Name
	})

	// If deleted current-context --> auto switch to another context
	if targetContext == config.CurrentContext {
		config.CurrentContext = config.Contexts[0].Name // Get the first of fucking context list
		fmt.Printf("Auto switched to context: %s\n", green(config.CurrentContext))
	}

	// Save
	err := saveConfig(kubeconfigPath, config)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted context: %s\n", red(targetContext))
	return nil
}

// For Delete User in Kubeconfig file.
func deleteUser(config *KubeConfig, path string, directArgs ...string) error {
	var targetUser string
	// get target user (interactive or direct mode)
	if len(directArgs) == 1 {
		targetUser = directArgs[0]
		// Validate
		if !itemExists(config, "user", targetUser) {
			fmt.Printf("User %s not found bro!\n", red(targetUser))
			return nil
		}

	} else {
		// Interactive - prompt for select user to delete
		var userList []string
		for _, u := range config.Users {
			userList = append(userList, u.Name)
		}
		prompt := promptui.Select{Label: "Select user", Items: userList}
		var err error // Need to define err var before using it
		_, targetUser, err = prompt.Run()
		if err != nil {
			// Use Helper func to handle Prompt Errors
			return handlePromptError(err)
		}
	}

	// check usedBy. Shared logic
	usedBy := getUsedBy(config, "user", targetUser)
	if len(usedBy) > 0 {
		fmt.Printf("Warning: Will ALSO delete context(s): %s\n", yellow(strings.Join(usedBy, ", ")))
		if !confirmDelete("Delete this user AND related contexts?") {
			fmt.Println("Cancelled.")
			return nil
		}

		// Delete related context FIRST
		for _, ctxName := range usedBy {
			// config.Contexts = removeContextByName(config.Contexts, ctxName)
			config.Contexts = deleteEntryByName(config.Contexts, ctxName, func(c Context) string {
				return c.Name
			})
			fmt.Printf("Deleted context: %s\n", red(ctxName))

			// Auto switch if deleted current context
			if ctxName == config.CurrentContext && len(config.Contexts) > 0 {
				config.CurrentContext = config.Contexts[0].Name
				fmt.Printf("Auto switched to context: %s\n", green(config.CurrentContext))
			}
		}
	}

	// Delete user
	// config.Users = removeUserByName(config.Users, targetUser)
	// Refactor this shit.
	config.Users = deleteEntryByName(config.Users, targetUser, func(u User) string {
		return u.Name
	})

	if err := saveConfig(path, config); err != nil {
		return err
	}

	fmt.Printf("Deleted user: %s\n", red(targetUser))
	return nil
}

// for delete cluster in kubeconfig file
func deleteCluster(config *KubeConfig, path string, directArgs ...string) error {
	var targetCluster string
	// get target cluster (interactive or direct mode)
	if len(directArgs) == 1 {
		targetCluster = directArgs[0]
		// Validate
		if !itemExists(config, "cluster", targetCluster) {
			fmt.Printf("Cluster %s not found bro!\n", red(targetCluster))
			return nil
		}
	} else {
		// Interactive - prompt for select user to delete
		var clusterList []string
		for _, c := range config.Clusters {
			clusterList = append(clusterList, c.Name)
		}
		prompt := promptui.Select{Label: "Select cluster", Items: clusterList}
		var err error // Need to define err var before using it
		_, targetCluster, err = prompt.Run()
		if err != nil {
			// Use Helper func to handle Prompt Errors
			return handlePromptError(err)
		}

	}

	// check usedBy. Shared logic
	usedBy := getUsedBy(config, "cluster", targetCluster)
	if len(usedBy) > 0 {
		fmt.Printf("Warning: Will ALSO delete context(s): %s\n", yellow(strings.Join(usedBy, ", ")))
		if !confirmDelete("Delete cluster AND related contexts?") {
			fmt.Println("Cancelled.")
			return nil
		}

		// Delete related context FIRST
		for _, ctxName := range usedBy {
			// config.Contexts = removeContextByName(config.Contexts, ctxName)
			config.Contexts = deleteEntryByName(config.Contexts, ctxName, func(c Context) string {
				return c.Name
			})
			fmt.Printf("Deleted context: %s\n", red(ctxName))

			// Auto switch if deleted current context
			if ctxName == config.CurrentContext && len(config.Contexts) > 0 {
				config.CurrentContext = config.Contexts[0].Name
				fmt.Printf("Auto switched to context: %s\n", green(config.CurrentContext))
			}
		}
	}

	// Delete cluster
	// config.Clusters = removeClusterByName(config.Clusters, targetCluster)
	// Refactor this shiet
	config.Clusters = deleteEntryByName(config.Clusters, targetCluster, func(c Cluster) string {
		return c.Name
	})

	if err := saveConfig(path, config); err != nil {
		return err
	}

	fmt.Printf("Deleted cluster: %s\n", red(targetCluster))
	return nil
}

func deleteOrphanData(config *KubeConfig, path string) error {
	// get all orphan data first
	orphanUsers := getOrphanUsers(config)
	orphanClusters := getOrphanClusters(config)

	// Check if any orphans exists, if not exit
	if len(orphanUsers) == 0 && len(orphanClusters) == 0 {
		fmt.Println(green("No orphan data found bro!"))
		return nil
	}

	// Show all orphans
	if len(orphanUsers) > 0 {
		fmt.Printf("Found %d orphan users: %s\n", len(orphanUsers), yellow(strings.Join(orphanUsers, ", ")))
	}

	if len(orphanClusters) > 0 {
		fmt.Printf("Found %d orphan clusters: %s\n", len(orphanClusters), yellow(strings.Join(orphanClusters, ", ")))
	}

	// Single confirm for both cluster and user
	if !confirmDelete("Delete ALL orphan items?") {
		fmt.Println("Cancelled.")
		return nil
	}

	// Delete all orphan users
	for _, u := range orphanUsers {
		// config.Users = removeUserByName(config.Users, u)
		config.Users = deleteEntryByName(config.Users, u, func(u User) string {
			return u.Name
		})
		fmt.Printf("Deleted user: %s\n", red(u))
	}

	// Delete all orphan clusters
	for _, c := range orphanClusters {
		// config.Clusters = removeClusterByName(config.Clusters, c)
		config.Clusters = deleteEntryByName(config.Clusters, c, func(c Cluster) string {
			return c.Name
		})
		fmt.Printf("Deleted cluster: %s\n", red(c))
	}

	// Save once at the end, after deletion!
	if err := saveConfig(path, config); err != nil {
		return err
	}

	fmt.Println(green("Deleted all orphan items successfully!"))
	return nil

}

func exportContext(config *KubeConfig, contextName string) *KubeConfig {
	// Find context
	var targetCtx Context
	for _, ctx := range config.Contexts {
		if ctx.Name == contextName {
			targetCtx = ctx
			break
		}
	}

	if targetCtx.Name == "" {
		fmt.Println(red("Context not found!"))
		return nil
	}

	// Find cluster
	var targetCluster Cluster
	for _, c := range config.Clusters {
		if c.Name == targetCtx.Context.Cluster {
			targetCluster = c
			break
		}
	}

	// Find user
	var targetUser User
	for _, u := range config.Users {
		if u.Name == targetCtx.Context.User {
			targetUser = u
			break
		}
	}

	// Create new kubeconfig only contain target context
	exportCtx := &KubeConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: contextName,
		Contexts:       []Context{targetCtx},
		Users:          []User{targetUser},
		Clusters:       []Cluster{targetCluster},
	}

	return exportCtx
}

func importContext(importPath string) {
	// Load current kubeconfig to import
	currentPath, err := getKubeconfigPath()
	if err != nil {
		fmt.Println(red("Error getting current kubeconfig path:"), err)
		return
	}
	currentConfig, err := loadConfig(currentPath)
	if err != nil {
		fmt.Println(red("Error loading current kubeconfig:"), err)
		return
	}

	// Load import file
	importConfig, err := loadConfig(importPath)
	if err != nil {
		fmt.Println(red("Error loading import kubeconfig:"), err)
		return
	}

	// Merge - reuse existing merge func xD
	mergeConfigs(importConfig, currentConfig)

	// Save
	saveConfig(currentPath, currentConfig)

	fmt.Printf("Imported context successfully to %s!\n", green(currentPath))
}

// New refactor func for delete Context/User/Cluster
// This would use generics to handle delete entries
// Doc: https://go.dev/doc/tutorial/generics
func deleteEntryByName[Entry any](items []Entry, name string, getName func(Entry) string) []Entry {
	var result []Entry
	// Start loop through
	// Same logic like removeByXXX. Check if name is not equal to name
	// Then append to result, so the name matched would be removed
	for _, item := range items {
		if getName(item) != name {
			result = append(result, item)
		}
	}

	return result
}

// This func is little hard for me. I don't care about writing comment in code. Haha
/* Example usage
config.Contexts = deleteEntryByName(
    config.Contexts,         // 1. Slice of Context struct
    targetContext,           // 2. Name (string) to remove
    func(ctx Context) string { return ctx.Name }  // 3. "How to get name from this struct"
)
*/

// Test
/*
func removeByName[T any](items []T, name string, getName func(T) string) []T {
	var result []T
	for _, item := range items {
		if getName(item) != name {
			result = append(result, item)
		}
	}
	return result
}
*/
