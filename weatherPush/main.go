package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"gopkg.in/yaml.v3"
)

// ============ 配置相关 ============

type QQConfig struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

type PushConfig struct {
	UserID string `yaml:"user_id"`
	City   string `yaml:"city"`
	Hour   int    `yaml:"hour"`
	Minute int    `yaml:"minute"`
}

type Config struct {
	QQ   QQConfig   `yaml:"qq"`
	Push PushConfig `yaml:"push"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	return cfg, nil
}

// ============ 机器人鉴权 ============

var api openapi.OpenAPI

func initBot(appID, appSecret string) error {
	creds := &token.QQBotCredentials{
		AppID:     appID,
		AppSecret: appSecret,
	}

	tokenSrc := token.NewQQBotTokenSource(creds)

	if err := token.StartRefreshAccessToken(context.Background(), tokenSrc); err != nil {
		return fmt.Errorf("start refresh token failed: %w", err)
	}

	api = botgo.NewOpenAPI(appID, tokenSrc).
		WithTimeout(5 * time.Second).
		SetDebug(true)

	log.Printf("Bot 初始化完成")
	return nil
}

// ============ 天气推送 ============

type WeatherScheduler struct {
	userID string
	city   string
	hour   int
	minute int
}

func NewWeatherScheduler(userID, city string, hour, minute int) *WeatherScheduler {
	return &WeatherScheduler{
		userID: userID,
		city:   city,
		hour:   hour,
		minute: minute,
	}
}

func (s *WeatherScheduler) Start() {
	go s.loop()
}

func (s *WeatherScheduler) loop() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, 0, 0, now.Location())

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		waitDuration := next.Sub(now)
		log.Printf("[天气推送] 下次推送时间: %s（还有 %s）", next.Format("2006-01-02 15:04:05"), waitDuration.Round(time.Second))

		time.Sleep(waitDuration)

		if err := s.push(); err != nil {
			log.Printf("[天气推送] 推送失败: %v", err)
		}
	}
}

func (s *WeatherScheduler) push() error {
	weather, err := fetchWeather(s.city)
	if err != nil {
		return fmt.Errorf("获取天气失败: %w", err)
	}

	content := fmt.Sprintf("☀️ 早上好！今日%s天气\n%s", s.city, weather)

	msg := &dto.MessageToCreate{
		Content: content,
		MsgType: 0,
	}

	if _, err := api.PostC2CMessage(context.Background(), s.userID, msg); err != nil {
		return fmt.Errorf("发送私信失败: %w", err)
	}

	log.Printf("[天气推送] 推送成功 → %s", s.userID)
	return nil
}

func fetchWeather(city string) (string, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=%%l:+%%C+%%t+💧%%h+🌬%%w&lang=zh", city)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ============ main ============

func main() {
	// 1. 加载配置
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化 Bot（鉴权）
	if err := initBot(cfg.QQ.AppID, cfg.QQ.Secret); err != nil {
		log.Fatalf("Bot 启动失败: %v", err)
	}

	// 3. 启动天气推送定时器
	ws := NewWeatherScheduler(
		cfg.Push.UserID,
		cfg.Push.City,
		cfg.Push.Hour,
		cfg.Push.Minute,
	)
	ws.Start()

	// 阻塞主线程
	select {}
}
