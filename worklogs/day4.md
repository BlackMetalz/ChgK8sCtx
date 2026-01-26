# Day 4 - 26Jan-2026


### Part 1
Allow user select namespace from current context instead of write manually.

Probably need to interact with client-go and K8S Api xD

Package need to be install:
```bash
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
go mod tidy # Dependencies hell
```

Lession learned:
- use `client-go`, k8s client lib to interact with K8s API
- create client from kubeconfig file
- get all namespace with `clientset.CoreV1().Namespaces().List())`

### Part 2
Improvement for CLI Arguments (Flags)
```bash
chg-k8s-ctx -n new_namespace # direct without asked
chg-k8s-ctx --rename old_ctx new_ctx # rename context
chg-k8s-ctx --delete ctx_name # delete ctx
```

Probably using `spf13/cobra` library.

Gemini recommend this for learning, `cobra-cli`
```
go install github.com/spf13/cobra-cli@latest
```

We need create `cmd` folder, oh this shit sound familiar xD

```go
// Value - create copy each time pass
var nsCmd = cobra.Command{...}  // Type: cobra.Command
rootCmd.AddCommand(&nsCmd)      // Must take address when pass

// Pointer - always reference
var nsCmd = &cobra.Command{...} // Type: *cobra.Command
rootCmd.AddCommand(nsCmd)       // Already pointer, pass directly
```

Go idiom:
```go
// ✅ Common pattern - declare as pointer directly
var nsCmd = &cobra.Command{...}

// ❌ Less common - declare value, then take address
var nsCmd cobra.Command = cobra.Command{...}
rootCmd.AddCommand(&nsCmd)
```

hmmm. There will be big rewrite for this shit since we are migrating to cobra.
First time, Gemini ask me to rewrite, move some shit like config.go, context.go to another internal folder with new package name. 

I thought it is not nescessary, but when ask again, they are from main package and in cmd package, we can not import main package? Oh, i just google for that shit, and yes, we can't import main package. 

Why? Go want a clarify that:
- Code runable, must be run from `main` package.
- Library/Reuseable --> need to be in other package, not fucking `main` package.

So my choice: KISS! --> Keep it simple, stupid! I don't want fancy codes but never gonna reach finish state.

After migrate to cobra, same shit. But benefits are:
1. Auto-generated `--help`:

```bash
go run . --help
go run . ctx --help
go run . ns --help
```

2. Better error msgs
```bash
go run . xyz
Error: unknown command "xyz" for "chg-k8s-ctx"
Run 'chg-k8s-ctx --help' for usage.
unknown command "xyz" for "chg-k8s-ctx"
exit status 1
```

3. Argument validation

in `ns_cmd.go`:
```go
Args: cobra.MaximumNArgs(1)
```

```bash
go run . ns 123 123
Error: accepts at most 1 arg(s), received 2
Usage:
  chg-k8s-ctx ns [namespace-name] [flags]

Flags:
  -h, --help   help for ns

accepts at most 1 arg(s), received 2
exit status 1
```

That argument is for future feature. But build in validation is great!

Added custom validator, can not switch to namespace `kube-system` for funny purpose xD