# Day 5 - 27Jan-2026

Finish CLI Args

- Change namespace without selection
```bash
chg-k8s-ctx ns new-namespace
```

Yes, cobra cli rock, no need to implement much, just fill
```bash
go run . ns --help
Switch to a different kubernetes namespace. 
Examples:
  # Interactive mode - show selection menu
  chg-k8s-ctx ns
  
  # Direct mode - switch immediately  
  chg-k8s-ctx ns kube-public
  chg-k8s-ctx ns my-namespace

Usage:
  chg-k8s-ctx ns [namespace] [flags]

Examples:
  chg-k8s-ctx ns
  chg-k8s-ctx ns kube-public

Flags:
  -h, --help   help for ns
```


- Rename context. Approach: KISS xD
So we would go for simple and stupid implementation
```bash
ctx --rename # Interactive mode

ctx --rename old-name new-name # Direct mode if we know old name and new name. Need validation for this!
```

A lot of joys while implement this feature. For the first time, i prove the code from AI is wrong xDD
![alt text](../images/day5.png)