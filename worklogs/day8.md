# Day 8 - 31/01/2026

### Refactor code without using AI xD
Let see how far I can go without AI. LOL

So in `context.go` I clear can see there are 3 func which can go for refactor 
```
- deleteContext
- deleteUser
- deleteCluster
```

In `helper.go`
Group 1:
```
- removeContextByName`
- removeUserByName
- removeClusterByName
```

Group 2:
```
- getOrphanUsers
- getOrphanClusters
```

Approach recommend by Gemini: Generics
URL: https://go.dev/doc/tutorial/generics . Lets read!

yeah, it is definitely what i'm missed. I remember Go generic appeared a lot in previous course i have ever watched. But not really understand that fucking basic somehow. WTF!

Little bit confused about constraint and syntax is not friendly!

Even thought I won't need help of AI to refactor, but in facts I still need AI to explain some fucking basic.

```
func removeByName[T any](items []T, name string, getName func(T) string) []T
//   ↑           ↑      ↑      ↑    ↑            ↑                       ↑
//   1           2      3      4    5            6                       7
```

Remember, they don't have `<` or `>` in syntax. LOL. 

###-Code-Meaning
```
1	removeByName	Function name
2	[T any]	Type parameter: T can be any type
3	items	First parameter name
4	[]T	Type of items = slice of T (example: []Context, []User)
5	name string	Second parameter: string
6	getName func(T) string	Third parameter: function receive T, return string
7	[]T	Return type: slice of T
```

What is `[]T`???
```
[]T = slice of type T

// when T = Context:
[]T → []Context

// when T = User:
[]T → []User

// when T = Cluster:
[]T → []Cluster
```

Compared with non-generic:
```
// Non-generic (3 functions):
func removeContextByName(contexts []Context, name string) []Context
func removeUserByName(users []User, name string) []User
func removeClusterByName(clusters []Cluster, name string) []Cluster

// Generic (1 function):
func removeByName[T any](items []T, name string, getName func(T) string) []T
//                              ↑                                          ↑
//                        []Context/[]User/[]Cluster                 return same type
```

And yeah, I just realized that generics isn't that powerful, we already handle more detail in specific func with removeByName. So `deleteEntryByName` is just simple remove item from slice, nothing more. In fact, we still need to keep exists handler.... But replace loop part with generic func. Therefore `deleteEntryByName` is a low-level helper!

Oh, we can remove `removeContextByName`, `removeUserByName`, `removeClusterByName` helper now. Haha

oh I see function in param, I have seen it in several video course I have been watched, but never understand it!
So it is called Higher-order function pattern. Function that take another function as argument. Python has it also but I never use it LOL.