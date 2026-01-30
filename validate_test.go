package main

import (
	"testing"
)

// Reuse helper from helper_test.go - createTestConfig()
// Note: Go test files in same package share scope

// ========================================
// Test validateConfig
// ========================================
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name       string
		config     *KubeConfig
		wantErrors int
	}{
		{
			name:       "valid config - no errors",
			config:     createTestConfig(),
			wantErrors: 0,
		},
		{
			name: "broken user reference",
			config: &KubeConfig{
				CurrentContext: "test-ctx",
				Contexts: []Context{
					{Name: "test-ctx", Context: struct {
						Cluster   string `yaml:"cluster"`
						User      string `yaml:"user"`
						Namespace string `yaml:"namespace"`
					}{Cluster: "test-cluster", User: "ghost-user"}},
				},
				Users:    []User{{Name: "real-user"}},
				Clusters: []Cluster{{Name: "test-cluster"}},
			},
			wantErrors: 1, // ghost-user not found
		},
		{
			name: "broken cluster reference",
			config: &KubeConfig{
				CurrentContext: "test-ctx",
				Contexts: []Context{
					{Name: "test-ctx", Context: struct {
						Cluster   string `yaml:"cluster"`
						User      string `yaml:"user"`
						Namespace string `yaml:"namespace"`
					}{Cluster: "ghost-cluster", User: "test-user"}},
				},
				Users:    []User{{Name: "test-user"}},
				Clusters: []Cluster{{Name: "real-cluster"}},
			},
			wantErrors: 1, // ghost-cluster not found
		},
		{
			name: "broken current-context",
			config: &KubeConfig{
				CurrentContext: "non-existent-context",
				Contexts: []Context{
					{Name: "real-ctx", Context: struct {
						Cluster   string `yaml:"cluster"`
						User      string `yaml:"user"`
						Namespace string `yaml:"namespace"`
					}{Cluster: "test-cluster", User: "test-user"}},
				},
				Users:    []User{{Name: "test-user"}},
				Clusters: []Cluster{{Name: "test-cluster"}},
			},
			wantErrors: 1, // current-context not found
		},
		{
			name: "multiple errors",
			config: &KubeConfig{
				CurrentContext: "broken-ctx",
				Contexts: []Context{
					{Name: "test-ctx", Context: struct {
						Cluster   string `yaml:"cluster"`
						User      string `yaml:"user"`
						Namespace string `yaml:"namespace"`
					}{Cluster: "ghost-cluster", User: "ghost-user"}},
				},
				Users:    []User{},
				Clusters: []Cluster{},
			},
			wantErrors: 3, // ghost-user, ghost-cluster, broken current-context
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateConfig(tt.config)
			if len(errors) != tt.wantErrors {
				t.Errorf("validateConfig() returned %d errors; want %d. Errors: %v",
					len(errors), tt.wantErrors, errors)
			}
		})
	}
}
