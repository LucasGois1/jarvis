package gateway

import (
	"context"

	"github.com/LucasGois1/jarvis/src/domain/entities"
)

type ChatGateway interface {
	CreateChat(ctx context.Context, chat *entities.Chat) error
	FindByID(ctx context.Context, id string) (*entities.Chat, error)
	SaveChat(ctx context.Context, chat *entities.Chat) error
}
