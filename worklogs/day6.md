# Day 6 - 29Jan-2026

### Implement Fuzzy search
This look fun to implement and nice to know that there is a library for fuzzy search in Go (1k4 stars)

```bash
go get github.com/sahilm/fuzzy
```

And yeah, it should be support mode direct switch to context, not interactive mode.

### Implement Switch Back
It is pure logic handles, just remember when you switch context, you need to save not only context, you need to save for history
And this time I think only support last context is great, no need to support multi history context.

### Implement kubeconfig flag
- Copy idea of kubectl, priority order: nothing much to explain xD
    1. --kubeconfig flag
    2. KUBECONFIG env var
    3. Default kubeconfig

### Implement merge command
Simple merge logic, just merge 2 files and save to output file xD