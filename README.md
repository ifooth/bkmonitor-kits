# bkmonitor-kits

> 蓝鲸监控 Golang 工具包

## 模块

### logger

日志库，封装了 go.uber.org/zap 和 lumberjack.v2 支持日志切割。

```golang
package main

import "github.com/TencentBlueKing/bkmonitor-kits/logger"

// 初始化日志库配置选项
func InitLogger() {
	logger.SetOptions(logger.Options{
		Filename:   "/data/log/myproject/applog",
		MaxSize:    1000, // 1GB
		MaxAge:     3,    // 3 days
		MaxBackups: 3,    // 3 backups
	})
}

func main() {
	// 生成环境的话可以试着自定义的日志配置 默认的输出流是标准输出
	InitLogger()

	logger.Info("This is the info level message.")
	logger.Warnf("This is the warn level message. %s", "oop!")
	logger.Error("Something error here.")
}
```

### host

监控主机标识。

### register

consul 域名注册。

### validator

监控数据上报校验。

## Contributing

我们诚挚地邀请你参与共建蓝鲸开源社区，通过提 bug、提特性需求以及贡献代码等方式，一起让蓝鲸开源社区变得更好。

![bkmonitor-kits](https://user-images.githubusercontent.com/19553554/126454082-d21b22f9-6df9-487f-82c1-a9dcd054f29a.png)


## License

基于 MIT 协议，详细请参考 [LICENSE](./LICENSE)
