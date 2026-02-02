# Gorm

> 转自https://zhuanlan.zhihu.com/p/662015722
> 别问我为什么不看官方文档，问就是半栏看着难受
> 而且官方示例复杂，不是很好上手，我太菜惹
> 有关 gorm 的更多资讯:
>
> - 开源地址: https://github.com/go-gorm/gorm
>
> * 中文教程: gorm.io/zh_CN/docs/

## 0 前言

> gorm 是 golang 中最流行的 orm 框架，为 go 语言使用者提供了简便且丰富的数据库操作 api

## 1 数据库

### 1.1 数据库

> 本章中，我们重点向大家介绍如何通过 gorm 创建 mysql db 实例以及完成 db 配置：

- 设置好连接 mysql 的 dsn（data source name）
- 通过 gorm.Config 完成 db 有关的自定义配置
- 通过 gorm.Open 方法完成 db 实例的创建

> 对应 mysql.go

与 database/sql 中原生的 sql.DB 实例不同，在创建 gorm.DB 实例时，默认情况下会向数据库服务端发起一次连接，以保证 dsn 的正确性.

另外想提的一个点是，在 gorm 体系之下，这个 DB 对象是绝对的核心，基本所有操作都是围绕着这个 DB 实例展开的，后续大家也会看到大量通过使用 DB 进行链式调用的代码风格，形如：

`db.Where(...).Order(...).WithContext(...).Find(...)`

## 1.2 配置

在创建 gorm.DB 实例时，可以通过 gorm.Config 进行自定义配置，其中各配置项含义如下：

```go
type Config struct {
    // gorm 中，针对单笔增、删、改操作默认会启用事务. 可以通过将该参数设置为 true，禁用此机制
    SkipDefaultTransaction bool
    // 表、列的命名策略
    NamingStrategy schema.Namer
    // 自定义日志模块
    Logger logger.Interface
    // 自定义获取当前时间的方法
    NowFunc func() time.Time
    // 是否启用 prepare sql 模板缓存模式
    PrepareStmt bool
    // 在 gorm 创建 db 实例时，会创建 conn 并通过 ping 方法确认 dsn 的正确性. 倘若设置此参数，则会禁用 db 初始化时的 ping 操作
    DisableAutomaticPing bool
    // 不启用迁移过程中的外联键限制
    DisableForeignKeyConstraintWhenMigrating bool
    // 是否禁用嵌套事务
    DisableNestedTransaction bool
    // 是否允许全局更新操作. 即未使用 where 条件的情况下，对整张表的字段进行更新
    AllowGlobalUpdate bool
    // 执行 sql 查询时使用全量字段
    QueryFields bool
    // 批量创建时，每个批次的数据量大小
    CreateBatchSize int
    // 条件创建器
    ClauseBuilders map[string]clause.ClauseBuilder
    // 数据库连接池
    ConnPool ConnPool
    // 数据库连接器
    Dialector
    // 插件集合
    Plugins map[string]Plugin
    // 回调钩子
    callbacks  *callbacks
    // 全局缓存数据，如 stmt、schema 等内容
    cacheStore *sync.Map
}
```

## 2 模型

### 2.1 gorm.Model

在定义持久化模型 PO(persist object) 时，推荐组合使用 gorm.Model 中预定义的几个通用字段，包括主键、增删改时间等：

```go
type PO struct {
    gorm.Model
}
package gorm
type Model struct {
    // 主键 id
    ID        uint `gorm:"primarykey"`
    // 创建时间
    CreatedAt time.Time
    // 更新时间
    UpdatedAt time.Time
    // 删除时间
    DeletedAt DeletedAt `gorm:"index"`
}
```

值得一提的是，在 gorm 体系中，一个 po 模型只要启用了 deletedAt 字段，则默认会开启软删除机制：在执行删除操作时，不会立刻物理删除数据，而是仅仅将 po 的 deletedAt 字段置为非空.

这里暂且点到为止，软删除的细节本文第 4 章中再作详细展开.

### 2.2 标签

下面我们介绍一下 po 模型中的常用标签：

```go
type PO struct{
   // 组合使用 gorm Model，引用 id、createdAt、updatedAt、deletedAt 等字段
   gorm.Model
  // 列名为 name；列类型字符串；使用该列作为唯一索引
   Name string `gorm:"column:name;type:varchar(15);unique_index"`
   // 该列默认值为 18
   Age int `gorm:"default:18"`
   // 该列值不为空
   Email string `gorm:"not null"`
   // 该列的数值逐行递增
   Num int `gorm:"auto_increment"`
}
```

几类常用的标签及对应的用途展示如下表：
标签|作用
-|-
primarykey|主键
unique_index|唯一键
index|键
auto_increment|自增列
column|列名
type|列类型
default|默认值
not null|非空

### 2.3 零值

在使用 po 模型时，可能会存在一个与零值有关的问题.

在 golang 中一些基础类型都存在对应的零值，即便用户未显式给字段赋值，字段在初始化时也会首先赋上零值. 比如 bool 类型的零值为 false；string 类型为 ""，int 类型为 0.

这样就会导致，在我们执行创建、更新等操作时，倘若 po 模型中存在零值字段，此时 gorm 无法区分到底是用户显式声明的零值，还是未显式声明而被编译器默认赋予的零值. 在无法区分的情况下，gorm 会统一按照后者，采取忽略处理的方式.

倘若此时我们想要明确是显式将字段设置为零值的，对应可以采取以下两种处理方式：

- 使用指针类型：

我们将 age 字段类型设定为 \*int，只要指针非空，就代表使用方进行了显式赋值.

```go
type PO struct{
   gorm.Model
   Age *int `gorm:"column:age"` // 默认值为 18
}
```

- 使用 sql.Nullxx 类型：

我们将 age 字段类型设定为 sql.NullInt64，只要 Valid 标识为 true，就代表使用方进行了显式赋值.

```go
type PO struct{
   gorm.Model
   Age sql.NullInt64 `gorm:"column:age"` // 默认值为 18
}
type NullInt64 struct {
    Int64 int64
    Valid bool // Valid is true if Int64 is not NULL
}
```
