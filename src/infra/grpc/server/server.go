package server

import (
	"net"

	"github.com/LucasGois1/jarvis/src/application/usecases/chatcompletionstream"
	"github.com/LucasGois1/jarvis/src/infra/grpc/pb"
	"github.com/LucasGois1/jarvis/src/infra/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	ChatCompletionUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfig            chatcompletionstream.ChatCompletionUseCaseConfigDTO
	ChatService           service.ChatService
	Port                  string
	AuthToken             string
	StreamChannel         chan chatcompletionstream.ChatCompletionOutputDTO
}

func (g *GRPCServer) AuthInterceptor(service interface{}, serverStream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	context := serverStream.Context()
	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if token[0] != g.AuthToken {
		return status.Error(codes.Unauthenticated, "authorization token is invalid")
	}

	return handler(service, serverStream)
}

func NewGRPCServer(usecase chatcompletionstream.ChatCompletionUseCase, chatConfig chatcompletionstream.ChatCompletionUseCaseConfigDTO, port string, authToken string, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *GRPCServer {
	chatService := service.NewChatService(usecase, chatConfig, streamChannel)

	return &GRPCServer{
		ChatCompletionUseCase: usecase,
		ChatConfig:            chatConfig,
		ChatService:           *chatService,
		Port:                  port,
		AuthToken:             authToken,
		StreamChannel:         streamChannel,
	}
}

func (s *GRPCServer) Start() {
	grpcOptions := []grpc.ServerOption{
		grpc.StreamInterceptor(s.AuthInterceptor),
	}
	grpcServer := grpc.NewServer(grpcOptions...)
	pb.RegisterChatServiceServer(grpcServer, &s.ChatService)

	lis, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
