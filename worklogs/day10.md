# Day 10 - Feb 20, 2026

Fix issue appeared after update to use `RunE` instead of `Run`

Before:
```bash
kienlt@kienlt-pc:/data/github/ChgK8sCtx$ go run . --kubeconfig testdata/kubeconfig ctx dev-cluster1
Error: Context dev-cluster1 does not exist
Usage:
  chgctx ctx [context-name] [flags]
  chgctx ctx [command]

Available Commands:
  export      Export a context to stdout (use: ctx export dev > dev.yaml)
  import      Import a context from file
  merge       Merge 2 kubeconfig file
  validate    Validate orphaned context/user/cluster in kubeconfig file

Flags:
      --cleanup          Delete orphan user/cluster in kubeconfig
  -c, --current          Show current context
  -x, --delete           Delete context
      --delete-cluster   Delete cluster in kubeconfig
      --delete-user      Delete user in kubeconfig
  -h, --help             help for ctx
  -l, --list             List all contexts
      --rename           Rename context

Global Flags:
  -d, --debug               Enable debug mode, show debug information xD
      --kubeconfig string   Path to kubeconfig file specified

Use "chgctx ctx [command] --help" for more information about a command.

Context dev-cluster1 does not exist
exit status 1
```

After:
```bash
kienlt@kienlt-pc:/data/github/ChgK8sCtx$ go run . --kubeconfig testdata/kubeconfig ctx dev-cluster1
Context dev-cluster1 does not exist
exit status 1
```