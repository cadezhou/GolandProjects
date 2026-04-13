# QQ Bot Projects

基于腾讯 [QQ 机器人 SDK (botgo)](https://github.com/tencent-connect/botgo) 的两个项目。

## 项目结构

```
.
├── qqbotMessage/     # QQ 机器人消息插件
└── weatherPush/      # 定时天气推送服务
```

## qqbotMessage

QQ 机器人核心插件，提供消息处理能力。

**功能：**
- 处理 @ 机器人消息
- 处理私信消息
- Webhook 回调支持

**启动：**
```bash
cd qqbotMessage
go run main.go
```

## weatherPush

定时天气推送服务，基于 qqbotMessage 实现。

**功能：**
- 定时获取天气信息（使用 wttr.in）
- 通过 QQ 机器人推送天气给指定用户

**配置 (config.yaml)：**
```yaml
qq:
  appid: "你的AppID"
  secret: "你的Secret"

push:
  user_id: "目标用户ID"
  city: "城市名"
  hour: 8        # 推送小时
  minute: 0      # 推送分钟
```

**启动：**
```bash
cd weatherPush
go run main.go
```

## 依赖

- [botgo](https://github.com/tencent-connect/botgo) - QQ 机器人 SDK
- [trpc-go](https://github.com/trpc-group/trpc-go) - TRPC 框架
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) - YAML 配置文件解析
