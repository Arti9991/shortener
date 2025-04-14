package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// PostBatch для сохранения множества оригинальных URL
// и создания укороченных в формате JSON.
func PostBatch(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Header.Get("content-type") != "application/json" {
			logger.Log.Info("Bad content-type header with this path!", zap.String("header", req.Header.Get("content-type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// добавляем счетчик для graceful shutdown
		hd.Wg.Add(1)

		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var InURLs models.InBuff
		// декодирование тела запроса.
		err = json.Unmarshal(body, &InURLs)
		if err != nil {
			logger.Log.Info("Bad request unmarshall", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//заполнение вспомогательной структуры хэшами.
		for i := range InURLs {
			InURLs[i].Hash = models.RandomString(8)
			InURLs[i].UserID = UserID
		}
		// сохранение URL в базу.
		OutBuff, err := hd.Dt.SaveTx(InURLs, hd.BaseAdr)
		if err != nil {
			logger.Log.Info("Error in SaveTx", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// сохранение URL в файл.
		err = hd.Files.FileSaveTx(InURLs, hd.BaseAdr)
		if err != nil {
			logger.Log.Info("Error in FileSaveTx", zap.Error(err))
		}

		// кодирование тела ответа.
		out, err := json.Marshal(OutBuff)
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		hd.Wg.Done()

		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(out)
	}
}
