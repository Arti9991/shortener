package storage

import (
	"github.com/Arti9991/shortener/internal/models"
)

type StorFunc interface {
	Save(key string, val string, UserID string) error
	Get(key string) (string, error)
	GetOrig(val string) (string, error)
	GetUser(UserID string, BaseAdr string) (models.UserBuff, error)
	SaveTx(InURLs models.InBuff, BaseAdr string) (models.OutBuff, error)
	Delete(keys []string, UserID string) error
	Ping() error
}
