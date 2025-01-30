package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
)

// хэндлер для пометки URL как удаленного в базе данных
func DeleteAddr(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			logger.Log.Info("Only DELETE requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// получение из контекста UserID и информации о регистрации
		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID
		IsExist := UserInfo.Register
		// если пользователь не существует, устанавливается соответствующий статус
		if !IsExist {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		// чтение тела запроса с URL подлежащими удалению
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// функция для запуска горутины и отправки структуры в канал
		ThreadDecode(body, UserID, hd.OutDelCh)

		res.WriteHeader(http.StatusAccepted)
	}
}

// функция с горутиной, считывающей данные из тела запроса, декдоированием из JSON
// и отправки данных в канал
func ThreadDecode(body []byte, UserID string, outCh chan models.DeleteURL) {

	go func() {
		var InURLs models.DeleteURL
		// декодирование тела запроса
		err := json.Unmarshal(body, &InURLs.ShortURL)
		if err != nil {
			logger.Log.Info("Bad request unmarshall", zap.Error(err))
			return
		}
		InURLs.UserID = UserID
		outCh <- InURLs
	}()
}
