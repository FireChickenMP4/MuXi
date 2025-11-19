# 还分不清 Cookie、Session、Token、JWT？看这一篇就够了

## 什么是认证（Authentication）

通俗地讲就是验证当前用户的身份，证明“你是你自己”（比如：你每天上下班打卡，都需要通过指纹打卡，当你的指纹和系统里录入的指纹相匹配时，就打卡成功）

**互联网中的认证：**

## 什么是授权（Authorization）

用户授予第三方应用访问该用户某些资源的权限

实现授权的方式有：cookie、session、token、OAuth

## 什么是凭证（Credentials）

实现认证和授权的前提是需要一种媒介（证书） 来标记访问者的身份

在战国时期，商鞅变法，发明了照身帖。照身帖由官府发放，是一块打磨光滑细密的竹板，上面刻有持有人的头像和籍贯信息。国人必须持有，如若没有就被认为是黑户，或者间谍之类的。

在现实生活中，每个人都会有一张专属的居民身份证，是用于证明持有人身份的一种法定证件。通过身份证，我们可以办理手机卡/银行卡/个人贷款/交通出行等等，这就是认证的凭证。

在互联网应用中，一般网站（如掘金）会有两种模式，游客模式和登录模式。游客模式下，可以正常浏览网站上面的文章，一旦想要点赞/收藏/分享文章，就需要登录或者注册账号。当用户登录成功后，服务器会给该用户使用的浏览器颁发一个令牌（token），这个令牌用来表明你的身份，每次浏览器发送请求时会带上这个令牌，就可以使用游客模式下无法使用的功能。

## 什么是 Cookie

cookie 重要的属性:

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_f3cd85d2706d488a92956d17f49e9b8a.png?x-oss-process=image/resize,w_1400/format,webp)

## 什么是 Session

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_02e077a924a24eb1abd1499c4e990976.png?x-oss-process=image/resize,w_1400/format,webp)

**session 认证流程：**

根据以上流程可知，SessionID 是连接 Cookie 和 Session 的一道桥梁，大部分系统也是根据此原理来验证用户登录状态。

## Cookie 和 Session 的区别

## 什么是 Token（令牌）

### Acesss Token

**特点：**

**token 的身份验证流程：**

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_497bd3c338ab459ca936c4155b041155.png?x-oss-process=image/resize,w_1400/format,webp)

每一次请求都需要携带 token，需要把 token 放到 HTTP 的 Header 里

基于 token 的用户认证是一种服务端无状态的认证方式，服务端不用存放 token 数据。用解析 token 的计算时间换取 session 的存储空间，从而减轻服务器的压力，减少频繁的查询数据库

token 完全由应用管理，所以它可以避开同源策略

### Refresh Token

另外一种 token——refresh token

refresh token 是专用于刷新 access token 的 token。如果没有 refresh token，也可以刷新 access token，但每次刷新都要用户输入登录用户名与密码，会很麻烦。有了 refresh token，可以减少这个麻烦，客户端直接用 refresh token 去更新 access token，无需用户进行额外的操作。

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_a0f8c2eb2add476f951d6c31b4265f9c.png?x-oss-process=image/resize,w_1400/format,webp)

Access Token 的有效期比较短，当 Acesss Token 由于过期而失效时，使用 Refresh Token 就可以获取到新的 Token，如果 Refresh Token 也失效了，用户就只能重新登录了。Refresh Token 及过期时间是存储在服务器的数据库中，只有在申请新的 Acesss Token 时才会验证，不会对业务接口响应时间造成影响，也不需要向 Session 一样一直保持在内存中以应对大量的请求。

## Token 和 Session 的区别

## 什么是 JWT

**JWT 的原理:**

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_7d6133baa25e40139aaf24b7ffdea50f.png?x-oss-process=image/resize,w_1400/format,webp)

**JWT 认证流程：**

```js
Authorization: Bearer;
```

服务端的保护路由将会检查请求头 Authorization 中的 JWT 信息，如果合法，则允许用户的行为

因为 JWT 是自包含的（内部包含了一些会话信息），因此减少了需要查询数据库的需要

因为 JWT 并不使用 Cookie 的，所以你可以使用任何域名提供你的 API 服务而不需要担心跨域资源共享问题（CORS）

因为用户的状态不再存储在服务端的内存中，所以这是一种无状态的认证机制

## JWT 的使用方式

客户端收到服务器返回的 JWT，可以储存在 Cookie 里面，也可以储存在 localStorage。

### 方式一

当用户希望访问一个受保护的路由或者资源的时候，可以把它放在 Cookie 里面自动发送，但是这样不能跨域，所以更好的做法是放在 HTTP 请求头信息的 Authorization 字段里，使用 Bearer 模式添加 JWT。

```js
GET / calendar / v1 / events;

Host: api.example.com;

Authorization: Bearer;
```

用户的状态不会存储在服务端的内存中，这是一种 无状态的认证机制

由于 JWT 是自包含的，因此减少了需要查询数据库的需要

JWT 的这些特性使得我们可以完全依赖其无状态的特性提供数据 API 服务，甚至是创建一个下载流服务。

因为 JWT 并不使用 Cookie ，所以你可以使用任何域名提供你的 API 服务而不需要担心跨域资源共享问题（CORS）

### 方式二

跨域的时候，可以把 JWT 放在 POST 请求的数据体里。

### 方式三

通过 URL 传输

```js
http://www.example.com/user?token=xxx
```

## Token 和 JWT 的区别

**相同：**

**区别：**

**常见的前后端鉴权方式：**

## 常见的加密算法

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_71c5d114d84c4d1d9ba62f63ef9e86db.png?x-oss-process=image/resize,w_1400/format,webp)

哈希算法(Hash Algorithm)又称散列算法、散列函数、哈希函数，是一种从任何一种数据中创建小的数字“指纹”的方法。哈希算法将数据重新打乱混合，重新创建一个哈希值。

哈希算法主要用来保障数据真实性(即完整性)，即发信人将原始消息和哈希值一起发送，收信人通过相同的哈希函数来校验原始数据是否真实。

哈希算法通常有以下几个特点：

【注意】：

1.  以上不能保证数据被恶意篡改，原始数据和哈希值都可能被恶意篡改，要保证不被篡改，可以使用 RSA 公钥私钥方案，再配合哈希值。
2.  哈希算法主要用来防止计算机传输过程中的错误，早期计算机通过前 7 位数据第 8 位奇偶校验码来保障（12.5% 的浪费效率低），对于一段数据或文件，通过哈希算法生成 128bit 或者 256bit 的哈希值，如果校验有问题就要求重传。

## 常见问题

使用 cookie 时需要考虑的问题

使用 session 时需要考虑的问题

使用 token 时需要考虑的问题

## 使用 JWT 时需要考虑的问题

因为 JWT 并不依赖 Cookie 的，所以你可以使用任何域名提供你的 API 服务而不需要担心跨域资源共享问题（CORS）

## 使用加密算法时需要考虑的问题

## 分布式架构下 session 共享方案

### 1\. session 复制

任何一个服务器上的 session 发生改变（增删改），该节点会把这个 session 的所有内容序列化，然后广播给所有其它节点，不管其他服务器需不需要 session ，以此来保证 session 同步

优点：可容错，各个服务器间 session 能够实时响应。

缺点：会对网络负荷造成一定压力，如果 session 量大的话可能会造成网络堵塞，拖慢服务器性能。

### 2\. 粘性 session /IP 绑定策略

采用 Ngnix 中的 ip_hash 机制，将某个 ip 的所有请求都定向到同一台服务器上，即将用户与服务器绑定。用户第一次请求时，负载均衡器将用户的请求转发到了 A 服务器上，如果负载均衡器设置了粘性 session 的话，那么用户以后的每次请求都会转发到 A 服务器上，相当于把用户和 A 服务器粘到了一块，这就是粘性 session 机制。

优点：简单，不需要对 session 做任何处理。

缺点：缺乏容错性，如果当前访问的服务器发生故障，用户被转移到第二个服务器上时，他的 session 信息都将失效。

适用场景：发生故障对客户产生的影响较小；服务器发生故障是低概率事件 。

实现方式：以 Nginx 为例，在 upstream 模块配置 ip_hash 属性即可实现粘性 session。

### 3\. session 共享（常用）

使用分布式缓存方案比如 Memcached 、Redis 来缓存 session，但是要求 Memcached 或 Redis 必须是集群

把 session 放到 Redis 中存储，虽然架构上变得复杂，并且需要多访问一次 Redis ，但是这种方案带来的好处也是很大的：

![image.png](https://ucc.alicdn.com/pic/developer-ecology/pmur6hy3nphhs_82ba56eeb77c42699074926c44764298.png?x-oss-process=image/resize,w_1400/format,webp)

### 4\. session 持久化

将 session 存储到数据库中，保证 session 的持久化

优点：服务器出现问题，session 不会丢失

缺点：如果网站的访问量很大，把 session 存储到数据库中，会对数据库造成很大压力，还需要增加额外的开销维护数据库。

只要关闭浏览器 ，session 真的就消失了？

不对。对 session 来说，除非程序通知服务器删除一个 session，否则服务器会一直保留，程序一般都是在用户做 log off 的时候发个指令去删除 session。

然而浏览器从来不会主动在关闭之前通知服务器它将要关闭，因此服务器根本不会有机会知道浏览器已经关闭，之所以会有这种错觉，是大部分 session 机制都使用会话 cookie 来保存 session id，而关闭浏览器后这个 session id 就消失了，再次连接服务器时也就无法找到原来的 session。如果服务器设置的 cookie 被保存在硬盘上，或者使用某种手段改写浏览器发出的 HTTP 请求头，把原来的 session id 发送给服务器，则再次打开浏览器仍然能够打开原来的 session。

恰恰是由于关闭浏览器不会导致 session 被删除，迫使服务器为 session 设置了一个失效时间，当距离客户端上一次使用 session 的时间超过这个失效时间时，服务器就认为客户端已经停止了活动，才会把 session 删除以节省存储空间。
