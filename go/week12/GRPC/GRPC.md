## 包需求

貌似项目需要 `go get google.golang.org/grpc`

## 编译 .proto 文件

proto安装好后，尝试编译你定义好的 .proto 文件，命令如下：
`protoc --go_out=. --go-grpc_out=. \*.proto`
这个路径，最后的这个要改成windows风格。。。貌似

--go*out=. 其中的. 是说你要编译的 .proto 文件目录为当前目录，按需修改
--go-grpc_out=.，其中的. 是说你生成 .pb.go 文件的目录，按需修改
*.proto，其中的 \_ 是说编译当前目录下的所有 .proto 文件，也可以单独指定为 xxx.proto 文件

## 错误问题&解决

protoc命令执行过程中，可能会遇到如下错误：

> PS D:\Work\Code\Go\src\test\proto> protoc --go_out=.
> –go-grpc_out=plugins=grpc:. xxx.proto protoc-gen-go: unable to determine Go import path for “xxx.proto”
> Please specify either:
> • a “go_package” option in the .proto source file, or
> • a “M” argument on the command line.

解决方法：在你的 .proto 文件中，添加如下代码option go_package = "./";

> emm我补充一句，好像这里改成./+包名是直接生成一个文件夹而且包名正确的，比较舒服

```go
// helloworld.proto

syntax = "proto3";

package helloworld;
option go_package = "./";

// 定义请求消息
message HelloRequest {
  string name = 1;
}

// 定义响应消息
message HelloReply {
  string message = 1;
}

// 定义服务
service Greeter {
  // 定义SayHello方法
  rpc SayHello (HelloRequest) returns (HelloReply);
}
```

## Myself

正如此，grpc就是类似协议一样的东西，然后支持go和py
这正好对应服务端和客户端，可以实现跨语言的“类似原生调用api”的体验
甚至说其实本身客户端就不必非要用gin框架，而是作为实打实的客户端，app，exe，cli都可以
_gRPC 的设计初衷就是让任何语言、任何类型的程序都能轻松调用远程服务。_

另外，有经验说，一般不会让客户端和服务端直接rpc通信，而是加了一个http网关，做了一层http转换
回答如此：
除了加访问权限，还有就是加访问权限的方式会更灵活和简单，而且你的外部客户端不一定支持grpc，所以一般用http网关。需不需要做转换其实看你的场景，因为用HTTP做鉴权，加一些中间件会比gRpc方便。

//其实说这玩意，就是我写的这种，最后客户端会是http
//跟直接http调用api差不多。。。

## gRPC-Gateway

什么是 gRPC-Gateway？
[gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) 是一个开源工具，它可以根据你的 Protocol Buffer 定义，自动生成一个反向代理服务器。这个代理服务器会把 HTTP/JSON 请求转换成 gRPC 请求，然后调用后端的 gRPC 服务，并把 gRPC 的响应再转成 JSON 返回给客户端。

简单来说，它让你能够同时提供 gRPC 和 RESTful API，并且只需要维护一套 .proto 接口定义。
// to be continue
