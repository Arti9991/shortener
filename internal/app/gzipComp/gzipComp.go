package gzipComp

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Arti9991/shortener/internal/logger"
	"go.uber.org/zap"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	res  http.ResponseWriter
	zres *gzip.Writer
}

func newCompressWriter(res http.ResponseWriter) *compressWriter {
	return &compressWriter{
		res:  res,
		zres: gzip.NewWriter(res),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.res.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zres.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.res.Header().Set("Content-Encoding", "gzip")
	}
	c.res.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zres.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func MiddlewareGzip(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ores := res

		//проверка заголовка с типом контента, для дальнейшего сжатия
		contentType := req.Header.Get("Content-Type")
		acceptedType := (strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/json"))
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip && acceptedType {
			logger.Log.Info("Client has gzip header")
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cres := newCompressWriter(res)
			// меняем оригинальный http.ResponseWriter на новый
			ores = cres
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cres.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip && acceptedType {
			logger.Log.Info("Client send compressed data in gzip format")
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(req.Body)
			if err != nil {
				logger.Log.Info("Error in compress reader", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			// меняем тело запроса на новое
			req.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ores, req)
	}
}
