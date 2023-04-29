package entities

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ChatConfig struct {
	Model            *Model
	Temperature      float32  // 0.0 to 1.0
	TopP             float32  // 0.0 to 1.0 - to a low value, like 0.1, the model will be very conservative in its word choices, and will tend to generate relatively predictable prompts
	N                int      // number of messages to generate
	Stop             []string // list of tokens to stop on
	MaxTokens        int      // number of tokens to generate
	PresencePenalty  float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	FrequencyPenalty float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, increasing the model's likelihood to talk about new topics.
}

type Status int

const (
	Ended Status = 1 << iota
	Active
	ValidStatus = Ended | Active
)

func NewStatus(status string) (Status, error) {
	switch status {
	case "Ended":
		return Ended, nil
	case "Active":
		return Active, nil
	default:
		return -1, errors.New("invalid status: " + status)
	}
}

type Chat struct {
	ID                   string
	UserID               string
	Status               Status
	TokenUsage           int
	Config               *ChatConfig
	InitialSystemMessage *Message
	Messages             []*Message
	HistoricalMessages   []*Message
}

func NewChat(UserID string, chatConfig *ChatConfig, initialSystemMessage *Message) (*Chat, error) {
	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               UserID,
		Status:               Active,
		Config:               chatConfig,
		InitialSystemMessage: initialSystemMessage,
	}

	chat.AddMessage(initialSystemMessage)

	if err := chat.validate(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *Chat) validate() error {
	if c.Status&ValidStatus == 0 {
		return errors.New("chat validation invalid status: " + fmt.Sprint(c.Status))
	}

	if c.Config == nil {
		return errors.New("chat validation invalid chat config: " + fmt.Sprint(c.Config))
	}

	if c.InitialSystemMessage == nil {
		return errors.New("chat validation invalid initial system message: " + fmt.Sprint(c.InitialSystemMessage))
	}

	return nil
}

func (c *Chat) AddMessage(message *Message) error {
	if c.Status == Ended {
		return errors.New("chat validation chat is closed")
	}

	for {
		if c.Config.Model.MaxTokens >= (message.Tokens + c.TokenUsage) {
			c.Messages = append(c.Messages, message)
			c.refreshTokenUsage()
			break
		}
		c.HistoricalMessages = append(c.HistoricalMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.refreshTokenUsage()
	}
	return nil
}

func (c *Chat) refreshTokenUsage() {
	var total int
	for _, message := range c.Messages {
		total += message.Tokens
	}

	c.TokenUsage = total
}

func (c *Chat) CountMessages() int {
	return len(c.Messages)
}

func (c *Chat) Close() {
	c.Status = Ended
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=Status
