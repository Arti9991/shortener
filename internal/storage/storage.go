package storage

import (
	"encoding/json"

	"github.com/Arti9991/shortener/internal/models"
)

type StorFunc interface {
	Save(string, string) error
	Get(string) (string, error)
	GetOrig(string) (string, error)
	SaveTx(*json.Decoder, string) (models.OutBuff, error)
	Ping() error
}
