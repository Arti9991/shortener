package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

func randomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}

func PostAddr(dt *storage.Data, BaseAdr string) http.HandlerFunc {
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
		dt.AddValue(hashStr, string(body))

		ansStr := BaseAdr + "/" + hashStr

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(ansStr))
		//logger.Log.Info("Response status is 201 Created. Response body size", zap.Int("size", len([]byte(ansStr))))
	}
}

func GetAddr(dt *storage.Data) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ident := path.Base(req.URL.String())
		redir := dt.GetURL(ident)

		if redir == "" {
			logger.Log.Info("There is no such identifier!", zap.String("ID", ident))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", redir)
		res.WriteHeader(http.StatusTemporaryRedirect)
		//logger.Log.Info("Response status is 307 TemporaryRedirect.", zap.String("location", res.Header().Get("Location")), zap.Int("size", len(redir)))

	}
}

func PostAddrJSON(dt *storage.Data, BaseAdr string) http.HandlerFunc {
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
		OutURL := &struct {
			ShortUrl string `json:"result"`
		}{}

		err := json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		hashStr := randomString(8)
		dt.AddValue(hashStr, IncomeURL.URL)

		OutURL.ShortUrl = BaseAdr + "/" + hashStr

		out, err := json.Marshal(OutURL)
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(out)
		//logger.Log.Info("Response status is 201 Created. Response body size", zap.Int("size", len([]byte(ansStr))))
	}
}
