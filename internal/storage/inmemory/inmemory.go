package inmemory

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/Arti9991/shortener/internal/storage/files"
)

type Data struct {
	storage.StorFunc
	sync.Mutex
	File      *files.FileData
	ShortUrls map[string]string
	UserKeys  map[string]string
}

// инициализация карты для хранения пар:
// ключ (сокращенный URL) - значение (исходный URL)
func NewData(file *files.FileData) *Data {
	dt := make(map[string]string)
	us := make(map[string]string)
	return &Data{File: file, ShortUrls: dt, UserKeys: us}
}

// добавление пары ключ (сокращенный URL) - значение (исходный URL)
func (d *Data) Save(key string, value string, UserID string) error {
	d.Lock()
	defer d.Unlock()
	_, ok := d.ShortUrls[key]
	if !ok {
		d.ShortUrls[key] = value
		d.UserKeys[UserID] = key
	}
	return nil
}

// получение оригнального URL по сокращенному
func (d *Data) Get(key string) (string, error) {
	d.Lock()
	defer d.Unlock()
	if d.ShortUrls[key] == "" {
		return "", errors.New("no such URL in memory")
	} else {
		return d.ShortUrls[key], nil
	}
}

// получение хэша по оригиналу URL
func (d *Data) GetOrig(val string) (string, error) {
	d.Lock()
	defer d.Unlock()
	for key, v := range d.ShortUrls {
		if v == val {
			return key, nil
		}
	}
	return "", errors.New("no such URL in map")
}

// множестевнное сохранение во внутреннюю память
func (d *Data) SaveTx(dec *json.Decoder, BaseAdr string) (models.OutBuff, error) {
	var IncomeURL models.BatchIncomeURL
	var OutBuff models.OutBuff

	for dec.More() {
		err := dec.Decode(&IncomeURL)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		hashStr := models.RandomString(8)

		// сохраняем URL в памяти
		_, ok := d.ShortUrls[hashStr]
		if !ok {
			d.ShortUrls[hashStr] = IncomeURL.URL
		}

		// сохранение URL в файле
		err = d.File.FileSave(hashStr, IncomeURL.URL)
		if err != nil {
			logger.Log.Info("Error in safe to File")
		}

		var OutURL models.BatchOutURL
		OutURL.ShortURL = BaseAdr + "/" + hashStr
		OutURL.CorrID = IncomeURL.CorrID

		OutBuff = append(OutBuff, OutURL)
	}
	return OutBuff, nil
}

// заглушка функции ping для реализации DuckType
func (d *Data) Ping() error {
	return nil
}

// сброс всех значений в карте
func (d *Data) ClearStor() {
	d.Lock()
	defer d.Unlock()
	for k := range d.ShortUrls {
		delete(d.ShortUrls, k)
	}
}
