# Bendan Bot

一个 Telegram 的 [@bendan_bot](https://t.me/bendan_bot)

## 功能

目前包含以下指令或功能：

- 对消息回复 `/*`。如 `/吃`，Bot 发送一个 `a 吃了 b ！`
- 对消息回复 `/* *`。如 `/吃 豆腐`，Bot 发送一个 `a 吃 b 豆腐！`
- 直接发送或回复 `/me *`。如 `/me 喝醉了`，Bot 发送一个 `a 喝醉了！`
- 对消息回复 `//pin`，Bot 将消息置顶（若有权限），最多保留 10 条由 Bot 置顶的消息
- 直接发送或回复 `//whoami`，Bot 发送指令调用者 ID、当前群组 ID
- 直接发送或回复 `/没关系`，生存鼓励师猫猫发送一条“没关系……”句子
- 识别具有跟踪参数的 URL，并自动净化它们
- 直接发送或回复 `？`，Bot 回复一个 `？`
- 直接发送或回复 `看看…`，Bot 半随机地回复一个内容
- 直接发送或回复 `是…吗`、`有…吗`、`是…还是…`、`有…还是…`，Bot 半随机地回复一个内容
- 直接发送或回复 `是不是…`、`有没有…`、`会不会…`、`能不能…`、`…行不行`、`这么有…`，Bot 半随机地回复一个内容
- 直接发送或回复 `//lang code` 执行代码（可多行），目前支持 Go、JavaScript，如 `//js Math.sin(45)`
- 转发 Telegram channel 消息到 Twitter、Mastodon

## 部署到 Vercel

你可以 Fork 本仓库，然后一键部署到 Vercel。需要在 Vercel ENV 中配置如下变量：

- `BOT_TOKEN`：Telegram Bot 的 Token
- `REFRESH_KEY`：用于第一次触发 Webhook 配置的密钥
- `DB_URI`：（可选）MongoDB 的连接 URI。用于置顶、转发功能，不需要这些功能则不配置
- `DB_NAME`：（可选）MongoDB 的数据库名称。同上，可不配置
- `FORWARD_CONFIG`：（可选）转发功能配置，JSON 字符串，具体字段见 [部署到 Server](#部署到-Server)。无需该功能则不配置

设置完成后，访问 `https://your-domain.vercel.app/?key=<refresh_key>` 对 Webhook 初始化即可。

## 部署到 Server

您需要做的是，在根目录创建 `.config` 文件，并做如下配置：

```json
{
  "bot_token": "...",
  "refresh_key": "secret",
  "db_uri": "mongodb+srv://...",
  "db_name": "bendan",
  "forward_config": {
    "group": {
      "id": "channel 关联群组 ID，可通过 //whoami 指令查询",
      "owner": "群组所有者 ID，用于后续转发交互鉴权"
    },
    "twitter": {
      "consumer_key": "Twitter 应用ID，在 developer.twitter.com 创建应用",
      "consumer_secret": "Twitter 应用密钥，同上",
      "user_token": "对应 Twitter 用户的 Access Token，单用户可在 developer.twitter.com 生成，或多用户基于应用的 OAuth1 授权",
      "user_secret": "user_token 对应的 secret"
    },
    "mastodon": {
      "endpoint": "如 https://mastodon.social",
      "token": "Preferences - Development 中创建应用得到"
    },
    "allowed_tags": ["随原文一起转发的 tag"]
  }
}
```

以上参数同 [部署到 Vercel](#部署到-vercel)，配置完成后直接运行主程序即可。

## 隐私

受 Telegram Bot API 限制，无法获取所有已置顶消息，以用于“仅保留 10 条置顶”功能。

因此需要将由 Bot 置顶的消息存储下来，但也仅存储了 `(置顶消息ID, 群组ID, 置顶时间)` 用于取消老的消息置顶，该 Bot **_不会存储任何与消息内容有关的信息_**。

## License

Bendan Bot is licensed under the [MIT License](https://github.com/sxyazi/bendan/blob/master/LICENSE).
