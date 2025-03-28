package handlers

import (
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/Arti9991/shortener/internal/storage/files"
)

// HandlersData структура со всей информацией для хэндлеров.
type HandlersData struct {
	Dt       storage.StorFunc      // интерфейс хранилища
	BaseAdr  string                // базовый адрес возвращаемого URL
	Files    *files.FileData       // данные о файле хранения URL (если он есть)
	OutDelCh chan models.DeleteURL // канал для отправки URL подлежащих удалению
}

// NewHandlersData инциализация структуры с параметрами хэндлеров.
func NewHandlersData(stor storage.StorFunc, base string, files *files.FileData, OutDelCh chan models.DeleteURL) *HandlersData {
	return &HandlersData{Dt: stor, BaseAdr: base, Files: files, OutDelCh: OutDelCh}
}
