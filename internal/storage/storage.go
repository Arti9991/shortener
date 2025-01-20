package storage

import (
	"github.com/Arti9991/shortener/internal/models"
)

type StorFunc interface {
	Save(string, string) error
	Get(string) (string, error)
	GetOrig(string) (string, error)
	SaveTx(models.InBuff, string) (models.OutBuff, error)
	Ping() error
}
