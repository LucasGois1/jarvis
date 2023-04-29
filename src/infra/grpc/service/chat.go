package service

import (
	"github.com/LucasGois1/jarvis/src/application/usecases/chatcompletionstream"
	"github.com/LucasGois1/jarvis/src/infra/grpc/pb"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionUseCaseConfigDTO
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewChatService(usecase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionUseCaseConfigDTO, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: usecase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}
}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {

	input := chatcompletionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      c.ChatConfigStream,
	}

	ctx := stream.Context()

	go func() {
		for msg := range c.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}
	}()

	_, err := c.ChatCompletionStreamUseCase.Execute(ctx, input)

	if err != nil {
		return err
	}

	return nil

}
