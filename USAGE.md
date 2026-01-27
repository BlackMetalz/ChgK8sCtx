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
| `ctx <name>` | Direct switch to context |
| `ctx -l` | List all contexts |
| `ctx -c` | Show current context |
| `ctx --rename` | Rename a context (interactive) |
| `ctx -x` | Delete a context (interactive) |
| `ctx --delete-user` | Delete user (cascade deletes related contexts) |
| `ctx --delete-cluster` | Delete cluster (cascade deletes related contexts) |
| `ctx --cleanup` | Delete orphan users/clusters |

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

```bash
# Use test kubeconfig
export KUBECONFIG=testdata/kubeconfig

# For orphan cleanup testing
export KUBECONFIG=testdata/kube_orphan
```

### testdata/kubeconfig
- 2 contexts: `dev-cluster`, `prod-cluster`
- Clean config for general testing

### testdata/kube_orphan
- 2 contexts: `dev-cluster`, `prod-cluster`
- 3 orphan users: `orphan-user-1`, `orphan-user-2`, `old-admin`
- 2 orphan clusters: `orphan-cluster-1`, `orphan-cluster-2`

---

## Manual Test Cases

### 1. Context Switch
```bash
go run . ctx              # Interactive
go run . ctx dev-cluster  # Direct
go run . ctx -c           # Verify current
```

### 2. Namespace Switch
```bash
go run . ns               # Interactive
go run . ns kube-system   # Direct
```

### 3. Rename Context
```bash
go run . ctx --rename     # Select & rename
go run . ctx -l           # Verify
```

### 4. Delete Context
```bash
go run . ctx -x           # Interactive delete
go run . ctx -l           # Verify deleted
```

### 5. Delete User (Cascade)
```bash
go run . ctx --delete-user
# Should show warning if user is used by contexts
# Confirms: "Delete user AND related contexts?"
```

### 6. Delete Cluster (Cascade)
```bash
go run . ctx --delete-cluster
# Should show warning if cluster is used by contexts
# Confirms: "Delete cluster AND related contexts?"
```

### 7. Cleanup Orphans
```bash
export KUBECONFIG=testdata/kube_orphan
go run . ctx --cleanup
# Expected output:
# Found 3 orphan users: orphan-user-1, orphan-user-2, old-admin
# Found 2 orphan clusters: orphan-cluster-1, orphan-cluster-2
# Delete ALL orphan items? [Yes/No]
```

---

## Debug Mode
```bash
go run . -d ctx           # Debug + context
go run . --debug ns       # Debug + namespace
```
