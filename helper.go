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
