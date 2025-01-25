package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
)

// хэндлер для получения оригинального URL по укороченному
func DeleteAddr(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			logger.Log.Info("Only DELETE requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var err error

		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID
		IsExist := UserInfo.Register

		if !IsExist {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil || string(body) == "" {
			logger.Log.Info("Bad request body", zap.String("body", string(body)))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		//var DeleteURLURLs models.DeleteBuff
		var DeleteURLURLs []string
		err = json.Unmarshal(body, &DeleteURLURLs)
		if err != nil {
			logger.Log.Info("Bad request unmarshall", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, DelURL := range DeleteURLURLs {
			err := hd.Dt.Delete(DelURL, UserID)
			if err != nil {
				logger.Log.Info("Error in delete function", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		res.WriteHeader(http.StatusAccepted)
	}
}
