package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func Ping(hd *handlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// if req.Method != http.MethodGet {
		// 	logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
		// 	res.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		db, err := sql.Open("pgx", hd.DBInfo)
		if err != nil {
			logger.Log.Info("Error in opening database!", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = db.Ping(); err != nil {
			logger.Log.Info("Error in ping database!", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
