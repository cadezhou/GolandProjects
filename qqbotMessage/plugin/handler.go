package plugin

import (
	"net/http"

	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/token"
)

// WebhookHandler 返回符合 tRPC HTTP handler 签名的处理函数
func WebhookHandler(creds *token.QQBotCredentials) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		webhook.HTTPHandler(w, r, creds)
		return nil
	}
}
