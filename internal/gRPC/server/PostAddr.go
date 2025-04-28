package protoServer

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Arti9991/shortener/internal/app/auth"
	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// PostAddr
func (s *ProtoServer) PostAddr(ctx context.Context, in *pb.PostAddrRequset) (*pb.PostAddrResponse, error) {
	var response pb.PostAddrResponse
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
	//генерация рандомной строки
	hashStr := models.RandomString(8)

	if in.IsJSON {
		var IncomeURL models.IncomeURL

		err := json.NewDecoder(bytes.NewBuffer([]byte(in.Addres))).Decode(&IncomeURL)
		if err != nil {
			return nil, status.Errorf(codes.Aborted,
				`Выставлен флаг JSON. Ошибка в докодировании JSON %s`, err.Error())
		}
		in.Addres = IncomeURL.URL
	}

	err := s.Hd.Dt.Save(hashStr, in.Addres, UserID)
	if err != nil {
		logger.Log.Info("Error in Save", zap.Error(err))
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			logger.Log.Info("URL already exicts! Getting shorten version", zap.String("income URL", in.Addres))
			hashStr2, err2 := s.Hd.Dt.GetOrig(in.Addres)
			if err2 != nil {
				logger.Log.Info("Error in GetOrig", zap.Error(err2))
				return nil, status.Errorf(codes.NotFound, `Ошибка при поиске URL %s`, in.Addres)
			}

			ansStr2 := s.Hd.BaseAdr + "/" + hashStr2
			response.Addres = ansStr2
			return &response, nil
		} else {
			return nil, status.Errorf(codes.Aborted, `Ошибка в базе данных %s`, err.Error())
		}
	}

	// сохранение URL в файл
	err = s.Hd.Files.FileSave(hashStr, in.Addres)
	if err != nil {
		logger.Log.Info("Error in FileSave", zap.Error(err))
	}

	ansStr := s.Hd.BaseAdr + "/" + hashStr

	if in.IsJSON {
		var OutcomeURL models.OutcomeURL

		OutcomeURL.ShortURL = ansStr

		out, err := json.Marshal(OutcomeURL)
		if err != nil {
			return nil, status.Errorf(codes.Aborted,
				`Выставлен флаг JSON. Ошибка в кодировании JSON %s`, err.Error())
		}
		ansStr = string(out)
	}

	response.Addres = ansStr

	return &response, nil
}
