# Day 5 - ??Jan-2026

Improvement for CLI Arguments (Flags)
```bash
chg-k8s-ctx -n new_namespace # direct without asked
chg-k8s-ctx --rename old_ctx new_ctx # rename context
chg-k8s-ctx --delete ctx_name # delete ctx
```

Probably using `spf13/cobra` library.