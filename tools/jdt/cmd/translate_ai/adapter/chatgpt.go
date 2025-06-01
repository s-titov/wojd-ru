package adapter

import (
	"context"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

const (
	interval   = 50 * time.Millisecond
	openaiURL  = "https://api.openai.com/v1/chat/completions"
	model      = "gpt-4.1"
	systemText = "Ты профессиональный переводчик. Переводи текст китайской MMORPG строго с китайского на русский язык. Переводить только текст. Строго запрещено изменять цифры: любые числа (арабские цифры 0-9) должны остаться без изменений, в числовом виде. Не добавляй пояснений, форматирование не изменяй. Перевод должен быть в одну строку. После перевода в строке не должно остаться китайских иероглифов."
)

type Adapter interface {
	Translate(ctx context.Context, text string) (string, error)
}

type adapter struct {
	client *resty.Client
}

func NewAdapter() Adapter {
	client := resty.New()
	client.
		SetBaseURL(openaiURL).
		SetAuthToken(os.Getenv("TOKEN")). // Обязательно установи переменную окружения TOKEN
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetRateLimiter(rate.NewLimiter(rate.Every(interval), 10))

	return &adapter{
		client: client,
	}
}

// Translate переводит одну строку с китайского на русский
func (a *adapter) Translate(ctx context.Context, text string) (string, error) {
	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemText},
			{Role: "user", Content: text},
		},
	}

	var respBody ChatResponse
	resp, err := a.client.R().
		SetContext(ctx).
		SetBody(reqBody).
		SetResult(&respBody).
		Post("")
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", resp.Error().(error)
	}

	if len(respBody.Choices) == 0 {
		return "", err
	}

	return respBody.Choices[0].Message.Content, nil
}
