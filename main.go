package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
	// Import this shit for unmarshal yaml
)

// init fucking struct for yaml?
type KubeConfig struct {
	CurrentContext string    `yaml:"current-context"`
	Clusters       []Cluster `yaml:"clusters"` // Link to cluster struct
	Users          []User    `yaml:"users"`    // Link to user struct
	Contexts       []Context `yaml:"contexts"` // Link to context struct
}

type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server       string `yaml:"server"`
		CertAuthData string `yaml:"certificate-authority-data"`
	}
}

type User struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertData string `yaml:"client-certificate-data"`
		ClientKeyData  string `yaml:"client-key-data"`
	}
}

type Context struct {
	Name string `yaml:"name"` // this is what we use to switch context
	// Like Kubectx style.
	Context struct {
		Cluster   string `yaml:"cluster"`
		User      string `yaml:"user"`
		Namespace string `yaml:"namespace"`
	}
}

func main() {
	fmt.Println("Hello, World!")

	// Look at the func ReadFile. We don't need to talk about param
	// Only understand this: ([]byte, error)
	fmt.Println("Try to read fucking file.")

	// Test
	_data, err := os.ReadFile("./testdata/kubeconfig")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Get fucking home dir

	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Println("Error getting home directory:", err)
	// 	return
	// }

	// kubeConfig := homeDir + "/.kube/config"

	// fmt.Println("Kubeconfig file path: ", kubeConfig)
	// _data, err := os.ReadFile(kubeConfig)

	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// Need to convert to string to print, if not it just show the fucking bytes only.
	// fmt.Println(string(_data))

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
		return
	}

	// Check type of config
	fmt.Printf("Type of config: %T\n", config)
	// Output: Type of config: main.KubeConfig
	// The type is always what you declare it as....
	// fmt.Println(config)
}
