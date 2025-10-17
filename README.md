# Cpay  - 区块链数字货币支付收款系统 🔒

一个基于 Go 语言开发的现代化加密货币支付网关系统，支持多种数字货币支付，提供完整的订单管理和自动化支付验证功能。


## 📌 项目亮点
- ✅ 支持多链路支付（波场Tron/Optimism/Base/Matic/Arbitrum/Avax）
- ✅ 支持多币种收款（USDT/USDC/其他任何ERC20或TRC20合约）
- ✅ API 接口 【包含客户端api和管理api】
- ✅ 安全可靠
- ✅ 高性能
- ✅ 补单功能

## 注意事项
 - 首次使用需要将etc/config.yaml.example重命名为config.yaml

## 📋 系统要求

- Go 1.24.4 或更高版本【二开推荐】
- SQLite 数据库 【可替换为MySQL等】
- Redis 【用于订单缓存和支付验证】 默认用内存

## 💻 命令行操作
```bash
1 - 运行程序
2 - 设置用户名
3 - 设置密码
4 - 设置Secret
5 - 查看默认信息
```

## 🏗️ 项目结构
```
cpay/
├── etc/                 # 配置文件目录
├── main.go                 # 程序入口
├── internal/               
│   ├── admin/             # 后台管理员模块
│   ├── api/               # API 模块
│   ├── blockchain/        # 区块链封装模块
│   ├── cron/               # 定时任务模块
│   ├── consumer/             # 消息队列消费者模块
│   ├── logger/             # 日志模块
│   ├── middleware/             # 中间件模块
│   ├── model/             # 数据模型模块
│       ├── cache/             # 缓存模块
│       ├── dao/             # 数据库操作模块
│       └── mdb/             # 数据库模型模块
│   ├── router/             # 路由模块
│   ├── runner/             # 常驻运行工人
│   ├── service/             # 服务模块
│   └── function.go        # 业务逻辑函数
```
## 📊 监控和日志

- **结构化日志**: 使用 Zap 日志库，支持日志轮转,按照链路分别生成日志文件

## 🔄 定时任务

系统包含以下定时任务：

- **支付检查**: 每 5 秒检查一次未支付订单的区块链状态
- **支付归集**: 定时将子钱包的金额归集到主钱包上,

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。


## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [GORM](https://gorm.io/) - ORM 库
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Cron](https://github.com/robfig/cron) - 定时任务库
- [Redis](https://redis.io/) - 缓存数据库
- [SQLite](https://www.sqlite.org/index.html) - 数据库


---

## 任务清单[TODO.md](TODO.md)


## 签名机制[SIGNATURE.md](SIGNATURE.md)


## 打赏
如果该项目对您有所帮助，希望可以请我喝一杯咖啡☕️
````
Usdt(trc20)打赏地址: TCB6JRqnkWdpbkRq4gPtWAf1fEu6B4Mgxn
````


**注意**: 本项目仅供学习和研究使用，请确保在合法合规的前提下使用本系统。
