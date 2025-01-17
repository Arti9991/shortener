package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

// хэндлер создания укороченного URL в формате JSON
func PostAddrJSON(hd *HandlersData) http.HandlerFunc {
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

		var IncomeURL models.IncomeURL
		var OutcomeURL models.OutcomeURL

		err := json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		hashStr := randomString(8)

		// сохранение URL в память или в базу
		err = hd.Dt.Save(hashStr, IncomeURL.URL)
		if err != nil {
			logger.Log.Info("Error in Save", zap.Error(err))
			if strings.Contains(fmt.Sprintf("%s", err), pgerrcode.UniqueViolation) {
				logger.Log.Info("URL already exicts! Getting shorten version", zap.String("income URL", IncomeURL.URL))
				hashStr, err2 := hd.Dt.GetOrig(IncomeURL.URL)
				if err2 != nil {
					logger.Log.Info("Error in GetOrig", zap.Error(err2))
					res.WriteHeader(http.StatusBadRequest)
					return
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

		// сохранение в файл
		err = hd.Files.FileSave(hashStr, IncomeURL.URL)
		if err != nil {
			logger.Log.Info("Error in FileSave", zap.Error(err))
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
