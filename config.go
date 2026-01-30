package main // Still same package because same folder xD

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// For context switch back (ctx -)
const historyFileName = "chg-k8s-ctx-history"

// Parse kubeconfig file to struct
func loadConfig(path string) (*KubeConfig, error) {
	// Look at the func ReadFile. We don't need to talk about param
	// Only understand this: ([]byte, error)
	// fmt.Println("Try to read fucking file.")

	// Test
	_data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err // Because we return 2 values
	}

	// Struct already created, not we need pass the fucking data to my struct
	var config KubeConfig // init fucking struct

	// err = yaml.Unmarshal(_data, &config) // If we used: :=  , it will not works. Because it works only for the first time
	err_1 := yaml.Unmarshal(_data, &config) // This works because we never used err_1 var before!
	// Look at those code above, we already used err var
	// Unmarshal take 2 params: []byte and *interface{}, return error
	// That is why we put _data which is fucking byte and &config which is fucking pointer to struct

	// Check for fucking error bro!
	if err_1 != nil {
		fmt.Println("Error unmarshalling YAML:", err_1)
		return nil, err_1 // Because we return 2 values
	}

	// List context of K8S
	// Now we have debug mode, we need to put this into debug
	if debugMode {
		fmt.Println("[DEBUG] Current context:", config.CurrentContext)

		// Print current namespace if not empty
		for _, ctx := range config.Contexts {
			if ctx.Name == config.CurrentContext {
				if ctx.Context.Namespace != "" {
					fmt.Println("[DEBUG] Default namespace:", ctx.Context.Namespace)
					break
				}
			}
		}
	}

	// Func need return, so we return what it needs
	return &config, nil

}

// Save kubeconfig file from struct to file
func saveConfig(path string, config *KubeConfig) error {
	// Oke, next step backup current file to .bak extension, overwrite if exist for simplicity
	// This should be fine even we run multiple times.
	// And yeah, only backup after user select the context.

	// in function, we need to read file again to back it up.
	_data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err // Because we return error
	}

	// Backup file path.
	bak_path := path + ".bak"

	// Only backup if .bak doesn't exists, first time only
	if _, err := os.Stat(bak_path); os.IsNotExist(err) {
		// .bak doesn't exist → create backup
		err = os.WriteFile(bak_path, _data, 0644)

		// Check for fucking error bro!
		if err != nil {
			fmt.Println("Warning: Failed to create backup:", err)
			return err // No fucking continue. No backup, no game.
		}
		// No need to print anything if backup is successful
	}

	// Write to file
	newData, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("Error marshalling in section write back to file")
		return err
	}
	err = os.WriteFile(path, newData, 0600)
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}

// Design for write config to file. Used for ctx merge command.
// Since I don't want to edit saveConfig() func xD
func writeConfig(path string, config *KubeConfig) error {
	newData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, newData, 0600)
}

// Get default kubeconfig path
func getKubeconfigPath() (string, error) {
	// 1. Kubeconfig
	// 2. KUBECONFIG env var
	// 3. Default kubeconfig

	if kubeconfigFlag != "" {
		if debugMode {
			fmt.Println("Using --kubeconfig flag: ", kubeconfigFlag)
		}
		// Change to stderr. prevent issue in export command. Hmm in fact it doesn't
		//fmt.Fprintf(os.Stderr, "Using --kubeconfig flag: %s\n", kubeconfigFlag)
		return kubeconfigFlag, nil
	}

	// Check if KUBECONFIG environment variable is set
	// This is a common Go idiom where the variable is declared and checked in the same line.
	// I'm trying to become GO idiom LOL
	if path := os.Getenv("KUBECONFIG"); path != "" {
		fmt.Println("Using KUBECONFIG environment variable: ", path)
		return path, nil
	} else {
		// Get home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		// Default kubeconfig path
		// fmt.Println("Using default kubeconfig path: ", homeDir+"/.kube/config")
		// return homeDir + "/.kube/config", nil

		// Use filepath.Join to create a path that is safe for any operating system.
		defaultPath := filepath.Join(homeDir, ".kube", "config")
		fmt.Println("Using default kubeconfig path:", defaultPath)
		return defaultPath, nil
	}
}

// Get history file Path - For multiple OS support
func getHistoryFilePath() (string, error) {
	kubeconfigPath, err := getKubeconfigPath() // DRY
	if err != nil {
		return "", err
	}

	// History file same folder with kubeconfig
	dir := filepath.Dir(kubeconfigPath)
	return filepath.Join(dir, historyFileName), nil
}
