package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// хэндлер создания укороченных URL для массива JSON
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

		dec := json.NewDecoder(req.Body)
		if _, err := dec.Token(); err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// сохранение URL в базу и в файл
		OutBuff, err := hd.Dt.SaveTx(dec, hd.BaseAdr)
		if err != nil {
			logger.Log.Info("Error in DBsaveTx", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := json.Marshal(OutBuff)
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(out)
	}
}
