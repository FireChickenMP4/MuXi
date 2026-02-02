整整 swagger
自己完全不会弄，所以干脆直接找了个例子得了
具体流程就是
先按照其语法写对应注释 @xxxx 什么的
然后`swag fmt`格式化文件
`swag init`就会生成对应的 docs 文件夹
然后记得引入生成的 docs 文件 `_ "swagger/docs"`
最后用 gin 添加 swagger 的路由`r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`

> 用的 gin-swagger
