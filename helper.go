package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
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

/*
// Yellow text helper xD
func yellow(s string) string {
	return "\x1b[33m" + s + "\x1b[0m"
}

// Red text helper xD
func red(s string) string {
	return "\x1b[31m" + s + "\x1b[0m"
}

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
