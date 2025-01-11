package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

type OutURL struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// хэндлер создания укороченных URL для массива JSON
func PostBatch(hd *handlersData) http.HandlerFunc {
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

		IncomeURL := &struct {
			CorrID string `json:"correlation_id"`
			URL    string `json:"original_url"`
		}{}
		var OutBuff []OutURL

		dec := json.NewDecoder(req.Body)
		if _, err := dec.Token(); err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		for dec.More() {
			err := dec.Decode(&IncomeURL)
			if err == io.EOF {
				break
			} else if err != nil {
				logger.Log.Info("Bad request body", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			hashStr := randomString(8)
			hd.dt.AddValue(hashStr, IncomeURL.URL)

			err = hd.Files.FileSave(hashStr, IncomeURL.URL)
			if err != nil {
				logger.Log.Info("Error in FileSave", zap.Error(err))
			}

			err = hd.DataBase.DBsaveTx(hashStr, IncomeURL.URL)
			if err != nil {
				logger.Log.Info("Error in DBsave", zap.Error(err))
			}

			var OutURL OutURL
			OutURL.ShortURL = hd.BaseAdr + "/" + hashStr
			OutURL.CorrID = IncomeURL.CorrID

			OutBuff = append(OutBuff, OutURL)
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
