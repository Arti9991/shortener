package handlers

import (
	"github.com/Arti9991/shortener/internal/storage/database"
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
	dt       *inmemory.Data
	BaseAdr  string
	Files    *files.FileData
	DataBase *database.DBStor
}

// инциализация структуры с параметрами хэндлеров
func NewHandlersData(stor *inmemory.Data, base string, files *files.FileData, DataBase *database.DBStor) *handlersData {
	return &handlersData{dt: stor, BaseAdr: base, Files: files, DataBase: DataBase}
}
