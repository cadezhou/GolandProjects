package main

import (
	"log"
	"net/http"

	"qqbotmessage/plugin"
)

func main() {
	// 1. 加载配置
	cfg, err := plugin.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化 Bot（鉴权 + 注册事件处理器）
	bot := plugin.NewBotPlugin(cfg.QQ.AppID, cfg.QQ.Secret)
	if err := bot.Start(); err != nil {
		log.Fatalf("Bot 启动失败: %v", err)
	}

	// 3. 注册 Webhook 回调路由 当访问/qqbot路径的时候 交给webhookhandler来处理
	http.HandleFunc(cfg.Server.Path, func(w http.ResponseWriter, r *http.Request) {
		plugin.WebhookHandler(bot.Creds())(w, r)
	})

	// 4. 启动 HTTP 服务（阻塞，直到进程退出）
	log.Printf("服务启动，监听 %s", cfg.Server.Addr())
	if err := http.ListenAndServe(cfg.Server.Addr(), nil); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
