package protoServer

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"sync"

	"github.com/Arti9991/shortener/internal/app/auth"
	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// DeleteAddr метод удаления адреса из базы
func (s *ProtoServer) DeleteAddr(ctx context.Context, in *pb.DeleteAddrRequest) (*pb.DeleteAddrResponse, error) {
	var response pb.DeleteAddrResponse

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
	// добавляем счетчик и запускаем горутину на отправку
	s.Hd.Wg.Add(1)
	SendDelete(s.Hd.Wg, in.Idents, UserID, s.Hd.OutDelCh)

	return &response, nil
}

// SendDelete функция отправки URL для удаления в горутине
func SendDelete(wg *sync.WaitGroup, URLs []string, UserID string, outCh chan models.DeleteURL) {

	go func() {
		defer wg.Done()
		var InURLs models.DeleteURL
		InURLs.ShortURL = URLs
		InURLs.UserID = UserID
		outCh <- InURLs
	}()
}
