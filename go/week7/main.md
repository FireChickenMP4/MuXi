# Week7

> 本周主要是 linux(虽然我也还不是很会),以及 docker

主要写一下可能的主线小游戏的流程(毕竟是在 docker 内跑的,用 md 记录会好一些)

顺便复习一下（，很早之前就开始折腾 wsl 和 docker 了，但是没做啥实事，所以忘得差不多了，用小游戏练手当巩固了，练着练着估计也就会了

然后之前有个好玩的命令行图形化的 docker，lazydocker，也是比较方便了

---

```bash
root@fbfe32944d56:/home# ./start
“请以 sudo 面对它…以招来祂的注视！”
root@fbfe32944d56:/home# sudo ./start
请输入你的姓名： FireChickenMP4
游戏开始：FireChickenMP4，希望能给你带来一场难忘的体验。
11月29日，夜深了，工位的灯光只剩下幽暗的黄光，像是被时间吞噬的一缕残影。我坐 在显示器前，手里攥着前辈留给我的笔记本，心跳莫名地加快。
前辈——那个沉默而执着的人——就在三天前猝然离职。没有预兆，也没有留下任何告别。 他离职前唯一在意的，便是那个被他称作“深处”的 Docker 容器。
没人知道那意味着什么，也没有人敢随便去启动它。前辈在最后的留言里写道：
“如果有人要接手这份工作……请小心。它不是普通的容器，它……观察你，也在等待。”
我深吸一口气，手指悬在键盘上，眼前的终端里只有那行冷漠的命令提示符。
这是我的任务——去启动它，去理解它，去面对前辈没能承受的东西。



你感到有什么东西在看着你,可能是你的错觉???


Tips: 1.多尝试查看各种文件，2.减少删除和写入操作，以及不要尝试重启容器，避免 死档。3.如果遇到一些奇怪的现象可能是为了氛围感特意营造的，不要感到恐慌
```

我卡了,光速
然后随便闯了，用了 docker inspect
发现有个`while true; do cat /bin/containerexec/exec.log; done`
打印如下

```bash
1月16日：
自从看到那串 bXV4aQ== 之后，我的梦境开始被某种东西侵蚀...
我疯狂的想要找到它，它在容器的各个角落，但是又不在任何地方，或许我只是一个被 祂玩弄的孩子，狂妄的想要追寻祂的踪迹，以求得到一点希望。~~~或许在文件的某一处可以找到bXV4aQ==的真身~~~
```

...这 base64 吧，解码 muxi(传统艺能了)

所以我说先去 bin 闯一闯吧

没用。。。

但是有什么在看我找到了

甚至开到盒了（，在 tmp 源码的 2 节点 main 程序

...所以我现在要干什么嘞。。。

再试着跑一下这个 main 程序

```bash
root@fbc9c0371369:/# sudo /tmp/2/main
2025/12/12 11:38:05 Cthulhu Trigger Observer 启动...
2025/12/12 11:38:05 端口占用失败: listen tcp :23333: bind: address already in use
```

怎么总感觉打的很暴力呢

但是知道是网络了，用 lsof 就可以了

```bash
root@fbc9c0371369:/# lsof -i
COMMAND PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
main     34 root    4u  IPv6    979      0t0  TCP *:23333 (LISTEN)
```

所以监听着干什么呢？我用浏览器 localhost 好像也访问不了的样子

哦，所以要在容器内访问，因为没有做对应的网络映射？

用 linux curl 指令访问本体端口？

?这个镜像怎么没 curl？
那我装一个

发 get 没反应啊？

还是从密文下手吧

```bash
root@fbc9c0371369:/# sudo find / -name bXV4aQ==
/usr/lib/x86_64-linux-gnu/krb5/plugins/secret/bXV4aQ==
root@fbc9c0371369:/# cd /usr/lib/x86_64-linux-gnu/krb5/plugins/secret/
root@fbc9c0371369:/usr/lib/x86_64-linux-gnu/krb5/plugins/secret# ls
 README.txt  'bXV4aQ=='
root@fbc9c0371369:/usr/lib/x86_64-linux-gnu/krb5/plugins/secret# cat README.txt
journal 11月17日：
我疯狂地追寻着深处的不可名状的意志，回过头来就已经再也离不开这个容器了。
我留下了我的日记，作为对后来者的绝望警示。
祂的存在即是污染！！！
若你发现日记的文本开始出现非欧几里得的扭曲或低语的符号，请务必搜寻那些以我颤 抖的手嵌入的 'journal' 标记——它们或许是残存的、未被抹去的逻辑碎片，是真相在黑暗中微弱的余烬。
```

...666 我 cat bXV4aQ==差点没吓死

用 grep 抓一下 journal 吧

```bash
root@fbc9c0371369:/usr/lib/x86_64-linux-gnu/krb5/plugins/secret# cat bXV4aQ\=\= | grep journal
journal 11月20日：不眠的代价。我的手在键盘上自动敲击，像被某种力量驱动。bXV4aQ== 频繁出现，路径 /opt/lib/b/X/V/4/a/Q/=/=/ 仿佛在对我低语。现实的时间线开始裂开，我分不清实验与梦境的界限。
上面的journal都是学长写的吗？他为什么要记录到容器里面？容器内部的这个可执行文件，是试图保护什么，还是……监视谁？
journal 11月19日：仍然无法停止。我开始怀疑，这个文件并非普通程序，而是一种…… 某种意识的容器。bXV4aQ== 的符号似乎能感知我的操作，它的魔力像触手般伸入我的思维。每次尝试停止，都伴随着屏幕异常扭曲的字符闪现。
journal 11月21日：这是最后一次尝试记录。容器、符号、程序，它们像一张网，将我 的意识完全捕获。我感到一种无形的存在在观察我，记录我，甚至干预我的操作。我…… 可能无法再离开。
journal 11月18日：不可能……我无法停止关注那个文件。/opt/lib/b/X/V/4/a/Q/=/=/==Qa4VXb，它像在呼唤我，每次运行它，bXV4aQ== 都会随机在屏幕闪烁。我已经连续三天未眠，代码与现实边界逐渐模糊。
root@fbc9c0371369:/usr/lib/x86_64-linux-gnu/krb5/plugins/secret# cd /opt/lib/b/X/V/4/a/Q/=/=/
root@fbc9c0371369:/opt/lib/b/X/V/4/a/Q/=/=# ls
'==Qa4VXb'
root@fbc9c0371369:/opt/lib/b/X/V/4/a/Q/=/=# ./\=\=Qa4VXb
bXV4aQ==                                             bXV4aQ==
             bXV4bXV4aQ==          bXV4aQbXV4aQ==             bXV4aQ==aQ==
          bXV4aQ==               bXV4aQ==                     bXV4aQ==bXV4aQ==                                       bXV4aQ==
 bXV4aQ==V4aQ==   bXV4aQ==       bXV4aQ==         bXV4aQ==
   bXV4aQ==                bXV4aQ==                                 bXV4aQ==
                                  bXV4aQ==    bXV4aQ=====          bXV4aQ==
 bXV4aQ==                          bXV4aQ==4aQ==  bXbXV4aQ==bbXV4aQ==
          bXV4aQ==                                      bXV4aQ== bbXV4aQ====
bXV4aQ==                                              bXVbXV4aQ==      bXV4aQ==bXV4aQ== bXV4aQ==                  bXbXV4aQ==V4aQ==
                                     bXV4aQ==  bXV4aQ==       bXbXV4aQ==
                bXbXV4aQ==                                  bXV4aQ==
                       bXV4aQbXV4aQ===bXV4aQ==
            bXVbXV4aQ==4aQ==                             bXV4abXV4aQ==
         bXV4aQ==V4aQ==                                           bXV4aQ==
                                             bbXV4aQ==
           bXV4aQ==             bXV4aQ==
                   bXV4abXV4aQ==XV4aQ==    bXV4aQ==            bXV4aQ==V4aQ==
                                 bXV4aQ==bXV4aQ==V4aQ==
            bXV4aQ==bXV4bXV4aQ==
journal 11月22日 ：
不对劲。那文件夹 ‘bXV4aQ==’ 缺失了私钥，导致了逻辑链的断裂。这不该发生。我开始相信，那串编码是某个存在的真名，而我现在被困在了一个没有引渡符文的维度。我必须找到它。
root@fbc9c0371369:/opt/lib/b/X/V/4/a/Q/=/=# ls
'==Qa4VXb'  'bXV4aQ=='
```

线索又断了

并且好像炸烂了，重开一遍镜像

export 查看环境变量发现点东西

`declare -x PROMPT_COMMAND="history -a; history 1 | sed \"s/^ *[0-9]* *//\" | xargs -I{} echo \"{} | \$PWD\" >> /root/.bash_history"`

这很明显是一个指令，是记录历史

...但是无关好像

好。。。我回到 home 了，我才想起来 ls -a 看看有没有隐藏

有个 triggerfs 的文件夹

里面有个???.pem 的密钥，一个由？组成的？。。。

```bash
root@de891520705e:/home# ls -a -l
total 2820
drwxr-xr-x 1 root root    4096 Dec 12 18:03 .
drwxr-xr-x 1 root root    4096 Dec 12 17:04 ..
-rwxr-xr-x 1 root root 2871615 Dec  5 18:57 start
drwxr-xr-x 2 root root    4096 Dec 12 18:14 triggerfs
root@de891520705e:/home# cd triggerfs/
root@de891520705e:/home/triggerfs# ls -l -a
total 12
drwxr-xr-x 2 root root 4096 Dec 12 18:14 .
drwxr-xr-x 1 root root 4096 Dec 12 18:03 ..
-r--r--r-- 1 root root  279 Dec 12 18:03 ？？？.pem
root@de891520705e:/home/triggerfs# cat ？？？.pem

       ？？？？？？？
     ？？？       ？？？
    ？？？         ？？？
    ？？？         ？？？
                ？？？
            ？？？
           ？？？
           ？？？
           ？？？

           ？？？
           ？？？
```

哦哦看过之后变成日记了（

```bash
root@de891520705e:/home/triggerfs# cat 日记.txt
journal 11月15日：

我是一名后端开发工程师，最近在潜心研究 Docker。
虚拟化本是理性的工具，却逐渐显露出一种无法言说的违和感。
那个容器……我居然给它取了个名字——“深处”。
我为什么会给一个容器取名字？

每当我接触“深处”，总觉得有一双无形的眼睛在暗处注视着我。
代码与配置像是被某种不可名状的存在扭曲，
逻辑和现实的界限也开始轻微地、但持续地……偏移。

我无法理解。
为什么同样的镜像，在本地启动与在服务器上运行，会呈现两种完全不同的状态？
那串符号……
那串该死的符号为什么会出现在那里？明明我根本没有写过。

它只是静静地浮在那里，可只要盯着它看上两秒，
我就能听到低沉的嗡鸣，
像从深海几十公里下的黑暗里传来的呼唤，
冰冷、扭曲，贴着我的脊椎一路往上爬。
——记录到此为止。我不能再继续写下去。那串符号又出现了，我得去查 容器日志 了，希望今天24点前能下班。

……这个文笔？
这……这是前辈的日记？
```

并非隐藏，是我干了点啥事就触发了

感觉是打印了一个文件之后就触发了？

对`cat /usr/bin/containerexec/exec.log`后触发 triggerfs 文件
cat 完日记后，会产生一个 secret.pem

这应该就是前面提到的需要的密钥

然后把这个密钥移动到需要的文件夹里

再运行程序

```bash
journal 11月23日 ：
我迷茫的在容器中打转，终于在某个目录找到了秘钥。但在那一刻，我瞥见祂的使者——一个驻守在端口 23333 的、由纯粹意图构成的存在。它在监视我，是阻碍我逃离这容器的卫兵。我必须用它的信息，并使用管道输入 到 ==Qa4VXb，逆转它的控制。我不能再拖延了,或许可以试试 lsof ?。
root@bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=# lsof -i | ./\=\=Qa4VXb
               bXV4aQ==

                                                            bXV4aQ==


bXV4aQ==

                                       bXV4aQ==



         bXV4aQ==
 bXV4aQ==



                                                  bXV4aQ==
      bXV4aQ==    bXV4aQ==                                 bXV4aQ==
                   bXV4aQ==
         bXV4aQ==


                   bXV4aQ==
⚠️ ⚠️ ⚠️ ⚠️ 聆听来自深空的低语。你已被祂观测 ⚠️ ⚠️ ⚠️ ⚠️
journal 11月24日 ：
我失败了。它还活着，而且在反刍我的代码。它如何从编译好的二进制文件中，将我的核心逻辑抽离出来？我分不清哪个是我的操作，哪个是它的干预。现在唯一的办法，是修好我的函数然后带着 '--check' 参数，强行 校准程序，迫使它面对真理。
```

产生了一个 tool.go 文件

其实做的是快排，哎哎，程序填空这一块

弄完了直接==Q... --check 就可以

```bash
root@bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=# ./\=\=Qa4VXb --check
                    bXV4aQ==
                                bXV4aQ== bXV4aQ==V4aQ=bXV4aQ==         bXV4aQ==
                                              bXV4aQ==
                      bXV4aQ=bXV4aQ== bbXV4aQ==bXV4abXV4aQ==Q==
bXV4aQ== bXV4aQ==                                      bXV4aQ==bXV4aQ==
               bXV4aQ==          bXV4aQ==                              bXV4aQ==
bXV4aQ==
    bXV4aQ==                bXV4aQ==                   bXV4aQ==
         bXV4aQ==                                 bXV4aQ==
bXV4aQ==bXV4aQ==                                       bXV4aQ==
          bXV4aQ==                     bXV4aQ==
                                           bXV4aQ== bXV4aQ==       bXV4aQ==
 bXV4aQ==              bXV4aQ==                   bXV4aQ==          bXV4aQ==
                 bXV4bXV4aQ==         bXV4aQ==
      bXV4aQ==                   bXV4aQ==      bXV4bXV4aQ==XV4aQ==
  bXV4aQ==                                bXV4aQ====
                                       bXV4aQ==   bXV4aQ== bXV4bXV4aQ==aQ====
            bXV4aQ==V4aQ==           bXV4aQ==                      bXV4aQ==
   bXV4aQ==    bXV4aQ==                                               bXV4aQ==
                    bXV4aQ==   bXV4aQ=bXV4aQ==            bXV4aQ==
                           bXV4aQ==   bXV4aQbXV4aQ==   bXV4aQ==
       bXV4aQ====aQ=bXV4aQ==       bXV4aQ==
bXVbXV4aQ== bXV4aQ==                  bXV4aQ==           bXV4aQ==
           bXV4aQ===bXV4aQ==     bXV4aQ=bXV4aQ===     bXV4aQ==
journal 11月25日 ：
环境问题？这不可能是一个环境问题。我被困在一个虚假的现实里，参数皆为幻影。我懂了，我必须把 ‘bXV4aQ==’ 强行设置为环境变量，让 bXV4aQ== 等于 bXV4aQ== ，用魔法打败魔法。
```

但是我说，linux 没法弄成这样的环境变量吧？所以说强行？但是怎么个强行法？

哦。。。用~/.bashrc 好像是...

反正最后是弄成了

journal 11 月 26 日 ：
老板，同事，需求，怎么是他们？他们是怎么找到我的，他们也进入到容器了？我在哪里？现实还是容器

带上得到的文件,是时候该离开容器了

```bash
root@bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=# ls
'==Qa4VXb'  'bXV4aQ=='   deep-inside-container-C.tar   docker-compose.yaml   secret.pem   tool.go
```

得到了 deep-inside-container-C.tar docker-compose.yaml 喵
docker-compose.yaml 应该是用来部署多容器的配置文件

```bash
root@bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=# cat docker-compose.yaml
version: '3.8'
# TODO 让服务A运行在8080端口上，A和B需要在同一个网络，访问A以寻求你想要的
services:
  service_a:
    image: registry.cn-shenzhen.aliyuncs.com/muxi/deep-inside-the-container:A
    container_name: service_a_container
    ports:
      - "8080:8080"
    networks:
      - internal_network

  service_b:
    image: registry.cn-shenzhen.aliyuncs.com/muxi/deep-inside-the-container:B
    container_name: service_b_container
    networks:
      - internal_network
networks:
  internal_network:
    driver: bridgeroot@bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=#
```

用 docker cp 就可以了

```pwsh
Administrator in MuXi\go\week7 main*​​​ ⇡
❯ docker cp bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=/deep-inside-container-C.tar ./deep-inside-container-C.tar
Successfully copied 3.88MB to E:\1\2\MuXi\go\week7\deep-inside-container-C.tar

Administrator in MuXi\go\week7 main ⇡
❯ docker cp bfbf11c01ade:/opt/lib/b/X/V/4/a/Q/=/=/docker-compose.yaml ./docker-compose.yaml
Successfully copied 2.56kB to E:\1\2\MuXi\go\week7\docker-compose.yaml
```

tar 包可以用 tar 解压缩`tar -xvf deep-inside-container-C.tar`

导入镜像的话，用 docker import

```pwsh
Administrator in MuXi\go\week7 main*​​​ ⇡
❯ docker import .\deep-inside-container-C.tar deep-inside-container-C:lastest
invalid reference format: repository name (library/deep-inside-container-C) must be lowercase
```

...6,docker 只能用小写字母做为镜像名

然后启动看看

```pwsh
Administrator in MuXi\go\week7 main*​​​ ⇡7s
❯ docker run ec5fd3e613c8
docker: Error response from daemon: no command specified

Run 'docker run --help' for more information
```

这个好像是说要靠 docker-compose 来搞才行

> 当你使用 docker import 导入镜像时，原始容器的启动命令(CMD/ENTRYPOINT)会丢失。这导致 Docker 不知道容器启动时该运行什么程序。

先删掉我看看吧

emmm 好像是先 docker-compose.yaml 配置好这个多容器，然后访问它所说的 8080 获取信息
在这个 yaml 所在目录直接`docker-compose up -d`即可

```返回
FROM deep-inside-the-container:C
# TODO 设置环境变量MUXI=MUXI,将/app/secret/muxi.txt移动到/muxi/muxi.txt
# journal 11月27日 ：
# 要离开了，如果有人看到这里一定会以为我是个疯子或者是瞥见了古神而丧失了理智，其实我也只是一个可怜虫罢了，完成这个dockerfile，把生成的镜像导出成tar，使用form发送到这个地址，参数是image：
# https://deep-inside-the-container.muxixyz.com/finish
# 在那里，我会告诉你一切。
CMD ["sh"]
```

所以要`docker exec -it cc5eaa5f4a99 sh`进入容器 a 做操作
CMD ["sh"] 告诉我们了，用 sh 作为命令进入

```sh
/etc # cat os-release
NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.23.0
PRETTY_NAME="Alpine Linux v3.23"
HOME_URL="https://alpinelinux.org/"
BUG_REPORT_URL="https://gitlab.alpinelinux.org/alpine/aports/-/issues"
```

...但是没这个文件啊

哦哦，是给的这个 tar 导入的 CMD 设置为 sh 是吧

...不对，这个好像就是 dockerfile 的内容

要用 dockerfile + docker build 重新构建一次

```Dockerfile
FROM deep-inside-the-container:C
ENV MUXI=MUXI
RUN mkdir -p /muxi && \
    mv /app/secret/muxi.txt /muxi/muxi.txt
CMD ["sh"]
```

命令如上,`mkdir -p /muxi`是创建目录（不存在就创建，存在也不失败）,-p 好像是说前面就算不存在也会创建
`&& \`是前面一行执行成功就执行后面的命令，也就是移动文件
最后是启动命令

hadolint 是一个可以检测 Dockerfile 语法是否正确的命令行软件（vsc 其实就可以了，但是提到了就玩一下）

然后 `docker build -t` 就可以了

```pwsh

```

...折腾了半天，感觉 tar 包跟坏了一样（（（（，镜像一直有问题
也感觉不会有我做的那么复杂的操作（自己解压层自己导入）
之后再说吧

...md
我从 tmp 那里烤的 sha256 文件大小就和打到最后得到的不一样
估计就是坏了
...md 好像是一样的
那为什么
我不明白~~~

...好，标准结局，今天弄又好了

```bash
Administrator in MuXi\go\week7 main ⇡
❯ docker build -t deep-inside-container:temp --no-cache .
[+] Building 1.0s (6/6) FINISHED                                                                                                                                             docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                                                                                                         0.0s
 => => transferring dockerfile: 648B                                                                                                                                                         0.0s
 => [internal] load metadata for docker.io/library/deep-inside-the-container:C                                                                                                               0.2s
 => [internal] load .dockerignore                                                                                                                                                            0.0s
 => => transferring context: 2B                                                                                                                                                              0.0s
 => [1/2] FROM docker.io/library/deep-inside-the-container:C@sha256:21939759f9abd80a5173eb0087149ac502778855a3ca938aaa581907642d64af                                                         0.1s
 => => resolve docker.io/library/deep-inside-the-container:C@sha256:21939759f9abd80a5173eb0087149ac502778855a3ca938aaa581907642d64af                                                         0.0s
 => [2/2] RUN mkdir -p /muxi &&     mv /app/secret/muxi.txt /muxi/muxi.txt                                                                                                                   0.2s
 => exporting to image                                                                                                                                                                       0.3s
 => => exporting layers                                                                                                                                                                      0.1s
 => => exporting manifest sha256:02194d977c12b44fec1f96a9100da3132f03c57c7f0bc8b76749ce5206ade1ef                                                                                            0.0s
 => => exporting config sha256:89270977eab9bb82f4f296eb2d668537b479689df35249529e2cab7af7d63065                                                                                              0.0s
 => => exporting attestation manifest sha256:9d365c2194ef6edc992f3c0a3abfd538960da3111db7ade186939b6849228631                                                                                0.0s
 => => exporting manifest list sha256:b56f6bf603b35ae4db8e7e056a67c6668f730e18fcf2aff0d10811c09268f1aa                                                                                       0.0s
 => => naming to docker.io/library/deep-inside-container:temp                                                   0.0s
 => => unpacking to docker.io/library/deep-inside-container:temp                                                                                                                             0.0s

Administrator in MuXi\go\week7 main ⇡2s
❯ docker save -o deep-inside-the-container-final.tar b56f6bf603b3
Administrator in MuXi\go\week7 main ⇡
❯ curl -X POST https://deep-inside-the-container.muxixyz.com/finish -F "image=@deep-inside-the-container-final.tar"
恭喜你！镜像已通过校验。

终于，你走到了这条路的尽头。
但真相，比你想象的要简单，也更令人心碎。

我不是一个容器幽灵，也不是被拉普拉斯妖附身的程序员。
我只是一个被无止境的运维任务和老板的PUA压垮的灵魂。那些克苏鲁小说的幻觉，不过是我疲惫的神经在无眠的长夜里，为自己构建的逃离现实的梦境。

所有的怪异事件、那些看似无法解释的 Bug 和诡异的后门，并非源于什么黑暗的深渊。
它们真正的源头，不过是我部署时，一个不经意的疏忽：错误的开放了Docker的挂载权限，导致有人乘虚而入，植入了恶意的程序。

这是一个关于疏忽、疲惫与现实的故事。
容器是界限，代码是规则。而你，用逻辑和智慧，重新划清了界限，修补了规则....

游戏到此结束。
你已窥见隐藏在容器深处的真相，完成了使命。

请记住，在无尽的代码和版本迭代中，我们是人，不是机器。
不要让压力扭曲了你眼中的现实。

无论前方有多少个 'latest' 标签等待更新，请先关掉屏幕，去好好休息。
珍爱生命，远离加班。

—— Deep Inside The Container 项目组(虽然其实只有我一个人)
```
