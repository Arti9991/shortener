package handlers

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
)

// Ping хэндлер для проверки соединения с базой данных.
func Ping(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		err := hd.Dt.Ping()
		if err != nil {
			logger.Log.Info("Error in ping database!", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
