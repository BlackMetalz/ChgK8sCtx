package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui" // Promptui for interactive prompt
	"go.yaml.in/yaml/v3"             // Import this shit for unmarshal yaml
)

// init fucking struct for yaml?
type KubeConfig struct {
	APIVersion     string    `yaml:"apiVersion"`  // Need this shit for correct kubeconfig file
	Kind           string    `yaml:"kind"`        // Need this shit for correct kubeconfig file
	Preferences    struct{}  `yaml:"preferences"` // empty struct for empty object
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
	// fmt.Println("Hello, World!")

	// Look at the func ReadFile. We don't need to talk about param
	// Only understand this: ([]byte, error)
	// fmt.Println("Try to read fucking file.")

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
	//fmt.Printf("Type of config: %T\n", config)
	// Output: Type of config: main.KubeConfig
	// The type is always what you declare it as....
	// fmt.Println(config)

	// List context of K8S
	fmt.Println("Current context:", config.CurrentContext)

	// List all context and append all cluster name to a slice of string
	// Init default slide of string
	var contextList []string
	// fmt.Println(strings.Repeat("=", 10) + "All context:" + strings.Repeat("=", 10))
	for _, ctx := range config.Contexts {
		// fmt.Println(abc.Name) // No need to print this shit.
		// Append to slice
		// Add a text: (current-context) next to current context to let user know this is current fucking context.
		if ctx.Name == config.CurrentContext {
			// We will need to remove this text when we write back to file in next step.
			contextList = append(contextList, ctx.Name+"(current-context)")
		} else {
			contextList = append(contextList, ctx.Name)
		}
	}

	// Ok, we have fucking promptui
	prompt := promptui.Select{
		Label: "Select k8s cluster",
		Items: contextList,
	}

	// Run fucking prompt
	_, result, err := prompt.Run()

	// Go idiom style
	switch err {
	case nil:
		// success, do nothing
	case promptui.ErrInterrupt:
		fmt.Println("Interrupted!")
		return
	case promptui.ErrEOF:
		fmt.Println("Cancelled!")
		return
	default:
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Check if user select same context
	if result == config.CurrentContext+"(current-context)" {
		fmt.Println("You already selected this context!")
		return
	}

	// if err != nil {
	// 	if err == promptui.ErrInterrupt {
	// 		// User click Ctrl+C
	// 		fmt.Println("Interrupted by Ctrl+C! Exiting...")
	// 	} else if err == promptui.ErrEOF {
	// 		// User click Ctrl+D or ESC
	// 		// Test in Mac, ESC doesn't trigger this. Only Ctrl+D works
	// 		fmt.Println("Cancelled by Ctrl+D or ESC! Exiting...")
	// 	} else {
	// 		// Other errors
	// 		fmt.Printf("Prompt failed: %v\n", err)
	// 	}
	// 	return
	// }

	fmt.Println("Your selection: ", result)

	// Write back to file section //
	// Remove (current-context) from result. StrimSuffix first, update struct later on.
	result = strings.TrimSuffix(result, "(current-context)")
	// Update struct
	config.CurrentContext = result

	// Marshal struct -> bytes
	newData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error marshalling in section write back to file")
		return
	}

	// Oke, next step backup current file to .bak extension, overwrite if exist for simplicity
	// This should be fine even we run multiple times.
	// And yeah, only backup after user select the context.
	err = os.WriteFile("./testdata/kubeconfig.bak", _data, 0644)
	if err != nil {
		fmt.Println("Warning: Failed to create backup:", err)
		// có thể continue hoặc return tùy bạn
	} else {
		fmt.Println("Current kubeconfig file was backed up to ./testdata/kubeconfig.bak")
	}

	// Write to file
	err = os.WriteFile("./testdata/kubeconfig", newData, 0600)
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}

	fmt.Println("Successfully switched to context: ", result)

}
