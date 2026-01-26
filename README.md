# ChgK8sCtx
Change K8S Context

No idea when I finish it or Can I finish it. LOL

# Code Strucutre (Updating....)
```
ChgK8sCtx/
├── main.go        # Entry point, CLI args
├── config.go      # Load/save kubeconfig
├── context.go     # switchContext, switchNamespace
├── namespace.go   # listNamespaces (K8s API)
├── helper.go      # getCurrentContextEntry, handlePromptError
└── types.go       # KubeConfig, Context, Cluster, User structs
```

# Feature/TODO
- make it works with idea of Kubectx and Kubens

# Installation
- I think we need some struct and yaml parser. Recommendation was: `go get go.yaml.in/yaml/v3`
```bash
kienlt@kienlt-pc:/data/ChgK8sCtx$ go get go.yaml.in/yaml/v3
go: downloading go.yaml.in/yaml/v3 v3.0.4
go: added go.yaml.in/yaml/v3 v3.0.4
```

# Structure
Flat structure:
- My project is fucking small
- No need to make it complex or fancy LOL

# for MacOS
- `xattr -d com.apple.quarantine ~/Downloads/chg-k8s-ctx-macos-latest`

# Tunnel from jump
```bash
ssh -N -f -L 8443:kienlt-lab-machine-1:6443 kienlt-lab-fci-jump
```

