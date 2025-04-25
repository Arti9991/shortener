package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
)

// GetStats хэндлер для получения статистики по количеству сохраненных URL и пользователей
func GetStats(hd *HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// проверяем сохраненный IP на пустую строку
		if hd.SubIP == "" {
			logger.Log.Info("Requset with this path was turned off")
			res.WriteHeader(http.StatusForbidden)
			return
		}
		// парсим сохраненный IP
		var err error
		_, subnet, err := net.ParseCIDR(hd.SubIP)
		if err != nil {
			logger.Log.Info("Error in parsing trusted subnet", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// забираем полученный IP из хэдера
		strIP := req.Header.Get("X-Real-IP")
		inpIP := net.ParseIP(strIP)
		// проверяем подходит ли он к нашей подсети
		if !subnet.Contains(inpIP) {
			logger.Log.Info("Trusted subnet do not contain this ip", zap.String("IP", inpIP.String()))
			res.WriteHeader(http.StatusForbidden)
			return
		}
		// получаем статистику
		OutBuff, err := hd.Dt.Stats()
		if err != nil {
			logger.Log.Info("Error im Stats method!", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// кодирование тела ответа
		out, err := json.MarshalIndent(OutBuff, "", "  ")
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}
