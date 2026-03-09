# High-Level Design: ChgK8sCtx

## Overview

`chgctx` is a CLI tool to switch Kubernetes contexts and namespaces, inspired by `kubectx`/`kubens`. Built in Go as a learning project.

```
User → CLI (Cobra) → Business Logic → kubeconfig file (~/.kube/config)
```

---

## Architecture

### Flat Package Structure

All code lives in `package main` — no sub-packages. Simple project, no need to over-engineer.

```
ChgK8sCtx/
├── main.go          # Entry point
├── root_cmd.go      # Root Cobra command (delegates to ctx logic by default)
├── ctx_cmd.go       # `chgctx ctx` subcommand + all ctx flags
├── ns_cmd.go        # `chgctx ns` subcommand
├── context.go       # Business logic: switch, rename, delete, export, import
├── config.go        # kubeconfig I/O: load, save, path resolution
├── helper.go        # Shared utilities: fuzzyFind, itemExists, color helpers
├── types.go         # KubeConfig, Context, Cluster, User structs
├── namespace.go     # Namespace listing via k8s API
├── merge_cmd.go     # `chgctx merge` subcommand
├── export_cmd.go    # `chgctx export` subcommand
├── validate_cmd.go  # `chgctx validate` subcommand
└── testdata/        # Fixture kubeconfigs for testing
```

---

## Command Routing

```
chgctx <args>
  │
  ├─ matches subcommand? ──yes──► route to subcommand (ns / merge / export / validate)
  │
  └─ no ──► rootCmd.RunE ──► ctxCmd.RunE(args)
                               │
                               ├─ --list / -l       → listContexts()
                               ├─ --current / -c    → showCurrentContext()
                               ├─ --rename          → renameContext()
                               ├─ --delete / -x     → deleteContext()
                               ├─ --delete-user     → deleteUser()
                               ├─ --delete-cluster  → deleteCluster()
                               ├─ --cleanup         → deleteOrphanData()
                               ├─ args[0] == "-"    → switch to previous context
                               ├─ exact match       → switchContext() direct
                               └─ no match          → fuzzyFindContext() → switchContext()
```

> **Root command uses `cobra.ArbitraryArgs`** so any string (e.g. `chgctx aws`, `chgctx gke`) passes through to fuzzy search without Cobra rejecting it as "unknown command".

### Subcommand / Context Name Conflict (Option C)

When a context name equals a subcommand name (e.g., context named `ns`):
- Cobra routes `chgctx ns` to `nsCmd` (subcommand wins — Cobra's design)
- `nsCmd` prints a yellow hint: `⚠ Hint: context 'ns' exists. Use: chgctx ctx ns`
- User can always escape via `chgctx ctx <name>` to force context switch

---

## Data Flow

### kubeconfig I/O (`config.go`)

```
getKubeconfigPath()
  priority: --kubeconfig flag > $KUBECONFIG env > ~/.kube/config

loadConfig(path) → *KubeConfig
  os.ReadFile → yaml.Unmarshal → KubeConfig struct

saveConfig(path, config)
  1. Read current file
  2. Backup to <path>.bak (first time only)
  3. yaml.Marshal(config) → os.WriteFile
```

### Context Switch Flow

```
switchContext(config, path)
  1. promptui.Select → user picks context
  2. savePreviousContext(path, oldCtx)   ← writes to chg-k8s-ctx-history file
  3. config.CurrentContext = newCtx
  4. saveConfig(path, config)
```

### Fuzzy Search (`helper.go`)

Uses `github.com/sahilm/fuzzy`:
```
fuzzyFindContext(config, query) → []string
  - 0 matches  → error: context not found
  - 1 match    → auto-switch (no prompt)
  - 2+ matches → error: ambiguous, list all matches
```

---

## Key Structs (`types.go`)

```go
KubeConfig
  ├── APIVersion     string
  ├── Kind           string
  ├── CurrentContext string
  ├── Contexts       []Context   { Name, Context{Cluster, User, Namespace} }
  ├── Users          []User      { Name, User{...} }
  └── Clusters       []Cluster   { Name, Cluster{Server, CA} }
```

---

## External Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `go.yaml.in/yaml/v3` | YAML parse/marshal |
| `github.com/manifoldco/promptui` | Interactive selection UI |
| `github.com/sahilm/fuzzy` | Fuzzy string matching |
| `k8s.io/client-go` | Namespace listing via K8s API |

---

## Testing Strategy

### Unit Tests

Tests live in `*_test.go` files, same package (`package main`).

- **In-memory fixtures**: `createTestConfig()` in `helper_test.go` returns a `*KubeConfig` struct — no file I/O, no real kubeconfig touched.
- **File I/O tests**: Use files from `testdata/` directory. Pass via `--kubeconfig` flag or `KUBECONFIG` env var — never read from `~/.kube/config`.

```bash
# Run all tests (safe — uses in-memory or testdata only)
go test -v ./...

# Manual smoke test against testdata (NOT real kubeconfig)
KUBECONFIG=testdata/kubeconfig go run . --list
KUBECONFIG=testdata/kubeconfig go run . -c
KUBECONFIG=testdata/kubeconfig go run . dev
```

### Test Data (`testdata/`)

| File | Purpose |
|------|---------|
| `kubeconfig` | Main fixture with multiple contexts |
| `kubeconfig.bak` | Backup fixture |
| `kubeconfig.broken` | For error-path testing |
| `kubeconfig.merge1/2` | For merge command testing |
| `kubeconfig.orphan` | For cleanup/orphan testing |
| `chg-k8s-ctx-history` | Previous context history fixture |

---

## Previous Context (ctx -)

```
~/.kube/chg-k8s-ctx-history   ← single-line file, stores last active context name

chgctx -
  loadPreviousContext()  → reads file
  savePreviousContext()  → writes current ctx BEFORE switching
  switch to stored ctx
```

---

## Notable Design Decisions

| Decision | Reason |
|----------|--------|
| Flat `package main` | Small project, simplicity > structure |
| `cobra.ArbitraryArgs` on root | Allow any string as context name without Cobra rejecting it |
| Backup on first save only | Preserve original; don't overwrite backup on every switch |
| Subcommand wins on name conflict | Cobra's routing is fixed; warn user + offer `chgctx ctx <name>` escape hatch |
| Generic `deleteEntryByName[T]` | Avoid copy-paste for context/user/cluster deletion (Go generics) |
