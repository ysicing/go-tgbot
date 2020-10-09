# simple bot

> 一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等

## 使用

### 配置文件

```yaml
tgbot: xxx # bot token
mode:
  debug: false
botapi: "https://api.telegram.org/bot%s/%s"
tgchan: "@chanid" # 频道name
tguser: userid # 用户id
```

### 具体使用

```
一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等

Usage:
  sb [flags]

Flags:
      --c string      msg信息或者文件路径 (default "simple bot")
      --chan          msg发送chan
  -h, --help          help for sb
      --type string   msg或者file，默认msg (default "msg")

# 发送msg
sb --c "sb" # 发送到个人
sb --chan --c "sb" # 发送到频道
# 发送file
sb --c ./sb.go --type file
```