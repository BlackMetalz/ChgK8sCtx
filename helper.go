package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/sahilm/fuzzy"
)

// getCurrentContextEntry returns pointer to the current context entry
// Returns nil if not found (shouldn't happen in valid kubeconfig, but who knows?)
func getCurrentContextEntry(config *KubeConfig) *Context {
	for i := range config.Contexts {
		if config.Contexts[i].Name == config.CurrentContext {
			return &config.Contexts[i] // return pointer to real struct
		}
	}
	return nil
}

// Handle Promt Error.
func handlePromptError(err error) error {
	// Go idiom style
	switch err {
	case nil:
		return nil
	case promptui.ErrInterrupt:
		return fmt.Errorf("interrupted by Ctrl+C")
	case promptui.ErrEOF:
		return fmt.Errorf("cancelled by Ctrl+D or ESC")
	default:
		return fmt.Errorf("prompt failed: %v", err)
	}
}

// Green text helper xD
func green(s string) string {
	return "\x1b[32m" + s + "\x1b[0m"
}

// Red text helper xD
func red(s string) string {
	return "\x1b[31m" + s + "\x1b[0m"
}

// Yellow text helper xD
func yellow(s string) string {
	return "\x1b[33m" + s + "\x1b[0m"
}

/*




// Blue text helper xD
func blue(s string) string {
	return "\x1b[34m" + s + "\x1b[0m"
}

// Magenta text helper xD
func magenta(s string) string {
	return "\x1b[35m" + s + "\x1b[0m"
}

// Cyan text helper xD
func cyan(s string) string {
	return "\x1b[36m" + s + "\x1b[0m"
}
*/

/*
// Remove context by name in slice of Context struct
// Logic is very simple, everyone include me could understand
// just put name of context you want to remove in function, it will return a Context slice without that context
func removeContextByName(contexts []Context, name string) []Context {
	var result []Context
	for _, ctx := range contexts {
		if ctx.Name != name {
			result = append(result, ctx)
		}
	}
	return result
}

// helper func for remove user from kubeconfig
// Hmm, seem like same context, opportunity for refactoring!
func removeUserByName(users []User, name string) []User {
	var result []User
	for _, user := range users {
		if user.Name != name {
			result = append(result, user)
		}
	}
	return result
}

// Same bro!
func removeClusterByName(clusters []Cluster, name string) []Cluster {
	var result []Cluster
	for _, c := range clusters {
		if c.Name != name {
			result = append(result, c)
		}
	}
	return result
}

*/

// Check if items exists in kubeconfig
// Works with context,cluster,user
func itemExists(config *KubeConfig, itemType, name string) bool {
	switch itemType {
	// Go idiom style
	case "context":
		for _, ctx := range config.Contexts {
			if ctx.Name == name {
				return true
			}
		}
	case "cluster":
		for _, c := range config.Clusters {
			if c.Name == name {
				return true
			}
		}
	case "user":
		for _, u := range config.Users {
			if u.Name == name {
				return true
			}
		}
	}
	return false
}

// Get contexts that use this user/cluster
func getUsedBy(config *KubeConfig, itemType, name string) []string {
	var result []string
	for _, ctx := range config.Contexts {
		switch itemType {
		case "cluster":
			if ctx.Context.Cluster == name {
				result = append(result, ctx.Name)
			}
		case "user":
			if ctx.Context.User == name {
				result = append(result, ctx.Name)
			}
		}
	}
	return result
}

// Confirm delete prompt
func confirmDelete(msg string) bool {
	prompt := promptui.Select{
		Label: msg,
		Items: []string{"Yes", "No"},
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return false
	}

	return selected == "Yes"
}

// For cleanup flag
// Logic: if it doesn't used by any context, it is orphan
func getOrphanUsers(config *KubeConfig) []string {
	var orphans []string
	for _, u := range config.Users {
		usedBy := getUsedBy(config, "user", u.Name)
		if len(usedBy) == 0 {
			orphans = append(orphans, u.Name)
		}
	}
	return orphans
}

func getOrphanClusters(config *KubeConfig) []string {
	var orphans []string
	for _, c := range config.Clusters {
		usedBy := getUsedBy(config, "cluster", c.Name)
		if len(usedBy) == 0 {
			orphans = append(orphans, c.Name)
		}
	}
	return orphans
}

// For Fuzzy Search
func fuzzyFindContext(config *KubeConfig, query string) []string {
	// get all context names
	var ctxNames []string
	for _, ctx := range config.Contexts {
		ctxNames = append(ctxNames, ctx.Name)
	}

	// Fuzzy search
	matches := fuzzy.Find(query, ctxNames)

	// Extract matched names
	var results []string
	for _, match := range matches {
		results = append(results, match.Str)
	}

	return results
}

// For history switch back
func savePreviousContext(kubeconfigPath string, ctx string) error {
	// Write to history file
	filePath := getHistoryFilePath(kubeconfigPath)

	// Open file with truncate mode to overwrite
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write context name to file
	// Alway replace/overwrite file, because support last history, not history list
	_, err = file.WriteString(ctx + "\n")
	if err != nil {
		return err
	}

	return nil
}

func loadPreviousContext(kubeconfigPath string) (string, error) {
	// read from history file
	filePath := getHistoryFilePath(kubeconfigPath)

	// Open file in read mode
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file content. Since this will have only 1 line as we expected from savePreviousContext()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil

}
