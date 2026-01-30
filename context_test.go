package main

import (
	"testing"
)

// ========================================
// Test exportContext
// ========================================
func TestExportContext(t *testing.T) {
	config := createTestConfig()

	t.Run("export existing context", func(t *testing.T) {
		result := exportContext(config, "dev-cluster")

		if result == nil {
			t.Fatal("exportContext() returned nil for existing context")
		}

		// Check exported config has correct structure
		if result.CurrentContext != "dev-cluster" {
			t.Errorf("exportContext().CurrentContext = %s; want dev-cluster", result.CurrentContext)
		}

		if len(result.Contexts) != 1 {
			t.Errorf("exportContext() contexts = %d; want 1", len(result.Contexts))
		}

		if len(result.Users) != 1 {
			t.Errorf("exportContext() users = %d; want 1", len(result.Users))
		}

		if len(result.Clusters) != 1 {
			t.Errorf("exportContext() clusters = %d; want 1", len(result.Clusters))
		}

		// Check correct items were exported
		if result.Contexts[0].Name != "dev-cluster" {
			t.Errorf("exportContext() context name = %s; want dev-cluster", result.Contexts[0].Name)
		}

		if result.Users[0].Name != "dev-admin" {
			t.Errorf("exportContext() user name = %s; want dev-admin", result.Users[0].Name)
		}

		if result.Clusters[0].Name != "dev-cluster" {
			t.Errorf("exportContext() cluster name = %s; want dev-cluster", result.Clusters[0].Name)
		}
	})

	t.Run("export non-existent context", func(t *testing.T) {
		result := exportContext(config, "non-existent-context")

		if result != nil {
			t.Errorf("exportContext() returned non-nil for non-existent context")
		}
	})

	t.Run("export each context", func(t *testing.T) {
		contexts := []string{"dev-cluster", "prod-cluster", "staging-cluster"}

		for _, ctxName := range contexts {
			result := exportContext(config, ctxName)
			if result == nil {
				t.Errorf("exportContext(%s) returned nil", ctxName)
				continue
			}
			if result.CurrentContext != ctxName {
				t.Errorf("exportContext(%s).CurrentContext = %s", ctxName, result.CurrentContext)
			}
		}
	})
}
