package main

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
