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
	COUNTRY_CODE  = "::COUNTRY_CODE::"
	DEFAULT_GREET = "Hello and welcome!"
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
	prompt := os.Getenv("PROMPT_GREET")

	if prompt == "" {
		return DEFAULT_GREET
	}

	prompt = strings.ReplaceAll(prompt, COUNTRY_CODE, countryCode)

	chatCompletion, err := p.Client.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		return DEFAULT_GREET
	}

	return chatCompletion.Choices[0].Message.Content
}

func (p *Prompt) Fact(c context.Context, countryCode string) string {
	prompt := os.Getenv("PROMPT_FACT")

	if prompt == "" {
		return ""
	}

	prompt = strings.ReplaceAll(prompt, COUNTRY_CODE, countryCode)

	chatCompletion, err := p.Client.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		return ""
	}

	return chatCompletion.Choices[0].Message.Content
}
