package handlers

import (
	"net/http"
	"path"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// хэндлер для получения оригинального URL по укороченному
func GetAddr(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var err error

		ident := path.Base(req.URL.String())
		redir, err := hd.Dt.Get(ident)
		if err == models.ErrorDeleted {
			logger.Log.Info("URL was delted", zap.String("ID", ident))
			res.WriteHeader(http.StatusGone)
			return
		} else if err != nil {
			logger.Log.Info("Error im GET method!", zap.String("ID", ident), zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", redir)
		res.WriteHeader(http.StatusTemporaryRedirect)
		//logger.Log.Info("Response status is 307 TemporaryRedirect.", zap.String("location", res.Header().Get("Location")), zap.Int("size", len(redir)))

	}
}
