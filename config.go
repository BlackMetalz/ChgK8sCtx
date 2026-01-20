package main // Still same package because same folder xD

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

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
	fmt.Println("Current context:", config.CurrentContext)

	// Func need return, so we return what it needs
	return &config, nil

}

func saveConfig(path string, config *KubeConfig) error {
	// Oke, next step backup current file to .bak extension, overwrite if exist for simplicity
	// This should be fine even we run multiple times.
	// And yeah, only backup after user select the context.

	// in function, we need to read file again to back it up.
	_data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err // Because we return 2 values
	}

	// Backup file
	bak_path := path + ".bak"
	err = os.WriteFile(bak_path, _data, 0644)
	if err != nil {
		fmt.Println("Warning: Failed to create backup:", err)
		return err // No fucking continue. No backup, no game.
	} else {
		fmt.Println("Current kubeconfig file was backed up to ", bak_path)
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

	fmt.Println("Successfully switched to context: ", config.CurrentContext)

	return nil
}
