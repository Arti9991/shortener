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
		//var err error

		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID
		IsExist := UserInfo.Register

		if !IsExist {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// body, err := io.ReadAll(req.Body)
		// if err != nil || string(body) == "" {
		// 	logger.Log.Info("Bad request body", zap.String("body", string(body)))
		// 	res.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		//var DeleteURLURLs models.DeleteBuff
		chDec, err := ThreadDecode(req.Body)
		if err != nil {
			logger.Log.Info("Error in ThreadDecode", zap.Error(err))
		}

		ThreadDelete(hd, UserID, chDec)

		res.WriteHeader(http.StatusAccepted)
	}
}

func ThreadDecode(Data io.ReadCloser) (chan string, error) {
	outCh := make(chan string)
	dec := json.NewDecoder(Data)
	// сдвигаем декодер
	_, err := dec.Token()
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(outCh)
		for dec.More() {
			var DeleteURLURLs string
			err := dec.Decode(&DeleteURLURLs) // (2)
			if err == io.EOF {
				return
			} else if err != nil {
				logger.Log.Info("Error in decode function", zap.Error(err))
				return
			}
			outCh <- DeleteURLURLs
		}
	}()
	return outCh, nil
}

func ThreadDelete(hd *HandlersData, UserID string, DelURL chan string) {
	for URL := range DelURL {
		go func(UsID string) {
			err := hd.Dt.Delete(URL, UserID)
			if err != nil {
				logger.Log.Info("Error in delete function", zap.Error(err))
				return
			}
		}(UserID)
	}
}
