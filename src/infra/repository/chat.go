package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/LucasGois1/jarvis/src/domain/entities"
	"github.com/LucasGois1/jarvis/src/infra/db"
)

type ChatRepositoryMySQL struct {
	db      *sql.DB
	Queries db.Queries
}

func NewChatRepositoryMySQL(dbt *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		db:      dbt,
		Queries: *db.New(dbt),
	}
}

func (repo *ChatRepositoryMySQL) CreateChat(ctx context.Context, chat *entities.Chat) error {
	err := repo.Queries.CreateChat(
		ctx,
		db.CreateChatParams{
			ID:               chat.ID,
			UserID:           chat.UserID,
			InitialMessageID: chat.InitialSystemMessage.ID,
			Status:           chat.Status.String(),
			TokenUsage:       int32(chat.TokenUsage),
			Model:            chat.Config.Model.Name,
			ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),
			Temperature:      float64(chat.Config.Temperature),
			TopP:             float64(chat.Config.TopP),
			N:                int32(chat.Config.N),
			Stop:             chat.Config.Stop[0],
			MaxTokens:        int32(chat.Config.MaxTokens),
			PresencePenalty:  float64(chat.Config.PresencePenalty),
			FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	)

	if err != nil {
		return errors.New("error creating chat: " + err.Error())
	}

	err = repo.Queries.AddMessage(
		ctx,
		db.AddMessageParams{
			ID:        chat.InitialSystemMessage.ID,
			ChatID:    chat.ID,
			Content:   chat.InitialSystemMessage.Content,
			Role:      chat.InitialSystemMessage.Role.String(),
			Tokens:    int32(chat.InitialSystemMessage.Tokens),
			CreatedAt: chat.InitialSystemMessage.CreatedAt,
			Erased:    false,
		},
	)

	if err != nil {
		return errors.New("error creating initial message: " + err.Error())
	}

	return nil
}

func (repo *ChatRepositoryMySQL) FindByID(ctx context.Context, id string) (*entities.Chat, error) {
	chat := &entities.Chat{}
	res, err := repo.Queries.FindChatByID(ctx, id)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	status, err := entities.NewStatus(res.Status)
	if err != nil {
		return nil, err
	}

	chat.ID = res.ID
	chat.UserID = res.UserID
	chat.Status = status
	chat.TokenUsage = int(res.TokenUsage)
	chat.Config = &entities.ChatConfig{
		Model: &entities.Model{
			Name:      res.Model,
			MaxTokens: int(res.ModelMaxTokens),
		},
		Temperature:      float32(res.Temperature),
		TopP:             float32(res.TopP),
		N:                int(res.N),
		Stop:             []string{res.Stop},
		MaxTokens:        int(res.MaxTokens),
		PresencePenalty:  float32(res.PresencePenalty),
		FrequencyPenalty: float32(res.FrequencyPenalty),
	}

	messages, err := repo.Queries.FindMessagesByChatID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		role, err := entities.NewRole(message.Role)

		if err != nil {
			return nil, err
		}

		chat.Messages = append(chat.Messages, &entities.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      role,
			Tokens:    int(message.Tokens),
			Model:     &entities.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	erasedMessages, err := repo.Queries.FindErasedMessagesByChatID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, message := range erasedMessages {
		role, err := entities.NewRole(message.Role)

		if err != nil {
			return nil, err
		}

		chat.HistoricalMessages = append(chat.HistoricalMessages, &entities.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      role,
			Tokens:    int(message.Tokens),
			Model:     &entities.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}
	return chat, nil

}
func (repo *ChatRepositoryMySQL) SaveChat(ctx context.Context, chat *entities.Chat) error {
	params := db.SaveChatParams{
		ID:               chat.ID,
		UserID:           chat.UserID,
		Status:           chat.Status.String(),
		TokenUsage:       int32(chat.TokenUsage),
		Model:            chat.Config.Model.Name,
		ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),
		Temperature:      float64(chat.Config.Temperature),
		TopP:             float64(chat.Config.TopP),
		N:                int32(chat.Config.N),
		Stop:             chat.Config.Stop[0],
		MaxTokens:        int32(chat.Config.MaxTokens),
		PresencePenalty:  float64(chat.Config.PresencePenalty),
		FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
		UpdatedAt:        time.Now(),
	}

	err := repo.Queries.SaveChat(
		ctx,
		params,
	)
	if err != nil {
		return err
	}
	// delete messages
	err = repo.Queries.DeleteChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}
	// delete erased messages
	err = repo.Queries.DeleteErasedChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}
	// save messages
	i := 0
	for _, message := range chat.Messages {
		err = repo.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,
				ChatID:    chat.ID,
				Content:   message.Content,
				Role:      message.Role.String(),
				Tokens:    int32(message.Tokens),
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt,
				OrderMsg:  int32(i),
				Erased:    false,
			},
		)
		if err != nil {
			return err
		}
		i++
	}

	// save erased messages
	i = 0
	for _, message := range chat.HistoricalMessages {
		err = repo.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,
				ChatID:    chat.ID,
				Content:   message.Content,
				Role:      message.Role.String(),
				Tokens:    int32(message.Tokens),
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt,
				OrderMsg:  int32(i),
				Erased:    true,
			},
		)
		if err != nil {
			return err
		}
		i++
	}
	return nil
}
