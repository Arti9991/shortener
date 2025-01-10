package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// хэндлер создания укороченного URL в формате JSON
func PostAddrJSON(hd *handlersData) http.HandlerFunc {
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
			URL string `json:"url"`
		}{}
		OutcomeURL := &struct {
			ShortURL string `json:"result"`
		}{}

		err := json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
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

		err = hd.DataBase.DBsave(hashStr, IncomeURL.URL)
		if err != nil {
			logger.Log.Info("Error in DBsave", zap.Error(err))
			if hd.DataBase.CodeIsUniqueViolation(err) {
				logger.Log.Info("URL already exicts! Getting shorten version")
				hashStr, err2 := hd.DataBase.DBgetOrig(IncomeURL.URL)
				if err2 != nil {
					logger.Log.Info("Error in GetOrig", zap.Error(err2))
				}
				OutcomeURL.ShortURL = hd.BaseAdr + "/" + hashStr

				out, err := json.Marshal(OutcomeURL)
				if err != nil {
					logger.Log.Info("Wrong responce body", zap.Error(err))
					res.WriteHeader(http.StatusBadRequest)
					return
				}

				res.Header().Set("content-type", "application/json")
				res.WriteHeader(http.StatusConflict)
				res.Write(out)
				return
			}
		}

		OutcomeURL.ShortURL = hd.BaseAdr + "/" + hashStr

		out, err := json.Marshal(OutcomeURL)
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
