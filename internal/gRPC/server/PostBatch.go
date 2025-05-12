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

// PostBatch метод для множественного сохранения URL в базе
func (s *ProtoServer) PostBatch(ctx context.Context, in *pb.PostBatchRequset) (*pb.PostBatchResponse, error) {
	var response pb.PostBatchResponse
	// получение из контекста UserID и информации о регистрации
	UserInfo := ctx.Value(models.CtxKey).(models.UserInfo)
	UserID := UserInfo.UserID
	// получение из контекста UserID и информации о регистрации
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
	}

	InURLs := make(models.InBuff, len(in.BatchURL))
	var err error

	//заполнение вспомогательной из данных запроса
	for i, val := range in.BatchURL {
		InURLs[i].Hash = models.RandomString(8)
		InURLs[i].UserID = UserID
		InURLs[i].CorrID = val.CorrID
		InURLs[i].URL = val.URL

	}
	// сохранение URL в базу.
	OutBuff, err := s.Hd.Dt.SaveTx(InURLs, s.Hd.BaseAdr)
	if err != nil {
		logger.Log.Info("Error in SaveTx", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, `Ошибка при сохранении %s`, err.Error())
	}

	// сохранение URL в файл.
	err = s.Hd.Files.FileSaveTx(InURLs, s.Hd.BaseAdr)
	if err != nil {
		logger.Log.Info("Error in FileSaveTx", zap.Error(err))
	}
	// заполнение ответной структуры
	for _, val := range OutBuff {
		var OutPart pb.BatchURL
		OutPart.CorrID = val.CorrID
		OutPart.URL = val.ShortURL
		response.BatchURL = append(response.BatchURL, &OutPart)
	}

	return &response, nil
}
