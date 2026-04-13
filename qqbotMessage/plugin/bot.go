package plugin

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"golang.org/x/oauth2"
)

type BotPlugin struct {
	appID     string
	appSecret string
	creds     *token.QQBotCredentials
	tokenSrc  oauth2.TokenSource
}

func NewBotPlugin(appID, appSecret string) *BotPlugin {
	return &BotPlugin{
		appID:     appID,
		appSecret: appSecret,
	}
}

func (p *BotPlugin) Start() error {
	// 1. 创建凭证
	p.creds = &token.QQBotCredentials{
		AppID:     p.appID,
		AppSecret: p.appSecret,
	}

	// 2. 创建 token source
	p.tokenSrc = token.NewQQBotTokenSource(p.creds)

	// 3. 启动 token 自动刷新
	if err := token.StartRefreshAccessToken(context.Background(), p.tokenSrc); err != nil {
		return fmt.Errorf("start refresh token failed: %w", err)
	}

	// 4. 初始化 OpenAPI
	api = botgo.NewOpenAPI(p.appID, p.tokenSrc).
		WithTimeout(5 * time.Second).
		SetDebug(true)

	// 5. 注册事件处理
	_ = event.RegisterHandlers(
		ATMessageHandler(),
		C2CMessageHandler(),
	)

	log.Printf("Bot 初始化完成")
	return nil
}

func (p *BotPlugin) Creds() *token.QQBotCredentials {
	return p.creds
}

// 全局 API 变量
var api openapi.OpenAPI

func GetAPI() openapi.OpenAPI {
	return api
}

// ATMessageHandler 处理 @ 机器人消息
func ATMessageHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[收到消息] GuildID: %s, ChannelID: %s, Content: %s",
			data.GuildID, data.ChannelID, data.Content)

		msg := &dto.MessageToCreate{
			Content:   "收到消息: " + data.Content,
			Timestamp: time.Now().UnixMilli(),
			MsgID:     data.ID,
		}

		if _, err := api.PostMessage(context.Background(), data.ChannelID, msg); err != nil {
			log.Printf("发送消息失败: %v", err)
			return err
		}

		return nil
	}
}

// C2CMessageHandler 处理私信消息
func C2CMessageHandler() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		log.Printf("[收到私信] UserID: %s, Content: %s", data.Author.ID, data.Content)

		msg := &dto.MessageToCreate{
			Content:   "收到私信: " + data.Content,
			Timestamp: time.Now().UnixMilli(),
		}

		if _, err := api.PostC2CMessage(context.Background(), data.Author.ID, msg); err != nil {
			log.Printf("发送私信失败: %v", err)
			return err
		}

		return nil
	}
}
