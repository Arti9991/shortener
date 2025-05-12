package protoserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"

	"github.com/Arti9991/shortener/internal/app/auth"
	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetAddrUser метод получения всех URL, сохраненных пользователем
func (s *ProtoServer) GetAddrUser(ctx context.Context, in *pb.GetAddrUserRequset) (*pb.GetAddrUserResponse, error) {
	var response pb.GetAddrUserResponse

	// получение из контекста UserID и информации о регистрации
	UserInfo := ctx.Value(models.CtxKey).(models.UserInfo)
	UserID := UserInfo.UserID
	// если пользователь не зарегистрирован, выдаем ему ID
	if !UserInfo.Register {
		UserID = models.RandomString(16)

		UserEnc, err := auth.EncodeUserID(UserID)
		if err != nil {
			logger.Log.Info("Error in Encoding", zap.Error(err))
			UserEnc = ""
		}
		mdOut := metadata.New(map[string]string{})
		mdOut.Set("UserID", UserEnc)
		err = grpc.SetHeader(ctx, mdOut)
		if err != nil {
			logger.Log.Info("Error in setting header", zap.Error(err))
		}
		logger.Log.Info("This is a new user")
		return &response, nil
	}

	OutBuff, err := s.Hd.Dt.GetUser(UserID, s.Hd.BaseAdr)
	if err == models.ErrorNoUserURL {
		logger.Log.Info("This user has no URL", zap.Error(err))
		return &response, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Aborted, `Ошибка в базе данных %s`, err.Error())
	}
	// заполнение ответной структуры
	for _, val := range OutBuff {
		var Save pb.UserURLs
		Save.OrigURL = val.OrigURL
		Save.ShortURL = val.ShortURL
		response.UserURLs = append(response.UserURLs, &Save)
	}

	return &response, nil
}
