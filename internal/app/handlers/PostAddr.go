package handlers

import (
	"io"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// хэндлер создания укороченного URL
func PostAddr(hd *handlersData) http.HandlerFunc {
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

		hashStr := randomString(8)
		hd.dt.AddValue(hashStr, string(body))
		err = hd.Files.FileSave(hashStr, string(body))
		if err != nil {
			logger.Log.Info("Error in FileSave", zap.Error(err))
		}

		ansStr := hd.BaseAdr + "/" + hashStr

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(ansStr))
		//logger.Log.Info("Response status is 201 Created. Response body size", zap.Int("size", len([]byte(ansStr))))
	}
}
