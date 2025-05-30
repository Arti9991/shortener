package protoserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"regexp"

	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddr получение исходного URL по укороченному
func (s *ProtoServer) GetAddr(ctx context.Context, in *pb.GetAddrRequset) (*pb.GetAddrResponse, error) {
	var response pb.GetAddrResponse

	re := regexp.MustCompile(`^.*/`)
	// Заменяем найденное на пустоту
	ident := re.ReplaceAllString(in.ShortAddr, "")

	origURL, err := s.Hd.Dt.Get(ident)
	if err == models.ErrorDeleted {
		return nil, status.Errorf(codes.Unavailable, `URL was delted %s`, ident)
	} else if err != nil {
		return nil, status.Errorf(codes.Aborted, `Ошибка в базе данных %s`, err.Error())
	}
	// заполнение ответного короткого URL
	response.Addres = origURL

	return &response, nil
}
