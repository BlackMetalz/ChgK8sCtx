# Inspired of Agile.
Just throw issue to backlog section. LOL
Backlog: Issue list that not gonna have a date to start/finish


### 1. Fuzzy Search (like fzf)
```bash
ctx dev     # matches: dev-cluster, dev-staging, my-dev-env
```
Library: github.com/sahilm/fuzzy

### 2. Context History / Switch Back

```bash
ctx -        # Switch back to previous context (like `cd -`)
```
Store last context in file/memory.

### 3. Context Aliases
```bash
ctx alias prod aws-eks-prod-cluster-very-long-name
ctx prod     # Switch to aliased context
```

### 4. Multi-kubeconfig Support
```bash
ctx --kubeconfig ~/.kube/config-work
ctx merge config1 config2 --output merged-config
```

### 5. Shell Prompt Integration???
```bash
ctx prompt   # Returns: "dev-cluster/kube-system"
```

### 6. Config Validation
```bash
ctx validate  # Check for broken references
# Output: "User 'old-admin' referenced by context but doesn't exist"
```

### 7. Export/Import Context
```bash
ctx export dev-cluster > dev.yaml
ctx import dev.yaml
```

### 8. Unit Tests (Go learning!)
Write tests for helper functions - learn testing package!

🔧 Refactoring Ideas:
1. Extract duplicate "context list building" → helper
2. Use errors.New() instead of fmt.Errorf() for static errors
3. Implement interfaces for testability

### No category
https://hashir.blog/2025/06/powerlevel10k-is-on-life-support-hello-starship/