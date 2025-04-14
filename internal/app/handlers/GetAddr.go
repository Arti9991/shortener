package handlers

import (
	"net/http"
	"path"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// GetAddr хэндлер для получения оригинального URL по укороченному.
func GetAddr(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// добавляем счетчик для graceful shutdown
		hd.Wg.Add(1)
		defer hd.Wg.Done()

		var err error
		// получаем индентификатор из URL запроса
		ident := path.Base(req.URL.String())
		// запрашиваем оригинальный URL и проверяем был ли он удален.
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
		// добавляем оригинальный URL в заголовок location.
		res.Header().Set("Location", redir)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
