package handlers

import (
	"net/http"
	"path"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// хэндлер для получения оригинального URL по укороченному
func GetAddr(hd *handlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ident := path.Base(req.URL.String())
		redir := hd.dt.GetURL(ident)

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
