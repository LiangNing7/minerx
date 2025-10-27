## 客户端规范

- 每个客户端至少包含以下4个文件：
  - client.go: Interface 定义，New, Functions
  - doc.go: 包 doc.go 文件
  - helper.go: 类似于 util.go 文件
  - model.go: 模型定义文件
  - options.go: 命令行参数定义，需要满足 IOptions 接口
- client.go: 
  - GetClient: 获取 global impl 实例
  - 使用sync.Once，确保只实例化一次
