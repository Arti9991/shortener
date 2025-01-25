package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
)

// хэндлер для получения оригинального URL по укороченному
func GetAddrUser(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
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
		OutBuff, err := hd.Dt.GetUser(UserID, hd.BaseAdr)
		//fmt.Println(err)
		if err == models.ErrorNoUserURL {
			res.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
			logger.Log.Info("Error im GET method!", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// } else if redir == "" {
		// 	logger.Log.Info("There is no such identifier!", zap.String("ID", ident))
		// 	res.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		// кодирование тела ответа
		out, err := json.Marshal(OutBuff)
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}

}
