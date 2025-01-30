package handlers

import (
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/Arti9991/shortener/internal/storage/files"
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

type HandlersData struct {
	Dt       storage.StorFunc
	BaseAdr  string
	Files    *files.FileData
	OutDelCh chan models.DeleteURL
}

// инциализация структуры с параметрами хэндлеров
func NewHandlersData(stor storage.StorFunc, base string, files *files.FileData, OutDelCh chan models.DeleteURL) *HandlersData {
	return &HandlersData{Dt: stor, BaseAdr: base, Files: files, OutDelCh: OutDelCh}
}
