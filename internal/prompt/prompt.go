package prompt

import (
	"context"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	COUNTRY_CODE = "::COUNTRY_CODE::"
)

type Prompt struct {
	Client *openai.Client
}

func NewPrompt() *Prompt {
	_ = godotenv.Load()

	client := openai.NewClient(
		option.WithAPIKey("OPENAI_API_KEY"),
	)

	return &Prompt{
		Client: &client,
	}
}

func (p *Prompt) Greet(c context.Context, countryCode string) string {
	// get PROMPT_GREET from .env
	greet := os.Getenv("PROMPT_GREET")

	if greet == "" {
		return "Hello and welcome!"
	}

	greet = strings.ReplaceAll(greet, COUNTRY_CODE, countryCode)

	chatCompletion, err := p.Client.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(greet),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}
