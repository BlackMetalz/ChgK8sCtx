# TODO - ChgK8sCtx

## Refactoring

### 1. Named Types for Nested Structs
**Priority: Medium**

Currently using anonymous structs in `types.go`:
```go
type Context struct {
    Name    string
    Context struct {  // Anonymous - verbose in tests!
        Cluster   string
        User      string
        Namespace string
    }
}
```

**Should refactor to named types:**
```go
type ContextDetails struct {
    Cluster   string `yaml:"cluster"`
    User      string `yaml:"user"`
    Namespace string `yaml:"namespace"`
}

type Context struct {
    Name    string         `yaml:"name"`
    Context ContextDetails `yaml:"context"`
}
```

**Benefits:**
- Cleaner test code
- Reusable type
- Better IDE autocomplete

**Files affected:**
- `types.go` - Define named types
- `helper_test.go` - Simplify createTestConfig()
- `validate_test.go` - Simplify test configs
- `merge_test.go` - Simplify test configs
- `context_test.go` - Simplify test configs

---

### 2. DRY - Remove/Cluster/User/Context functions
**Priority: Low**

Functions `removeContextByName`, `removeUserByName`, `removeClusterByName` are almost identical.
Could use generics (Go 1.18+) or interface-based approach.

---

### 3. Error handling consistency
**Priority: Low**

Some functions return `error`, some print and return. 
Consider consistent error handling pattern.

---

## Future Features

### From backlog.md
- [ ] Shell Prompt Integration
- [ ] More comprehensive validation rules
