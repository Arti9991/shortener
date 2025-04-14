package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// GetAddrUser хэндлер для получения всех оригинальных URL
// сохраненных пользователем.
func GetAddrUser(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// добавляем счетчик для graceful shutdown
		hd.Wg.Add(1)

		var err error
		// получение из контекста UserID и информации о регистрации.
		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID
		IsExist := UserInfo.Register
		// установка заголовка ответа для незарегистрированного пользователя.
		if !IsExist {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		// получение всех сокращенных URL для данного пользователя из базы или памяти.
		OutBuff, err := hd.Dt.GetUser(UserID, hd.BaseAdr)
		if err == models.ErrorNoUserURL {
			res.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
			logger.Log.Info("Error im GET method!", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// кодирование тела ответа
		out, err := json.Marshal(OutBuff)
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		hd.Wg.Done()
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}

}
