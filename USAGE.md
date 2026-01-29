# ChgK8sCtx - Usage Guide

A `kubectx`/`kubens` clone for learning Go.

## Quick Start

```bash
# Build
go build -o chg-k8s-ctx .

# Or run directly
go run . <command>
```

## Commands

### Context (`ctx`)

| Command | Description |
|---------|-------------|
| `ctx` | Interactive context switch |
| `ctx <name>` | Direct switch to context (supports fuzzy search) |
| `ctx -` | Switch to previous context (like `cd -`) |
| `ctx -l` | List all contexts |
| `ctx -c` | Show current context |
| `ctx --rename` | Rename a context (interactive) |
| `ctx -x` | Delete a context (interactive) |
| `ctx --delete` | Delete a context (interactive) |
| `ctx --delete-user` | Delete user (cascade deletes related contexts) |
| `ctx --delete-cluster` | Delete cluster (cascade deletes related contexts) |
| `ctx --cleanup` | Delete orphan users/clusters |

#### Fuzzy Search
```bash
ctx dev          # Exact match: dev-cluster
ctx gke          # Fuzzy match: gke-cluster  
ctx aws          # Fuzzy match: aws-eks-cluster
ctx cluster      # Multiple matches → shows list
```

### Namespace (`ns`)

| Command | Description |
|---------|-------------|
| `ns` | Interactive namespace switch |
| `ns <name>` | Direct switch to namespace |

### Global Flags

| Flag | Description |
|------|-------------|
| `-d, --debug` | Enable debug output |

---

## Test Data

### Base Files (do not modify directly)
- `testdata/kubeconfig.original` - Clean config for general testing
- `testdata/kubeconfig.orphan` - Config with orphan users/clusters

### Test Pattern
```bash
# Reset before each test
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig
```

---

## Runnable Test Scripts

### Setup
```bash
cd /Users/kienlt/data/github.com/ChgK8sCtx
go build -o chg-test .
```

### 1. Context Switch (Exact + Fuzzy)
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx -c                 # Current: staging-cluster
./chg-test ctx dev-cluster        # Exact match
./chg-test ctx -c                 # Verify: dev-cluster
./chg-test ctx prod               # Fuzzy match → prod-cluster
./chg-test ctx -c                 # Verify: prod-cluster
```

### 2. Context History (`ctx -`)
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx dev-cluster        # Switch to dev
./chg-test ctx prod-cluster       # Switch to prod (saves dev)
./chg-test ctx -                  # Switch back to dev
./chg-test ctx -c                 # Verify: dev-cluster
./chg-test ctx -                  # Toggle back to prod
./chg-test ctx -c                 # Verify: prod-cluster
```

### 3. Fuzzy Search Edge Cases
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx gke                # Fuzzy: gke-cluster
./chg-test ctx aws                # Fuzzy: aws-eks-cluster
./chg-test ctx cluster            # Multiple matches → shows list
./chg-test ctx nonexistent        # Error: not found
```

### 4. Cleanup Orphans
```bash
cp testdata/kubeconfig.orphan testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx -l                 # List all contexts
./chg-test ctx --cleanup          # Shows orphan users/clusters
# Expected:
# Found 3 orphan users: orphan-user-1, orphan-user-2, old-admin
# Found 2 orphan clusters: orphan-cluster-1, orphan-cluster-2
```

### 5. Error Cases
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx staging-cluster    # Switch first
./chg-test ctx staging-cluster    # Already on context → warning
./chg-test ctx notfound           # Context not found → error
```

### 6. Rename Context
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx -l                           # List before
./chg-test ctx --rename dev-cluster my-dev  # Direct rename
./chg-test ctx -l                           # Verify renamed
./chg-test ctx --rename                     # Interactive rename
```

### 7. Delete Context
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx -l                 # List before
./chg-test ctx -x                 # Interactive delete
./chg-test ctx --delete           # Interactive delete (alias)
./chg-test ctx -l                 # Verify deleted
```

### 8. Cascade Delete (User/Cluster)
```bash
cp testdata/kubeconfig.original testdata/kubeconfig
export KUBECONFIG=testdata/kubeconfig

./chg-test ctx --delete-user      # Shows warning, confirms cascade
# Expected: "Warning: Will ALSO delete context(s): ..."
# Confirms: "Delete this user AND related contexts?"

./chg-test ctx --delete-cluster   # Same pattern for cluster
# If deleting current context's user/cluster → auto-switches
```

---

## Debug Mode
```bash
go run . -d ctx           # Debug + context
go run . --debug ns       # Debug + namespace
```
