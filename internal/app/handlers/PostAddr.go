package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// PostAddr хэндлер для сохранения оригинального URL и создание укороченного.
func PostAddr(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil || string(body) == "" {
			logger.Log.Info("Bad request body", zap.String("body", string(body)))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// получение из контекста UserID и информации о регистрации
		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID

		//генерация рандомной строки
		hashStr := models.RandomString(8)

		// сохранение URL в базу или в память
		err = hd.Dt.Save(hashStr, string(body), UserID)
		if err != nil {
			logger.Log.Info("Error in Save", zap.Error(err))
			if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
				logger.Log.Info("URL already exicts! Getting shorten version", zap.String("income URL", string(body)))
				hashStr2, err2 := hd.Dt.GetOrig(string(body))
				if err2 != nil {
					logger.Log.Info("Error in GetOrig", zap.Error(err2))
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				ansStr := hd.BaseAdr + "/" + hashStr2

				res.Header().Set("content-type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte(ansStr))
				return
			}
		}

		// сохранение URL в файл
		err = hd.Files.FileSave(hashStr, string(body))
		if err != nil {
			logger.Log.Info("Error in FileSave", zap.Error(err))
		}

		ansStr := hd.BaseAdr + "/" + hashStr

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(ansStr))
	}
}
