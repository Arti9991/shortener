package handlers

import (
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

func MainPage(dt *storage.Data, BaseAdr string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(req.Body)
		if err != nil || string(body) == "" {
			logger.Log.Info("Bad request body", zap.String("method", string(body)))
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
