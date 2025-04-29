package protoServer

import (
	// импортируем пакет со сгенерированными protobuf-файлами

	"context"
	"encoding/json"
	"fmt"

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

// PostAddr
func (s *ProtoServer) PostBatch(ctx context.Context, in *pb.PostBatchRequset) (*pb.PostBatchResponse, error) {
	var response pb.PostBatchResponse
	// получение из контекста UserID и информации о регистрации
	UserInfo := ctx.Value(models.CtxKey).(models.UserInfo)
	UserID := UserInfo.UserID

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
		fmt.Println(UserID)
	}

	var InURLs models.InBuff
	var err error
	// декодирование тела запроса.
	err = json.Unmarshal(in.Income, &InURLs)
	if err != nil {
		logger.Log.Info("Bad request unmarshall", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, `Ошибка в докодировании JSON %s`, err.Error())
	}
	//заполнение вспомогательной структуры хэшами.
	for i := range InURLs {
		InURLs[i].Hash = models.RandomString(8)
		InURLs[i].UserID = UserID
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

	// кодирование тела ответа.
	out, err := json.Marshal(OutBuff)
	if err != nil {
		logger.Log.Info("Wrong responce body", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, `Ошибка в кодировании JSON %s`, err.Error())
	}
	fmt.Println(string(out))
	response.Outcome = out

	return &response, nil
}
