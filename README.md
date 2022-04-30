# Bendan Bot

一个 Telegram 的 [@bendan_bot](https://t.me/bendan_bot)

## 功能

目前包含以下指令：

- 对消息回复 `/*`。如 `/吃`，Bot 发送一个 `a 吃了 b ！`
- 对消息回复 `/* *`。如 `/吃 豆腐`，Bot 发送一个 `a 吃 b 豆腐！`
- 直接发送或回复 `/me *`。如 `/me 喝醉了`，Bot 发送一个 `a 喝醉了！`
- 对消息回复 `//pin`，Bot 将消息置顶（若有权限），最多保留 10 条由 Bot 置顶的消息

## 部署到 Vercel

你可以 Fork 本仓库，然后一键部署到 Vercel。需要在 Vercel ENV 中配置如下变量：

- `BOT_TOKEN`：Telegram Bot 的 Token
- `DB_URI`：MongoDB 的连接 URI。用于置顶功能，如果不需要此功能，可以不配置
- `DB_NAME`：MongoDB 的数据库名称。同上，可以不配置
- `REFRESH_KEY`：用于第一次触发 Webhook 配置的密钥

设置完成后，访问 `https://your-domain.vercel.app/?key=<refresh_key>` 对 Webhook 初始化即可。

## 部署到 Server

您需要做的是，在根目录创建 `.config` 文件，并做如下配置：

```json
{
  "bot_token": "...",
  "db_uri": "mongodb+srv://...",
  "db_name": "bendan",
  "refresh_key": "secret"
}
```

以上参数同 [部署到 Vercel](#部署到-vercel)，配置完成后直接运行主程序即可。

## 隐私

受 Telegram Bot API 限制，无法获取所有已置顶消息，以用于“仅保留10条置顶”功能。

因此需要将由 Bot 置顶的消息存储下来，但也仅存储了 `(置顶消息ID, 群组ID, 置顶时间)` 用于取消老的消息置顶，该 Bot ***不会存储任何与消息内容有关的信息***。

## License

Bendan Bot is licensed under the [MIT License](https://github.com/sxyazi/bendan/blob/master/LICENSE).
