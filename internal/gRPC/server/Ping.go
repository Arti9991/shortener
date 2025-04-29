package protoServer

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"

	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddUser реализует интерфейс добавления пользователя.
func (s *ProtoServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse

	err := s.Hd.Dt.Ping()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, `Ошибка в базе данных %s`, err.Error())
	}

	return &response, nil
}
