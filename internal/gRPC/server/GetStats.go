package protoServer

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"net"

	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetStats метод получения статистики для сервера
func (s *ProtoServer) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	var response pb.GetStatsResponse

	var err error
	var incomeIP string

	// проверяем сохраненный IP на пустую строку
	if s.Hd.SubIP == "" {
		return nil, status.Errorf(codes.PermissionDenied, `Метод недоступен`)
	}
	// парсим сохраненный IP
	_, subnet, err := net.ParseCIDR(s.Hd.SubIP)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, `Ошибка в парсинге доверенного IP`)
	}
	// получем информации об IP из метаданных
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			incomeIP = values[0]
		}
	} else {
		return nil, status.Errorf(codes.PermissionDenied, `Нет метаданных`)
	}
	inpIP := net.ParseIP(incomeIP)
	// проверяем подходит ли он к нашей подсети
	if !subnet.Contains(inpIP) {
		return nil, status.Errorf(codes.PermissionDenied, `Данный IP не входит в подсеть: %s`, incomeIP)
	}
	// получаем статистику
	OutBuff, err := s.Hd.Dt.Stats()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, `Ошибка в базе данных %s`, err.Error())
	}

	response.NumURLs = int64(OutBuff.NumUrls)
	response.NumUsers = int64(OutBuff.NumUsers)

	return &response, nil
}
