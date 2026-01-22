# Day 3 - Jan22-2026

### Midnight
- Fucking lazy, I haven't made clear progress for this fucking project.
- Creating new function for switch default namespace and include it in main function.
- Migrate switchContext also to new function and include it in main function. Much better than before.
- Start to handle change namespace


### Daytime
- Switch namespace ok, put some data to test works well, it means parse kubeconfig successfully.
```bash
testdata % kubectl config get-contexts
CURRENT   NAME              CLUSTER           AUTHINFO        NAMESPACE
          aws-eks-cluster   aws-eks-cluster   aws-eks-user    default
*         dev-cluster       dev-cluster       dev-admin       kienlt
          gke-cluster       gke-cluster       gke-user        default
          prod-cluster      prod-cluster      prod-admin      production
          staging-cluster   staging-cluster   staging-admin   staging
```

- I'm thinking to create handler for testdata and default kubeconfig (~/.kube/config) to make it easier to switch between them. Because from day 1, i used hardcode kubeconfig for easy test.
- It needs to support both KUBECONFIG and default kubeconfig (~/.kube/config)
- I made it works. Code still trash but it's okay.
- Looking ready to use, but still ask Gemini to review
    - package: I only do it in future if I see it is needed
    - Kubeconfig Path: Gemini point an issue: `return homeDir + "/.kube/config", nil`. This will not works well in Window.
    - Error Handling: Wrap to `fmt.Errorf("wrapper message: %w", err)`. Hmm, interesting...
    - Struct and pointer: +1 time to remember when to use pointer and not to use pointer.
```
   * Gợi ý cải thiện:
       * Khi nào nên truyền struct vào function bằng value (func do(c Config)) và khi nào bằng pointer (func do(c *Config))?
       * Quy tắc chung: Nếu function cần sửa đổi struct đó (ví dụ saveConfig cần thay đổi CurrentContext), bro phải dùng pointer (*Config). Nếu function chỉ đọc dữ liệu từ struct (ví dụ getCurrentContextName), bro
         có thể dùng value.
       * Tuy nhiên: Với các struct lớn, việc truyền bằng pointer sẽ hiệu quả hơn vì Go không phải copy toàn bộ struct mà chỉ copy một địa chỉ memory (rất nhỏ). Hầu hết các Go developer sẽ mặc định dùng pointer cho
         struct để đảm bảo tính nhất quán và hiệu năng.
```