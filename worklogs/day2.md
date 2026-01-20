# Day 2 - Jan20-2026

### Midnight
- I don't want to build 2 separated binary like kubectx and kubens. Just not clone xD. I want to build 1 binary that can do both for learning purpose.
- So i have idea:
    - if user run binary, it will lets user pick 2 options: Change context or Change namespace.
    - if user pick Change context, it will show list of context and let user pick one.
    - if user pick Change namespace, Ask user to enter new default namespace for current context (No validation for now!)
    - Otherwise if user define args like:
        - `go run main.go ns` ==> Ask user to enter new default namespace for current context (No validation for now!).
        - `go run main.go ctx` ==> Give what user see in Day 1.

### Daytime
- Refactor code, create new func loadConfig and saveConfig in file config.go. Little fun, i still remember how that shit work convert to function and why function need to return value xD.

- Now we have to run: `go run .` instead of `go run main.go` because we have multiple files. The `.` mean compile whole fucking package in current folder.

- Hmmm, structure need to be refactor as well. So Gemini recommend a file called `types.go` instead of `structures.go` xD