## xy-dianping 项目
基于 Go + tRPC + Fx + Redis + Gorm + MySQL + RabbitMQ 技术的仿大众点评项目

参考 Java 版本实现：https://github.com/viego1999/xy-dianping

---
项目组织结构目录如下：
```
project-root/  
    ├── cmd/                       # 存放应用入口文件  
    │   └── main.go               # 主程序入口  
    ├── internal/                 # 内部包，通常不被外部项目引用  
    │   ├── config/               # 配置相关  
    │   │   └── config.go         # 配置文件读取和解析  
    │   ├── middleware/           # 中间件  
    │   │   └── logging.go        # 日志记录中间件  
    │   ├── service/              # 服务层  
    │   │   ├── user_service.go   # 用户服务实现  
    │   │   └── ...               # 其他服务实现  
    │   ├── repo/                 # 存储库层  
    │   │   ├── user_repository.go # 用户存储库实现  
    │   │   └── ...               # 其他存储库实现  
    │   └── ...                   # 其他内部包  
    ├── pkg/                      # 公用包，可被其他项目引用  
    │   ├── util/                 # 工具函数  
    │   │   └── string_utils.go  # 字符串处理工具  
    │   └── ...                   # 其他公用包  
    ├── models/                   # 数据模型定义  
    │   ├── user.go               # 用户模型  
    │   └── ...                   # 其他模型  
    ├── api/                      # API定义和路由  
    │   ├── routes.go             # 路由定义  
    │   ├── v1/                   # API版本目录  
    │   │   ├── user/             # 用户相关API  
    │   │   │   ├── user.go       # 用户API处理器  
    │   │   │   └── ...           # 其他用户API处理器  
    │   │   └── ...               # 其他API版本目录  
    │   └── ...                   # 其他API定义  
    ├── docs/                     # 文档  
    │   ├── api_docs.md           # API文档  
    │   └── ...                   # 其他文档  
    ├── go.mod                    # Go模块依赖文件  
    └── go.sum                    # Go模块依赖校验文件
```