package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Объявление zap логгера.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализация zap логгера (уровень логгирования INFO)
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zaplog, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zaplog
	return nil
}

// MiddlewareLogger обработчик для zap логгера с логированием полученных и отправленных запросов
func MiddlewareLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		reslog := loggingResponseWriter{
			ResponseWriter: res, //встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&reslog, req)
		duration := time.Since(start)
		Log.Info("got incoming HTTP request",
			zap.String("URI", req.RequestURI),
			zap.String("method", req.Method),
		)
		Log.Info("responce on request",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
			zap.Duration("duration", duration),
		)
	})
}

// переопределение методов write и WriteHeader для удобного использования middleware
type (
	// структура для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// реализация http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter //встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// Переопределение функции для интерфейса.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	//запись ответа, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// Переопределение функции для интерфейса.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	//запись кода статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode //код статуса
}
