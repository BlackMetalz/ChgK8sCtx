# Day 1 - Jan19-2026

### Midnight
- Init the fucking project, add GEMINI and CLAUDE markdown for not generating fucking code. I want to learn, not to be lazy.
- Asked Gemini to gen testdata/kubeconfig for demo purpose. I used to manage 15++ k8s clusters..
- Readfile success, still remember init func that return []byte and error for example, for call it we need to assign it to a variable and check error. (`holyFuck, err := os.ReadFile("./testdata/kubeconfig"`)
- Function or Struct, will always required first letter is upper case for expose or simply mean called from other, don't expose mean private.
- The `_` is the blank identifier - tells Go "I don't care about this value, throw it away."
- The type is always what you declare it as....

### Daytime
- Since I'm able to read the fucking kubeconfig file now, i'm able to parse it to a struct.
- Time to read(list context), lets user pick context, update current context and write back to file!
- Learned something like `strings.Repeat("=",10)` ==> Output: `====================`
- Print specific field of struct is not big issue, I'm already familiar with this 5-6 years ago....

- Next issue about select context, I want feature like use arrow key (Arrow keys (↑↓)) to select context and press enter to switch context, don't want to write context name manually.( Idea by K9S). Searching solution for it and recommended was Terminal UI (TUI). 

- The solution I picked is `https://github.com/manifoldco/promptui/`, because I read this example and I feel I understand LOL `https://github.com/manifoldco/promptui/blob/master/_examples/select/main.go`, it seems to be fit with my requirement.

- Implement it pretty easy: Just add to import then run: `go mod tidy` to download it. And then copy paste the example code to my main.go, change the variable name and run it. It works! Holy fuck!

- Output was amazing:
```bash
kienlt@Luongs-MacBook-Pro ChgK8sCtx % go run main.go
Current context:
dev-cluster
==========All context:==========
dev-cluster
staging-cluster
prod-cluster
aws-eks-cluster
gke-cluster
Use the arrow keys to navigate: ↓ ↑ → ←
? Select k8s cluster:
  ▸ dev-cluster
    staging-cluster
    prod-cluster
    aws-eks-cluster
    gke-cluster
---> After pick a context

✔ staging-cluster
Your selection:  staging-cluster
```

- Show string "(current-context)" next to current context to let user know what is current context
- Added backup after user selection, update current context. 

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