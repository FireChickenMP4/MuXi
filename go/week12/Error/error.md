# Go 错误处理

> 主要是看[Go 错误处理指北：如何优雅的处理错误？](https://segmentfault.com/a/1190000045397595)

## Go 中为什么没有Exception

Python 的 Exception 异常处理机制是主流编程语言中最为流行的方式，可是 Go 为什么采用了 Error 机制呢？

Go 官方的 [FAQ: Why does Go not have exceptions?](https://go.dev/doc/faq#exceptions) 中给出了解释：

> 我们认为，将异常与控制结构耦合在一起（如 try-catch-finally 语句）会导致代码变得复杂。同时，这也往往会促使程序员将太多普通的错误（比如打开文件失败）标记为异常。
>
> Go 采用了一种不同的处理方式。对于普通的错误处理，Go 函数支持多返回值机制使得在不覆盖返回值的情况下，能够轻松地报告错误。[Go 还提供了一个标准的错误类型，再加上其他特性](https://go.dev/blog/error-handling-and-go)，使得错误处理变得简洁而又与其他语言截然不同。
>
> Go 还提供了一些内置函数，用于标识和恢复真正的异常情况。恢复机制只会在函数状态因错误而被销毁时执行，这足以处理灾难性错误，同时不需要额外的控制结构。使用得当时，可以写出简洁的错误处理代码。
>
> 详情请参考 [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover) 一文。另外，博客文章 [Errors are values](https://go.dev/blog/errors-are-values) 展示了一种整洁的错误处理方式，说明了由于错误只是值，Go 语言的全部能力都可以用于处理错误。

说白了，Go 官方认为 Error 机制更简单有效，且符合 Go 语言大道至简的调性。

## 构造错误

既然要讲解如何处理错误，那么就先从如何构造一个错误说起吧。

我们知道，Go 的 error 实际上就是一个普通的接口，普普通通：

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```

得益于 Go 函数支持多值返回，可以很方便返回一个错误：

```go
func foo() (string, error) {
    return "",nil
}
```

> NOTE:
> 当函数返回多个值时，error 作为最后一个返回值是约定俗成的惯用法。如果你不这么做，代码当然能成功编译，但你有更好的选择。

Go 有两种构造错误方式

```go
err1 := errors.New("I'm err1")
err2 := fmt.Errorf("I'm err%d",2)
```

都是最终返回errorString类型的指针

```go
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```

> NOTE:
> 其实 fmt.Errorf 内部也是调用 errors.New 来创建 error。当然，在 Go 1.13 版本以后，fmt.Errorf 可能会在特定条件下返回 wrapError 类型错误。

## 处理错误

现在我们已经可以构造一个错误，接下来看看如何优雅的处理错误。

### 错误处理方法

怎么说呢，经典...`if err != nil`起手

```go
data, err := foo()
if err != nil {
    //处理错误
    return
}
//正常逻辑
```

一切错误处理都要用`if err != nil`起手1

### Sentinel error

预定义的错误值：`Sentinel error`，一般翻译成`哨兵错误`。这是错误处理惯用法，在Go内置包种有大量应用。
比如：

> https://github.com/golang/go/blob/go1.23.1/src/bufio/bufio.go#L22

```go
var (
    ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
    ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
    ErrBufferFull        = errors.New("bufio: buffer full")
    ErrNegativeCount     = errors.New("bufio: negative count")
)
```

或者：

> https://github.com/golang/go/blob/go1.23.1/src/io/io.go#L29

```go
var ErrShortWrite = errors.New("short write")
var errInvalidWrite = errors.New("invalid write result")
var ErrShortBuffer = errors.New("short buffer")
var EOF = errors.New("EOF")
var ErrUnexpectedEOF = errors.New("unexpected EOF")
var ErrNoProgress = errors.New("multiple Read calls return no data or error")
```

这些都叫 Sentinel error，绝大多数 Sentinel error 都会被定义为包级别公开变量，可以看到也有内置的 errInvalidWrite 并没有对外公开。

每个 `error` 变量都以前缀 `Err` 开头，这是约定俗成的做法。`io.EOF` 是个特例，因为 `EOF` 是另一种约定用法，它的全拼是 end of file，表示文件结束，应用非常广泛，可以算作专有名词了。

我们可以这样处理Sentinel error

```go
if err != nil {
    if err == bufio.ErrBufferFull {
        //处理特定错误
    }
    //处理错误
}
//以及可以用switch case
f, err := os.Open("example.txt")
if err != nil {
    return
}

b := bufio.NewReader(f)

data, err := b.Peek(10)
if err != nil {
    switch err {
    case bufio.ErrNegativeCount:
        // do something
        return
    case bufio.ErrBufferFull:
        // do something
        return
    default:
        // do something
        return
    }
}
fmt.Println(string(data))
```

示例中 b.Peek(10) 可能会返回 ErrNegativeCount 或 ErrBufferFull 错误变量，因为它们是依赖包中可导出的公开变量，所以我们可以在自己的代码中使用这些变量来识别返回了哪个特定的错误消息。

也就是说，这些 Sentinel error 变量会成为包 API 的一部分，用于错误处理。

如果没有 Sentinel error 的存在，我们可能需要通过字符串匹配的方式来识别错误类型：

```go
if err != nil {
    if strings.Contains(err.Error(), "buffer full") {

    }
}
```

_我个人完全不赞成这种写法，不到万不得已，千万不要写成这种代码。_

*记住：error 接口上的 Error 方法适用于人类，而非代码。*只有我们需要查看错误信息，或者记录日志的时候，才应该使用 Error 方法。

此外，你可能在标准库中见到过如下类似代码：

> https://github.com/golang/go/blob/go1.23.1/src/os/error.go#L16

```go
var (
    // ErrInvalid indicates an invalid argument.
    // Methods on File will return this error when the receiver is nil.
    ErrInvalid = fs.ErrInvalid // "invalid argument"

    ErrPermission = fs.ErrPermission // "permission denied"
    ErrExist      = fs.ErrExist      // "file already exists"
    ErrNotExist   = fs.ErrNotExist   // "file does not exist"
    ErrClosed     = fs.ErrClosed     // "file already closed"
)
```

os.ErrInvalid 实际上等价于 fs.ErrInvalid，这种为 Sentinel error 重新赋值的操作也很常见。为了保持良好的分层架构，我们自己的代码设计也可以这样做。

另外，Sentinel error 还有一种看似“另类”的用法，表示错误没有发生，比如 path/filepath.SkipDir：

> https://github.com/golang/go/blob/go1.23.1/src/path/filepath/path.go#L259

```go
// SkipDir is used as a return value from [WalkDirFunc] to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = errors.New("skip this directory")
```

根据注释我们可以了解到，SkipDir 变量用作 WalkDirFunc 的返回值，以指示将跳过调用中指定的目录，它并不表示一个错误。

所以这里 SkipDir 仅作为哨兵，而非错误。其实 io.EOF 也是哨兵，并且它们都没有以 Err 来命名。

这也是我认为 Sentinel error 存在二义性的地方，我个人认为绝大多数情况下不应该这么使用，尽量避免这种用法。

### 常量错误

因为 Sentinel error 是一个变量，所以我们可以随意改变它的值：

```go
oldEOF := io.EOF
io.EOF = errors.New("MyEOF")
fmt.Println(oldEOF == io.EOF) // false
```

这是一个很可怕的事情。

所以 Sentinel error 的确不是一个好的设计，起码也应该将其定义成一个常量。

但问题是在 Go 中我们无法直接将 errors.New 的返回值赋值给一个常量。

如下示例：

```go
const ErrMyEOF = errors.New("MyEOF")
```

这将得到编译报错：

```go
errors.New("MyEOF") (value of type error) is not constant
```

为了解决这个问题，我们可以自定义 error 类型：

```go
type Error string
func (e Error) Error() string { return string(e) }
```

Error 类型底层类型为 string，所以可以直接赋值给一个常量：

```go
const ErrMyEOF = Error("MyEOF")
```

现在常量 ErrMyEOF 不可改变。

但是，这又会引入另外一个新的问题。以下示例代码，执行结果为 true：

```go
const ErrMyEOF = Error("MyEOF")
const ErrNewMyEOF = Error("MyEOF")
fmt.Println(ErrMyEOF == ErrNewMyEOF) // true
```

这与 Go 内置的 errors.New 表现并不相同。

以下示例代码，执行结果为 false：

```go
myEOF = errors.New("EOF")
fmt.Println(io.EOF == myEOF) // false
```

造成二者表现不同的原因是：内置的 errors.New 函数返回 errorString 的指针类型 &errorString{text}，而我们构造的自定义 Error 实际上是 string 类型。

errors.New 返回指针类型是有意而为之的，目的就是在判断两个错误值是否相等时，会比较两个对象是否为同一个对象，而不是比较 Error 方法所返回的字符串内容是否相等。如果仅比较字符串内容是否相等，则我们随便使用 errors.New 函数创建的错误就可以实现与预置的 Sentinel error 相等。

所以常量错误并不常见，我个人其实也不太推荐一定要追求把错误定义为常量，适当引入的编码规范更加切合实际。

尽管 errorString 类型仅包含一个字段 s string，但它还是被有意设计成 struct 而非简单的 string 类型别名，否则 Sentinel error 实用价值将大大折扣。

### 定制错误类型

与使用 errors.New 创建出来的 \*errorString 错误值相比，定制错误类型往往能提供更多的上下文信息。

Go 内置库中就有这样的例子，比如错误类型 os.PathError：

> https://github.com/golang/go/blob/go1.23.1/src/io/fs/fs.go#L250

```go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op   string
    Path string
    Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
```

> NOTE:
> 错误类型命名通常以 Error 结尾，这是约定俗成的惯用法。
> PathError 类型不仅能够记录错误，还会记录导致出现错误的操作和文件路径。在出现错误时，更方便排查问题。

有了新的错误类型后，最大的好处是可以通过类型断言，来判断错误的类型。如果断言成立，则可以根据错误类型对当前错误做更为精细的控制。

示例如下：

```go
// 尝试打开一个不存在的文件
_, err := os.Open("nonexistent.txt")
if err != nil {
    // 使用类型断言检查是否为 *os.PathError 类型
    if pathErr, ok := err.(*os.PathError); ok {
        fmt.Printf("Failed to %s file: %s\n", pathErr.Op, pathErr.Path)
        fmt.Println("Error message:", pathErr.Err)
    } else {
        // 其他类型的错误处理
        fmt.Println("Error:", err)
    }
}
```

可以发现，为了实现错误类型的断言检查，PathError 类型必须是公开类型。

其实无论是 Sentinel error，还是自定义的错误类型，它们都存在同样的问题，都会成为包 API 的一部分，被公开出去。这很可能导致包 API 的快速膨胀。并且，如果代码分层设计不好，很容易出现循环依赖问题。

### Opaque error

Opaque error 是 Go 语言布道师 [Dave Cheney](https://github.com/davecheney) 在 [Gocon Spring 2016](https://dave.cheney.net/paste/gocon-spring-2016.pdf) 演讲中提出的一种叫法，姑且把它翻译为 `不透明的错误处理`。

Opaque error 非常简单，它是最灵活的错误处理策略，因为它需要代码和调用者之间的耦合最少。

示例如下：

```go
import "github.com/quux/"

func fn() error {
        x, err := bar.Foo()
        if err != nil {
                return err
        }
        // use x
}
```

这就是 Opaque error 的全部内容了：只需返回错误，而不对其内容做出任何假设。

没错，遇到错误后直接 return err 的做法就是 Opaque error。

显然，这种代码看似优雅，却过于理想。现实中我们仍有很多情况下还是需要知道错误内容，然后决定是否对其进行处理。

### 错误值比较

比较两个错误值是否相等的操作，一般结合 Sentinel error 一同使用：

```go
if err != nil {
    if err == bufio.ErrBufferFull {
        // handle ErrBufferFull
    }
    // handle err
}
```

先使用 if err != nil 与 nil 比较来判定是否存在错误，如果有错误，更进一步，使用 if err == bufio.ErrBufferFull 来判定错误是否为某个 Sentinel error。

当可能出现多种错误时，还可以使用 switch...case... 来判定错误值：

```go
if err != nil {
    switch err {
    case bufio.ErrNegativeCount:
        // do something
        return
    case bufio.ErrBufferFull:
        // do something
        return
    default:
        // do something
        return
    }
}
```

### 类型断言

Go 支持两种类型断言，[Type Assertion](https://go.dev/doc/effective_go#interface_conversions) 和 [Type Switch](https://go.dev/doc/effective_go#type_switch)。

Go 的类型断言语法可以直接应用于错误处理，因为 error 本身就是一个普通的接口。

断言一个错误的类型，其实前文中我们已经见过了：

```go
// 尝试打开一个不存在的文件
_, err := os.Open("nonexistent.txt")
if err != nil {
    // 使用类型断言检查是否为 *os.PathError 类型
    if pathErr, ok := err.(*os.PathError); ok {
        fmt.Printf("Failed to %s file: %s\n", pathErr.Op, pathErr.Path)
        fmt.Println("Error message:", pathErr.Err)
    } else {
        // 其他类型的错误处理
        fmt.Println("Error:", err)
    }
}
```

如果改用 switch...case... 可以这样写：

```go
// 尝试打开一个不存在的文件
_, err := os.Open("nonexistent.txt")
if err != nil {
    // 使用 switch type 检查错误类型
    switch e := err.(type) {
    case *os.PathError:
        fmt.Printf("Failed to %s file: %s\n", e.Op, e.Path)
        fmt.Println("Error message:", e.Err)
    default:
        // 其他类型的错误处理
        fmt.Println("Error:", err)
    }
}
```

值得一提的是，在使用 Type Switch 语法时，是禁止使用 fallthrough 关键字的，否则编译报错 cannot fallthrough in type switch。

这种情况 case 语句只能使用逗号并提供多个选项：

```go
if err != nil {
    switch err.(type) {
    case *os.PathError, *os.LinkError:
        // do something
    default:
        // do something
    }
}
```

这两种方法的最大缺点就是我们需要导入指定的错误类型，如示例中的 os.PathError 或 os.LinkError。这会导致我们的代码与错误所在的包存在较强的依赖关系。

### 行为断言

随着 Go 语言的演进，大家对 Go 的错误处理又有了新的理解。以前断言错误类型，现在社区中则更推荐断言错误行为。

Go 语言布道师 [Dave Cheney](https://github.com/davecheney) 在他的文章 [Inspecting errors](https://dave.cheney.net/2014/12/24/inspecting-errors) 中提出了断言错误行为而不是类型。

> NOTE:
> 没错，Dave Cheney 的名字再一次出现，后文还会出现😄。这位大佬对 Go 社区的贡献很大，尤其是错误处理，著名的 pkg/errors 包就是他开发的。

```go
func isTimeout(err error) bool {
    type timeout interface {
        Timeout() bool
    }
    te, ok := err.(timeout)
    return ok && te.Timeout()
}
```

函数 isTimeout 用来判定一个错误对象是否表示 Timeout，内部通过断言错误对象是否实现了 timeout 接口来实现。

我们不再假设错误的类型，而是假设其实现了某个接口，并且 timeout 接口是一个临时接口，并不是从其他包中导入的接口类型。这样就真正的实现了包之间的解耦，错误类型无需公开，它们不再必须是包 API 的一部分。

net.Error 就是一个比较不错的实践：

> https://github.com/golang/go/blob/go1.23.1/src/net/net.go#L415

```go
// An Error represents a network error.
type Error interface {
    error
    Timeout() bool // Is the error a timeout?

    // Deprecated: Temporary errors are not well-defined.
    // Most "temporary" errors are timeouts, and the few exceptions are surprising.
    // Do not use this method.
    Temporary() bool
}
```

客户端代码可以断言错误是否为 net.Error 类型，然后再根据行为区分暂时性网络错误和永久性网络错误。

例如，一个爬虫程序在遇到临时错误时可以短暂休眠并重试，否则放弃这个请求，直接处理错误。

示例代码如下：

```go
if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
    time.Sleep(1e9)
    continue
}
if err != nil {
    log.Fatal(err)
}
```

当然也可以写成这样

```go
if nerr, ok := err.(interface{
    Temporary() bool
}); ok {
    time.Sleep(1e9)
    continue
}
if err != nil {
    log.Fatal(err)
}
```

这样就实现了我们的代码与错误所在的包之间最大化的解耦。

当然这两种写法其实都可以，看个人喜好。

```deepseek
你提到的 isTimeout 函数是文章作者为了演示“断言错误行为”而自己编写的示例函数，它并不存在于 Go 标准库或任何特定包里。它的作用是在你的代码中定义一个局部的行为接口（比如 timeout 接口），然后通过类型断言来检查传入的错误是否具有 Timeout() bool 方法，从而判断该错误是否代表超时。

为什么说“不再必须是包 API 的一部分”？
传统上，如果我们想判断一个错误是否属于某种特定类型（例如 *os.PathError），我们就必须导入定义该类型的包，并且该类型必须是公开的（即包 API 的一部分）。这样代码就和这个包产生了强耦合。

而采用“断言行为”的方式时，你可以在自己的代码中定义一个只包含所需方法的本地接口（例如 timeout 接口），然后通过类型断言来检查错误是否实现了这个接口。这个接口完全是你自己定义的，不需要从任何外部包导入。这样一来：

错误的具体类型可以保持私有（不导出），因为它只需要实现方法，而不需要被外部知道。

你的代码不依赖任何外部包的具体类型，只依赖一个行为约定（即接口），大大降低了耦合。

文章随后举了 net.Error 的例子，但请注意：net.Error 本身是 net 包导出的公开接口，它仍然是包 API 的一部分。作者用它来说明“错误可以定义行为接口”，但真正的解耦做法是像 isTimeout 那样，自己定义一个本地接口来断言行为，而不依赖 net.Error 这样的导出接口。

所以，isTimeout 就是你可以在自己的代码中自由编写的函数，它体现了“断言错误行为”的解耦思想。
```

### 暂存错误状态

// to be continue
