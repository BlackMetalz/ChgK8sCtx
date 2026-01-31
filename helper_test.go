package main

import (
	"testing"
)

// Helper function to create test config
func createTestConfig() *KubeConfig {
	return &KubeConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "dev-cluster",
		Contexts: []Context{
			{Name: "dev-cluster", Context: struct {
				Cluster   string `yaml:"cluster"`
				User      string `yaml:"user"`
				Namespace string `yaml:"namespace"`
			}{Cluster: "dev-cluster", User: "dev-admin"}},
			{Name: "prod-cluster", Context: struct {
				Cluster   string `yaml:"cluster"`
				User      string `yaml:"user"`
				Namespace string `yaml:"namespace"`
			}{Cluster: "prod-cluster", User: "prod-admin"}},
			{Name: "staging-cluster", Context: struct {
				Cluster   string `yaml:"cluster"`
				User      string `yaml:"user"`
				Namespace string `yaml:"namespace"`
			}{Cluster: "staging-cluster", User: "staging-admin"}},
		},
		Users: []User{
			{Name: "dev-admin"},
			{Name: "prod-admin"},
			{Name: "staging-admin"},
			{Name: "orphan-user"}, // Not used by any context
		},
		Clusters: []Cluster{
			{Name: "dev-cluster"},
			{Name: "prod-cluster"},
			{Name: "staging-cluster"},
			{Name: "orphan-cluster"}, // Not used by any context
		},
	}
}

// ========================================
// Test itemExists
// ========================================
func TestItemExists(t *testing.T) {
	config := createTestConfig()

	tests := []struct {
		name     string
		itemType string
		itemName string
		want     bool
	}{
		// Context tests
		{"context exists - dev", "context", "dev-cluster", true},
		{"context exists - prod", "context", "prod-cluster", true},
		{"context not exists", "context", "non-existent", false},

		// User tests
		{"user exists - dev-admin", "user", "dev-admin", true},
		{"user exists - orphan", "user", "orphan-user", true},
		{"user not exists", "user", "ghost-user", false},

		// Cluster tests
		{"cluster exists - dev", "cluster", "dev-cluster", true},
		{"cluster exists - orphan", "cluster", "orphan-cluster", true},
		{"cluster not exists", "cluster", "ghost-cluster", false},

		// Invalid item type
		{"invalid item type", "invalid", "dev-cluster", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := itemExists(config, tt.itemType, tt.itemName)
			if got != tt.want {
				t.Errorf("itemExists(%s, %s) = %v; want %v",
					tt.itemType, tt.itemName, got, tt.want)
			}
		})
	}
}

// ========================================
// Test getUsedBy
// ========================================
func TestGetUsedBy(t *testing.T) {
	config := createTestConfig()

	tests := []struct {
		name     string
		itemType string
		itemName string
		want     []string
	}{
		{"user used by context", "user", "dev-admin", []string{"dev-cluster"}},
		{"cluster used by context", "cluster", "prod-cluster", []string{"prod-cluster"}},
		{"orphan user - no usage", "user", "orphan-user", nil},
		{"orphan cluster - no usage", "cluster", "orphan-cluster", nil},
		{"non-existent user", "user", "ghost-user", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getUsedBy(config, tt.itemType, tt.itemName)
			if len(got) != len(tt.want) {
				t.Errorf("getUsedBy(%s, %s) = %v; want %v",
					tt.itemType, tt.itemName, got, tt.want)
			}
		})
	}
}

// ========================================
// Test getOrphanUsers
// ========================================
func TestGetOrphanUsers(t *testing.T) {
	config := createTestConfig()

	orphans := getOrphanUsers(config)

	if len(orphans) != 1 {
		t.Errorf("getOrphanUsers() = %v; want 1 orphan", orphans)
	}

	if len(orphans) > 0 && orphans[0] != "orphan-user" {
		t.Errorf("getOrphanUsers() = %v; want [orphan-user]", orphans)
	}
}

// ========================================
// Test getOrphanClusters
// ========================================
func TestGetOrphanClusters(t *testing.T) {
	config := createTestConfig()

	orphans := getOrphanClusters(config)

	if len(orphans) != 1 {
		t.Errorf("getOrphanClusters() = %v; want 1 orphan", orphans)
	}

	if len(orphans) > 0 && orphans[0] != "orphan-cluster" {
		t.Errorf("getOrphanClusters() = %v; want [orphan-cluster]", orphans)
	}
}

// ========================================
// Test deleteEntryByName (Generic function!)
// ========================================
func TestDeleteEntryByName(t *testing.T) {
	t.Run("delete context by name", func(t *testing.T) {
		config := createTestConfig()
		original := len(config.Contexts) // 3

		result := deleteEntryByName(config.Contexts, "dev-cluster", func(c Context) string {
			return c.Name
		})

		if len(result) != original-1 {
			t.Errorf("deleteEntryByName() contexts len = %d; want %d", len(result), original-1)
		}

		// Check that dev-cluster is actually removed
		for _, ctx := range result {
			if ctx.Name == "dev-cluster" {
				t.Errorf("deleteEntryByName() should have removed dev-cluster")
			}
		}
	})

	t.Run("delete user by name", func(t *testing.T) {
		config := createTestConfig()
		original := len(config.Users) // 4

		result := deleteEntryByName(config.Users, "orphan-user", func(u User) string {
			return u.Name
		})

		if len(result) != original-1 {
			t.Errorf("deleteEntryByName() users len = %d; want %d", len(result), original-1)
		}
	})

	t.Run("delete cluster by name", func(t *testing.T) {
		config := createTestConfig()
		original := len(config.Clusters) // 4

		result := deleteEntryByName(config.Clusters, "orphan-cluster", func(c Cluster) string {
			return c.Name
		})

		if len(result) != original-1 {
			t.Errorf("deleteEntryByName() clusters len = %d; want %d", len(result), original-1)
		}
	})

	t.Run("delete non-existent item", func(t *testing.T) {
		config := createTestConfig()
		original := len(config.Contexts)

		result := deleteEntryByName(config.Contexts, "non-existent", func(c Context) string {
			return c.Name
		})

		// Should return same length - nothing deleted
		if len(result) != original {
			t.Errorf("deleteEntryByName() should not delete non-existent item")
		}
	})
}

// ========================================
// Test getCurrentContextEntry
// ========================================
func TestGetCurrentContextEntry(t *testing.T) {
	config := createTestConfig()

	// Test found case
	entry := getCurrentContextEntry(config)
	if entry == nil {
		t.Errorf("getCurrentContextEntry() = nil; want dev-cluster")
	}
	if entry != nil && entry.Name != "dev-cluster" {
		t.Errorf("getCurrentContextEntry().Name = %s; want dev-cluster", entry.Name)
	}

	// Test not found case
	config.CurrentContext = "non-existent"
	entry = getCurrentContextEntry(config)
	if entry != nil {
		t.Errorf("getCurrentContextEntry() = %v; want nil", entry)
	}
}

// ========================================
// Test fuzzyFindContext
// ========================================
func TestFuzzyFindContext(t *testing.T) {
	config := createTestConfig()

	tests := []struct {
		name     string
		query    string
		wantLen  int
		wantName string // First match expected
	}{
		{"exact match", "dev-cluster", 1, "dev-cluster"},
		{"partial match - dev", "dev", 1, "dev-cluster"},
		{"partial match - prod", "prod", 1, "prod-cluster"},
		{"partial match - cluster", "cluster", 3, ""}, // All contain "cluster"
		{"no match", "xyz", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyFindContext(config, tt.query)
			if len(got) != tt.wantLen {
				t.Errorf("fuzzyFindContext(%s) len = %d; want %d",
					tt.query, len(got), tt.wantLen)
			}
			if tt.wantName != "" && len(got) > 0 && got[0] != tt.wantName {
				t.Errorf("fuzzyFindContext(%s)[0] = %s; want %s",
					tt.query, got[0], tt.wantName)
			}
		})
	}
}

// ========================================
// Test color helpers (simple sanity check)
// ========================================
func TestColorHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) string
		input    string
		contains string
	}{
		{"green", green, "test", "test"},
		{"red", red, "error", "error"},
		{"yellow", yellow, "warning", "warning"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if got == "" {
				t.Errorf("%s(%s) = empty string", tt.name, tt.input)
			}
			// Just check it contains the original text
			if len(got) <= len(tt.input) {
				t.Errorf("%s(%s) should add ANSI codes", tt.name, tt.input)
			}
		})
	}
}
