## Go 的切片：

### 在功能和使用场景上：

> **Go 的切片（slice） ≈ C++ 的 `std::vector<T>`**  
> 而 **Go 的数组（array） ≈ C++ 的原生数组（如 `int arr[5]`）或 `std::array<T, N>`**

---

### 🆚 详细对比

| 特性           | Go 数组 (`[N]T`)       | Go 切片 (`[]T`)               | C++ 原生数组 (`T arr[N]`)          | C++ `std::array<T, N>`    | C++ `std::vector<T>`                 |
| -------------- | ---------------------- | ----------------------------- | ---------------------------------- | ------------------------- | ------------------------------------ |
| 大小固定？     | ✅ 是                  | ❌ 否（动态）                 | ✅ 是                              | ✅ 是                     | ❌ 否（动态）                        |
| 值语义？       | ✅ 赋值/传参时拷贝全部 | ✅ 切片头是值，但共享底层数组 | ✅（但作为函数参数会退化为指针！） | ✅ 完整拷贝               | ✅ vector 对象是值，但内部管理堆内存 |
| 底层存储       | 栈 or 静态区           | 底层数组在堆上（通常）        | 栈 or 静态区                       | 栈（通常）                | 堆（元素在堆上）                     |
| 可变长度？     | ❌                     | ✅（通过 `append`）           | ❌                                 | ❌                        | ✅（通过 `push_back`, `resize` 等）  |
| 自动扩容？     | ❌                     | ✅（`append` 触发）           | ❌                                 | ❌                        | ✅（内部自动 realloc）               |
| 安全边界检查？ | ✅（运行时 panic）     | ✅（运行时 panic）            | ❌（无）                           | ✅（`.at()` 有，`[]` 无） | ✅（`.at()` 有，`[]` 无）            |

---

### 🔍 为什么说 **Go 切片 ≈ C++ `std::vector`**？

1. **动态大小**：
   - `vector` 可以 `push_back`，Go 切片可以 `append`
2. **自动内存管理**：
   - 两者都自动处理底层内存分配和扩容
3. **引用语义（共享数据需显式）**：
   - `vector` 赋值是深拷贝（除非用引用）
   - Go 切片赋值是“浅拷贝”（共享底层数组），但这是设计选择，类似 `vector` 的迭代器或指针行为
4. **容量（capacity）概念**：
   - `vector` 有 `.size()` 和 `.capacity()`
   - Go 切片有 `len()` 和 `cap()`

> 💡 实际上，Go 切片的设计灵感部分来自像 `vector` 这样的动态数组容器。

---

### ⚠️ 重要区别

虽然功能相似，但语义上有关键不同：

| 方面         | Go 切片                                                                | C++ `std::vector`                         |
| ------------ | ---------------------------------------------------------------------- | ----------------------------------------- |
| **默认行为** | 赋值/传参 → 共享底层数组（浅拷贝 header）                              | 赋值/传参 → 深拷贝整个 vector（包括数据） |
| **扩容影响** | 如果多个切片共享底层数组，一个 `append` 可能导致“意外不共享”（因扩容） | `vector` 是独立对象，扩容不影响其他副本   |
| **空状态**   | `nil` slice vs empty slice（行为略有不同）                             | `vector` 默认构造即为空，无 `nil` 概念    |

#### 示例：C++ vector 是深拷贝

```cpp
std::vector<int> a = {1, 2, 3};
std::vector<int> b = a;      // 深拷贝！
b[0] = 99;
// a 仍是 {1,2,3}
```

而 Go 切片是“浅拷贝 header”：

```go
a := []int{1, 2, 3}
b := a        // 浅拷贝 slice header，共享底层数组
b[0] = 99
// a 变成 [99, 2, 3]！
```

> 所以严格来说：**Go 切片的行为更像 C++ 中“指向 vector 内部的指针 + 长度信息”，而不是 vector 本身**。  
> 但从**用途**上看，程序员用切片来实现动态数组，正如用 `vector` 一样。

---

### ✅ 总结

- **Go 的数组** ↔ **C/C++ 的固定大小数组**（如 `int[5]` 或 `std::array`）
- **Go 的切片** ↔ **C++ 的 `std::vector`**（在用途和动态性上）

但要注意：

- Go 切片的**共享语义**比 `vector` 更“隐式”（容易意外共享或意外断开共享）
- C++ `vector` 更强调**值语义**（拷贝即独立）

所以你在编程时：

- 想要动态、可增长的序列？→ 用 **Go 切片**（就像 C++ 用 `vector`）
- 想要固定大小、栈分配、值语义的结构？→ 用 **Go 数组**（就像 C++ 用 `std::array`）

## make 和 new 的区别 来自 the way to go

> 我其实还是有点不太理解这一段想说什么，但是大概是这么个事

---

### 7.2.4 new() 和 make() 的区别

> 看起来二者没有什么区别，都在堆上分配内存，但是它们的行为不同，适用于不同的类型。

- new(T) 为每个新的类型 T 分配一片内存，初始化为 0 并且返回类型为`\*T` 的内存地址：这种方法 返回一个指向类型为 T，值为 0 的地址的指针，它适用于值类型如数组和结构体（参见第 10 章）；它相当于 `&T{}`。
- make(T) 返回一个类型为 T 的初始值，它只适用于 3 种内建的引用类型：切片、map 和 channel（参见第 8 章，第 13 章）。

换言之，new 函数分配内存，make 函数初始化；下图给出了区别：

在图 7.3 的第一幅图中：

```go
var p *[]int = new([]int) // *p == nil; with len and cap 0
p := new([]int)
```

在第二幅图中， `p := make([]int, 0) `，切片 已经被初始化，但是指向一个空的数组。

以上两种方式实用性都不高。下面的方法：

```go
var v []int = make([]int, 10, 50)
```

或者

```go
v := make([]int, 10, 50)
```

这样分配一个有 50 个 int 值的数组，并且创建了一个长度为 10，容量为 50 的 切片 v，该 切片 指向数组的前 10 个元素。

---

### 7.2.6 bytes 包

类型 `[]byte `的切片十分常见，Go 语言有一个 bytes 包专门用来解决这种类型的操作方法。

bytes 包和字符串包十分类似（参见第 4.7 节）。而且它还包含一个十分有用的类型 Buffer:

```go
import "bytes"

type Buffer struct {
    ...
}
```

这是一个长度可变的 bytes 的 buffer，提供 Read 和 Write 方法，因为读写长度未知的 bytes 最好使用 buffer。

Buffer 可以这样定义：

```go
var buffer bytes.Buffer。
```

或者使用 new 获得一个指针：

```go
var r *bytes.Buffer = new(bytes.Buffer)。
```

或者通过函数：

```go
func NewBuffer(buf []byte) *Buffer
```

创建一个 Buffer 对象并且用 buf 初始化好；NewBuffer 最好用在从 buf 读取的时候使用。

通过 buffer 串联字符串

类似于 Java 的 StringBuilder 类。

在下面的代码段中，我们创建一个 buffer，通过`buffer.WriteString(s)`方法将字符串 s 追加到后面，最后再通过`buffer.String() `方法转换为 string：

```go
var buffer bytes.Buffer
for {
    if s, ok := getNextString(); ok { //method getNextString() not shown here
        buffer.WriteString(s)
    } else {
        break
    }
}
fmt.Print(buffer.String(), "\n")
```

这种实现方式比使用 += 要更节省内存和 CPU，尤其是要串联的字符串数目特别多的时候。

### 7.6.8 切片和垃圾回收

---

切片的底层指向一个数组，该数组的实际体积可能要大于切片所定义的体积。只有在没有任何切片指向的时候，底层的数组内层才会被释放，这种特性有时会导致程序占用多余的内存。

示例 函数 `FindDigits` 将一个文件加载到内存，然后搜索其中所有的数字并返回一个切片。

```go
var digitRegexp = regexp.MustCompile("[0-9]+")

func FindDigits(filename string) []byte {
b, \_ := ioutil.ReadFile(filename)
return digitRegexp.Find(b)
}
```

这段代码可以顺利运行，但返回的 `[]byte` 指向的底层是整个文件的数据。只要该返回的切片不被释放，垃圾回收器就不能释放整个文件所占用的内存。换句话说，一点点有用的数据却占用了整个文件的内存。

想要避免这个问题，可以通过拷贝我们需要的部分到一个新的切片中：

```go
func FindDigits(filename string) []byte {
b, \_ := ioutil.ReadFile(filename)
b = digitRegexp.Find(b)
c := make([]byte, len(b))
copy(c, b)
return c
}
```
