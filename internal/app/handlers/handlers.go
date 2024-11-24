package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arti9991/shortener/internal/app/storage"
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

func MainPage(dt *storage.Data) http.HandlerFunc {
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
		fmt.Printf("\n\n\nBody: %s\t", string(body))

		ansStr := randomString(8)

		fmt.Printf("reqBody: %s\n\n\n", ansStr)
		fmt.Printf("reqURL + Body: %#v + %s\n\n\n", req.Host, ansStr)

		dt.AddValue(string(body), ansStr)

		ansStr = "http://" + req.Host + "/" + ansStr

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
		ident := req.URL.String()
		ident = strings.ReplaceAll(ident, "/", "")
		fmt.Printf("Id: %#v\t", ident)

		redir := dt.GetURL(ident)

		fmt.Printf("Redir: %#v\n", redir)

		if redir == "" {
			http.Error(res, "There is no such identifier!", http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", redir)
		res.WriteHeader(http.StatusTemporaryRedirect)
		body := "Data in =======================\n\r"
		body += fmt.Sprintf("Id: %#v\t", ident)
		body += fmt.Sprintf("Redir: %#v\n", redir)
		body += "Header responce:\n"
		for k, v := range res.Header() {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}
		res.Write([]byte(body))

	}
}