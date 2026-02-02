# Mermaid 示例代码

> 这是一种，我觉得类似 md 的文本图示语言，比较优雅吧，所以大概演示一下，以下 Qwen 生成的一些示例代码，能学个七七八八的语法

Mermaid 有以下几种常见的图标示例:

1. 流程图 (Flowchart)
2. 时序图/序列图 (Sequence Diagram)
3. 甘特图 (Gantt Chart)
4. 状态图 (State Diagram)
5. 饼图 (Pie Chart)
6. 类图 (Class Diagram)

```mermaid
flowchart TD
    A[开始] --> B{用户登录}
    B -->|成功| C[进入主页]
    B -->|失败| D[显示错误]
    C --> E[选择功能]
    E --> F[数据分析]
    E --> G[报告生成]
    F --> H[导出结果]
    G --> H
    H --> I[结束]

    classDef success fill:#9f9,stroke:#090;
    classDef error fill:#f99,stroke:#900;
    class C,F,G,H success;
    class D error;
```

```mermaid
sequenceDiagram
    participant A as 用户
    participant B as 前端
    participant C as API服务器
    participant D as 数据库

    A->>B: 输入查询条件
    B->>C: 发送API请求
    C->>D: 查询数据
    D-->>C: 返回结果
    C-->>B: 格式化响应
    B-->>A: 显示结果

    Note over A,B: 用户交互
    Note over C,D: 服务器处理
```

```mermaid
gantt
    title 项目时间表
    dateFormat  YYYY-MM-DD
    section 需求分析
    需求收集     :a1, 2023-10-01, 7d
    需求评审     :a2, after a1, 3d

    section 开发阶段
    前端开发     :b1, 2023-10-12, 10d
    后端开发     :b2, 2023-10-12, 12d

    section 测试
    单元测试     :c1, after b1, 5d
    集成测试     :c2, after c1, 5d

    section 部署
    上线准备     :d1, after c2, 3d
    正式发布     :d2, after d1, 1d

    critical 需求评审到开发 :crit, a2, 2023-10-11, 1d
```

```mermaid
pie
	title 项目时间分配
	"前端开发" : 35
	"后端开发" : 40
	"测试" : 15
	"文档编写" : 10
```

```mermaid
stateDiagram-v2
	[*] --> 准备中
	准备中 --> 运行中 : 开始执行
	运行中 --> 暂停 : 暂停按钮
	运行中 --> 结束 : 完成执行
	暂停  --> 运行中 : 继续执行
	暂停 --> 结束 : 取消执行
	结束 --> [*]
```

```mermaid
classDiagram
	class Animal {
		+String name
		+int age
		+void eat()
		+void sleep()
	}

	class Dog {
		+String breed
		+void bark()
	}

	class Cat {
		+String color
		+void meow()
	}

	Animal

	class Pet {
		<>
		+void play()
	}

	Dog ..|> Pet
	Cat ..|> Pet
```
