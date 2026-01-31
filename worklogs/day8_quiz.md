# Fucking 🧪 Quiz Time!
# Just keep vietnamese here. LOL

Câu 1:
Trong function signature này:

```
func deleteEntryByName[Entry any](<items []Entry, name string, getName func(Entry>) string) []Entry
```

Nếu bro gọi với config.Users (slice của User struct), thì Entry sẽ là gì?



Câu 2:
Tại sao không thể viết như này:

```
func deleteEntryByName[Entry any](<items []Entry, name string>) []Entry {
    for _, item := range items {
        if item.Name != name {  // ← Tại sao lỗi?
            // ...
        }
    }
}
```

Câu 3:
Giải thích đoạn code này làm gì:
```
func(c Context) string { return c.Name }
```

Trả lời, và đây là feedback:

Câu 1: ✅ Đúng!
Entry sẽ là  User struct. Go compiler sẽ tự suy ra từ config.Users là []User → Entry = User.

Và đúng, [T any] hoặc [Entry any] là dấu hiệu nhận biết generic function!


Câu 2: Giải thích! 📚
```
// Lỗi này:
if item.Name != name  // ❌ ERROR!
```

Tại sao lỗi?

Vì `Entry any` nghĩa là Entry có thể là BẤT KỲ TYPE NÀO - kể cả int, string, hay bất kỳ struct nào.
Go compiler không biết Entry có field Name hay không!
```
// Ví dụ: Entry có thể là:
type User struct { Name string }     // ← Có .Name ✅
type Cluster struct { Name string }  // ← Có .Name ✅
type int                              // ← KHÔNG có .Name ❌
type string                           // ← KHÔNG có .Name ❌
```

Vì any quá general, Go không dám giả định item.Name tồn tại!
Giải pháp: Dùng getName function để CHỈ CHO GO cách lấy name:

```
getName(item)  // ← Gọi function để lấy name, an toàn!
```


Câu 3: ✅ Đúng!
Đây là anonymous function (function không tên), nhận Context struct, return string là c.Name.

# Feedback
Score: 2.5/3 🎯
Bro đã hiểu concept chính! Câu 2 là phần khó nhất - đó chính là lý do tại sao Generics trong Go cần "helper function" để access fields! 💪

