package handlers

import (
	"context"
	"sync"

	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/Arti9991/shortener/internal/storage/files"
)

// HandlersData структура со всей информацией для хэндлеров.
type HandlersData struct {
	Dt       storage.StorFunc
	Files    *files.FileData
	OutDelCh chan models.DeleteURL
	BaseAdr  string
	Ctx      context.Context
	Wg       *sync.WaitGroup
}

// NewHandlersData инциализация структуры с параметрами хэндлеров.
func NewHandlersData(stor storage.StorFunc,
	base string, files *files.FileData,
	OutDelCh chan models.DeleteURL,
	ctx context.Context,
	wg *sync.WaitGroup) *HandlersData {
	return &HandlersData{Dt: stor, BaseAdr: base, Files: files, OutDelCh: OutDelCh, Ctx: ctx, Wg: wg}
}
