package handlers

import (
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
	"golang.org/x/exp/rand"
)

func randomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}

type handlersData struct {
	dt      *inmemory.Data
	BaseAdr string
	Files   *files.FileData
	DBInfo  string
}

// инциализация структуры с параметрами хэндлеров
func NewHandlersData(stor *inmemory.Data, base string, files *files.FileData, DBAddress string) *handlersData {
	return &handlersData{dt: stor, BaseAdr: base, Files: files, DBInfo: DBAddress}
}

// // хэндлер создания укороченного URL
// func PostAddr(hd *handlersData) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodPost {
// 			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		body, err := io.ReadAll(req.Body)
// 		if err != nil || string(body) == "" {
// 			logger.Log.Info("Bad request body", zap.String("body", string(body)))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		hashStr := randomString(8)
// 		hd.dt.AddValue(hashStr, string(body))
// 		err = hd.Files.FileSave(hashStr, string(body))
// 		if err != nil {
// 			logger.Log.Info("Error in FileSave", zap.Error(err))
// 		}

// 		ansStr := hd.BaseAdr + "/" + hashStr

// 		res.Header().Set("content-type", "text/plain")
// 		res.WriteHeader(http.StatusCreated)
// 		res.Write([]byte(ansStr))
// 		//logger.Log.Info("Response status is 201 Created. Response body size", zap.Int("size", len([]byte(ansStr))))
// 	}
// }

// хэндлер для получения оригинального URL по укороченному
// func GetAddr(hd *handlersData) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodGet {
// 			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		ident := path.Base(req.URL.String())
// 		redir := hd.dt.GetURL(ident)

// 		if redir == "" {
// 			logger.Log.Info("There is no such identifier!", zap.String("ID", ident))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		res.Header().Set("Location", redir)
// 		res.WriteHeader(http.StatusTemporaryRedirect)
// 		//logger.Log.Info("Response status is 307 TemporaryRedirect.", zap.String("location", res.Header().Get("Location")), zap.Int("size", len(redir)))

// 	}
// }

// хэндлер создания укороченного URL в формате JSON
// func PostAddrJSON(hd *handlersData) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodPost {
// 			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if req.Header.Get("content-type") != "application/json" {
// 			logger.Log.Info("Bad content-type header with this path!", zap.String("header", req.Header.Get("content-type")))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		IncomeURL := &struct {
// 			URL string `json:"url"`
// 		}{}
// 		OutURL := &struct {
// 			ShortURL string `json:"result"`
// 		}{}

// 		err := json.NewDecoder(req.Body).Decode(&IncomeURL)
// 		if err != nil {
// 			logger.Log.Info("Bad request body", zap.Error(err))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		hashStr := randomString(8)
// 		hd.dt.AddValue(hashStr, IncomeURL.URL)
// 		err = hd.Files.FileSave(hashStr, IncomeURL.URL)
// 		if err != nil {
// 			logger.Log.Info("Error in FileSave", zap.Error(err))
// 		}

// 		OutURL.ShortURL = hd.BaseAdr + "/" + hashStr

// 		out, err := json.Marshal(OutURL)
// 		if err != nil {
// 			logger.Log.Info("Wrong responce body", zap.Error(err))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		res.Header().Set("content-type", "application/json")
// 		res.WriteHeader(http.StatusCreated)
// 		res.Write(out)
// 	}
// }

// func Ping(hd *handlersData) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodGet {
// 			logger.Log.Info("Only GET requests are allowed with this path!", zap.String("method", req.Method))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 	}
// }
