package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// var QuerryPrepare = `INSERT INTO urls (id, hash_id, income_url)
// 	VALUES  (DEFAULT, $1, $2);`

// хэндлер создания укороченного URL в формате JSON
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
			Corr_ID string `json:"correlation_id"`
			URL     string `json:"url"`
		}{}
		OutURL := &struct {
			Corr_ID  string `json:"correlation_id"`
			ShortURL string `json:"short_url"`
		}{}
		var OutBuff []byte
		// var stmt *sql.Stmt
		// var err error

		// if !hd.DataBase.InFiles {
		// 	stmt, err = hd.DataBase.DB.Prepare(QuerryPrepare)
		// 	if err != nil {
		// 		logger.Log.Info("Error in DB prepare", zap.Error(err))
		// 		hd.DataBase.InFiles = true
		// 	}
		// }
		// defer stmt.Close()
		dec := json.NewDecoder(req.Body)
		for {
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
			// if !hd.DataBase.InFiles {
			// 	_, err := stmt.Exec(hashStr, IncomeURL.URL)
			// 	if err != nil {
			// 		logger.Log.Info("Error in DB Save", zap.Error(err))
			// 		hd.DataBase.InFiles = true
			// 	}
			// }

			OutURL.ShortURL = hd.BaseAdr + "/" + hashStr
			OutURL.Corr_ID = IncomeURL.Corr_ID

			tmp, err := json.Marshal(OutURL)
			if err != nil {
				logger.Log.Info("Wrong responce body", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			OutBuff = append(OutBuff, tmp...)
		}

		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(OutBuff)
	}
}
