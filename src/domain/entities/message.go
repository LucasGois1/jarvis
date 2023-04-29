package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	tiktoken_go "github.com/j178/tiktoken-go"
)

type Role int

const (
	User Role = 1 << iota
	System
	Assistant
	ValidRoles = User | System | Assistant
)

func NewRole(role string) (Role, error) {
	switch role {
	case "user":
		return User, nil
	case "system":
		return System, nil
	case "assistant":
		return Assistant, nil
	default:
		return -1, errors.New("invalid role: " + role)
	}
}

type Message struct {
	ID        string
	Role      Role
	Content   string
	Tokens    int
	Model     *Model
	CreatedAt time.Time
}

func NewMessage(role Role, content string, model *Model) (*Message, error) {
	totalTokens := tiktoken_go.CountTokens(model.Name, content)

	message := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Model:     model,
		Tokens:    totalTokens,
		CreatedAt: time.Now(),
	}

	if err := message.validate(); err != nil {
		return nil, err
	}

	return message, nil
}

func (m *Message) validate() error {
	if m.Role&ValidRoles == 0 {
		return errors.New("message validation invalid role: " + fmt.Sprint(m.Role))
	}

	if m.CreatedAt.IsZero() {
		return errors.New("message validation invalid created_at: " + m.CreatedAt.String())
	}

	if m.Model == nil {
		return errors.New("message validation invalid model: " + fmt.Sprint(m.Model))
	}

	if m.Tokens > m.Model.MaxTokens {
		return errors.New("message validation max tokens exceeded")
	}

	return nil
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=Role
