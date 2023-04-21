package chatcompletionstream

import (
	"context"
	"errors"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/LucasGois1/jarvis/src/domain/entities"
	"github.com/LucasGois1/jarvis/src/domain/gateway"
)

type ChatCompletionUseCaseConfigDTO struct {
	Model                string
	Temperature          float32  // 0.0 to 1.0
	TopP                 float32  // 0.0 to 1.0 - to a low value, like 0.1, the model will be very conservative in its word choices, and will tend to generate relatively predictable prompts
	N                    int      // number of messages to generate
	Stop                 []string // list of tokens to stop on
	MaxTokens            int      // number of tokens to generate
	PresencePenalty      float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	FrequencyPenalty     float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, increasing the model's likelihood to talk about new topics.
	InitialSystemMessage string
}

type ChatCompletionInputDTO struct {
	ChatID      string
	UserID      string
	UserMessage string
	Config      ChatCompletionUseCaseConfigDTO
}

type ChatCompletionOutputDTO struct {
	ChatID  string
	UserID  string
	Content string
}

type ChatCompletionUseCase struct {
	ChatGateway gateway.ChatGateway
	Client      *openai.Client
	Stream      chan ChatCompletionOutputDTO
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, client *openai.Client, stream chan ChatCompletionOutputDTO) *ChatCompletionUseCase {
	return &ChatCompletionUseCase{
		ChatGateway: chatGateway,
		Client:      client,
		Stream:      stream,
	}
}

func (uc *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDTO) (*ChatCompletionOutputDTO, error) {
	chat, err := uc.ChatGateway.FindByID(ctx, input.ChatID)
	if err != nil {
		if err.Error() == "chat not found" {
			// create new chat (entity)
			chat, err = createNewChat(input)
			if err != nil {
				return nil, errors.New("error creating new chat: " + err.Error())
			}
			// save on database***
			err = uc.ChatGateway.CreateChat(ctx, chat)
			if err != nil {
				return nil, errors.New("error persisting new chat: " + err.Error())
			}
		} else {
			return nil, errors.New("error fetching existing chat: " + err.Error())
		}
	}
	userMessage, err := entities.NewMessage(entities.User, input.UserMessage, chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating user message: " + err.Error())
	}
	err = chat.AddMessage(userMessage)
	if err != nil {
		return nil, errors.New("error adding new message: " + err.Error())
	}

	messages := []openai.ChatCompletionMessage{}
	for _, msg := range chat.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role.String(),
			Content: msg.Content,
		})
	}

	resp, err := uc.Client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:            chat.Config.Model.Name,
			Messages:         messages,
			MaxTokens:        chat.Config.MaxTokens,
			Temperature:      chat.Config.Temperature,
			TopP:             chat.Config.TopP,
			PresencePenalty:  chat.Config.PresencePenalty,
			FrequencyPenalty: chat.Config.FrequencyPenalty,
			Stop:             chat.Config.Stop,
			Stream:           true,
		},
	)
	if err != nil {
		return nil, errors.New("error creating chat completion: " + err.Error())
	}

	var fullResponse strings.Builder

	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.New("error streaming response: " + err.Error())
		}

		fullResponse.WriteString(response.Choices[0].Delta.Content)
		r := ChatCompletionOutputDTO{
			ChatID:  chat.ID,
			UserID:  input.UserID,
			Content: fullResponse.String(),
		}

		uc.Stream <- r
	}

	assistant, err := entities.NewMessage(entities.Assistant, fullResponse.String(), chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating assistant message: " + err.Error())
	}

	err = chat.AddMessage(assistant)
	if err != nil {
		return nil, errors.New("error adding new message: " + err.Error())
	}

	err = uc.ChatGateway.SaveChat(ctx, chat)
	if err != nil {
		return nil, errors.New("error saving chat: " + err.Error())
	}

	return &ChatCompletionOutputDTO{
		ChatID:  chat.ID,
		UserID:  input.UserID,
		Content: fullResponse.String(),
	}, nil
}
func createNewChat(input ChatCompletionInputDTO) (*entities.Chat, error) {
	model := entities.NewModel(input.Config.Model, input.Config.MaxTokens)

	chatConfig := &entities.ChatConfig{
		Temperature:      input.Config.Temperature,
		TopP:             input.Config.TopP,
		N:                input.Config.N,
		Stop:             input.Config.Stop,
		MaxTokens:        input.Config.MaxTokens,
		PresencePenalty:  input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
		Model:            model,
	}

	initialMessage, err := entities.NewMessage(entities.System, input.Config.InitialSystemMessage, model)
	if err != nil {
		return nil, errors.New("error creating initial message: " + err.Error())
	}

	chat, err := entities.NewChat(input.UserID, chatConfig, initialMessage)
	if err != nil {
		return nil, errors.New("error creating new chat: " + err.Error())
	}

	return chat, nil
}
