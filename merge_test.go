package main

import (
	"testing"
)

// ========================================
// Test mergeConfigs
// ========================================
func TestMergeConfigs(t *testing.T) {
	t.Run("merge non-overlapping configs", func(t *testing.T) {
		source := &KubeConfig{
			Contexts: []Context{
				{Name: "source-ctx", Context: struct {
					Cluster   string `yaml:"cluster"`
					User      string `yaml:"user"`
					Namespace string `yaml:"namespace"`
				}{Cluster: "source-cluster", User: "source-user"}},
			},
			Users:    []User{{Name: "source-user"}},
			Clusters: []Cluster{{Name: "source-cluster"}},
		}

		target := &KubeConfig{
			Contexts: []Context{
				{Name: "target-ctx", Context: struct {
					Cluster   string `yaml:"cluster"`
					User      string `yaml:"user"`
					Namespace string `yaml:"namespace"`
				}{Cluster: "target-cluster", User: "target-user"}},
			},
			Users:    []User{{Name: "target-user"}},
			Clusters: []Cluster{{Name: "target-cluster"}},
		}

		mergeConfigs(source, target)

		// Target should now have 2 contexts, 2 users, 2 clusters
		if len(target.Contexts) != 2 {
			t.Errorf("mergeConfigs() contexts = %d; want 2", len(target.Contexts))
		}
		if len(target.Users) != 2 {
			t.Errorf("mergeConfigs() users = %d; want 2", len(target.Users))
		}
		if len(target.Clusters) != 2 {
			t.Errorf("mergeConfigs() clusters = %d; want 2", len(target.Clusters))
		}
	})

	t.Run("merge with duplicates - should skip", func(t *testing.T) {
		source := &KubeConfig{
			Contexts: []Context{
				{Name: "same-ctx", Context: struct {
					Cluster   string `yaml:"cluster"`
					User      string `yaml:"user"`
					Namespace string `yaml:"namespace"`
				}{Cluster: "same-cluster", User: "same-user"}},
			},
			Users:    []User{{Name: "same-user"}},
			Clusters: []Cluster{{Name: "same-cluster"}},
		}

		target := &KubeConfig{
			Contexts: []Context{
				{Name: "same-ctx", Context: struct {
					Cluster   string `yaml:"cluster"`
					User      string `yaml:"user"`
					Namespace string `yaml:"namespace"`
				}{Cluster: "same-cluster", User: "same-user"}},
			},
			Users:    []User{{Name: "same-user"}},
			Clusters: []Cluster{{Name: "same-cluster"}},
		}

		mergeConfigs(source, target)

		// Should still have only 1 of each (duplicates skipped)
		if len(target.Contexts) != 1 {
			t.Errorf("mergeConfigs() with duplicates contexts = %d; want 1", len(target.Contexts))
		}
		if len(target.Users) != 1 {
			t.Errorf("mergeConfigs() with duplicates users = %d; want 1", len(target.Users))
		}
		if len(target.Clusters) != 1 {
			t.Errorf("mergeConfigs() with duplicates clusters = %d; want 1", len(target.Clusters))
		}
	})

	t.Run("merge empty source into target", func(t *testing.T) {
		source := &KubeConfig{}
		target := createTestConfig()
		originalLen := len(target.Contexts)

		mergeConfigs(source, target)

		if len(target.Contexts) != originalLen {
			t.Errorf("mergeConfigs() empty source changed target")
		}
	})
}
