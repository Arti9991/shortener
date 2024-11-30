package handlers

import (
	"io"
	"net/http"
	"path"
	"time"

	"github.com/Arti9991/shortener/internal/storage"
	"golang.org/x/exp/rand"
)

func randomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(uint64(time.Now().UnixNano()))
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}

func MainPage(dt *storage.Data, BaseAdr string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(req.Body)
		if err != nil || string(body) == "" {
			http.Error(res, "The body is empty!", http.StatusBadRequest)
			return
		}

		hashStr := randomString(8)
		dt.AddValue(hashStr, string(body))

		ansStr := BaseAdr + "/" + hashStr

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(ansStr))
	}
}

func GetAddr(dt *storage.Data) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Only Get requests are allowed!", http.StatusBadRequest)
			return
		}

		ident := path.Base(req.URL.String())
		redir := dt.GetURL(ident)

		if redir == "" {
			http.Error(res, "There is no such identifier!", http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", redir)
		res.WriteHeader(http.StatusTemporaryRedirect)

	}
}
