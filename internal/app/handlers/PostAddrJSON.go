package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// PostAddr хэндлер для сохранения оригинального URL и создание укороченного в формате JSON.
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

		UserInfo := req.Context().Value(models.CtxKey).(models.UserInfo)
		UserID := UserInfo.UserID
		//fmt.Println(UserID)
		//генерация рандомной строки
		hashStr := models.RandomString(8)

		// сохранение URL в память или в базу
		err = hd.Dt.Save(hashStr, IncomeURL.URL, UserID)
		if err != nil {
			logger.Log.Info("Error in Save", zap.Error(err))
			if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
				logger.Log.Info("URL already exicts! Getting shorten version", zap.String("income URL", IncomeURL.URL))
				hashStr2, err2 := hd.Dt.GetOrig(IncomeURL.URL)
				if err2 != nil {
					logger.Log.Info("Error in GetOrig", zap.Error(err2))
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				OutcomeURL.ShortURL = hd.BaseAdr + "/" + hashStr2

				out, err2 := json.Marshal(OutcomeURL)
				if err2 != nil {
					logger.Log.Info("Wrong responce body", zap.Error(err2))
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
