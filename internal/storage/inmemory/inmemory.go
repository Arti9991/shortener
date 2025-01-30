package inmemory

import (
	"errors"
	"sync"

	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/Arti9991/shortener/internal/storage/files"
)

type Data struct {
	storage.StorFunc
	sync.Mutex
	File      *files.FileData
	ShortUrls map[string]string
	UserKeys  map[string][]string
}

// инициализация карты для хранения пар:
// ключ (сокращенный URL) - значение (исходный URL)
func NewData(file *files.FileData) *Data {
	dt := make(map[string]string)
	us := make(map[string][]string)
	return &Data{File: file, ShortUrls: dt, UserKeys: us}
}

// добавление пары ключ (сокращенный URL) - значение (исходный URL)
func (d *Data) Save(key string, value string, UserID string) error {
	d.Lock()
	defer d.Unlock()
	_, ok := d.ShortUrls[key]
	if !ok {
		d.ShortUrls[key] = value
		d.UserKeys[UserID] = append(d.UserKeys[UserID], key)
	}
	return nil
}

// получение оригнального URL по сокращенному
func (d *Data) Get(key string) (string, error) {
	d.Lock()
	defer d.Unlock()
	if d.ShortUrls[key] == "" {
		return "", models.ErrorNoURL
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
	return "", models.ErrorNoURL
}

// получение всех сокращенных и оригинальных URL для конкретного пользователя
func (d *Data) GetUser(UserID string, BaseAdr string) (models.UserBuff, error) {
	d.Lock()
	defer d.Unlock()
	var OutBuff models.UserBuff
	for _, hash := range d.UserKeys[UserID] {
		//var UserURL models.UserURL
		orig := d.ShortUrls[hash]
		short := BaseAdr + "/" + hash
		OutBuff = append(OutBuff, models.UserURL{ShortURL: short, OrigURL: orig})
	}
	if len(OutBuff) == 0 {
		return nil, models.ErrorNoUserURL
	}
	return OutBuff, nil
}

// множестевнное сохранение во внутреннюю память
func (d *Data) SaveTx(InURLs models.InBuff, BaseAdr string) (models.OutBuff, error) {

	var OutBuff models.OutBuff

	for _, income := range InURLs {
		hashStr := income.Hash

		// сохраняем URL в памяти
		_, ok := d.ShortUrls[hashStr]
		if !ok {
			d.ShortUrls[hashStr] = income.URL
			d.UserKeys[income.UserID] = append(d.UserKeys[income.UserID], hashStr)
		}

		// // сохранение URL в файле
		// err = d.File.FileSave(hashStr, IncomeURL.URL)
		// if err != nil {
		// 	logger.Log.Info("Error in safe to File")
		// }

		var OutURL models.BatchOutURL
		OutURL.ShortURL = BaseAdr + "/" + hashStr
		OutURL.CorrID = income.CorrID

		OutBuff = append(OutBuff, OutURL)
	}
	return OutBuff, nil
}

// заглушка функции ping для реализации DuckType
func (d *Data) Ping() error {
	return nil
}

// заглушка для реализации интерфейса хранилища
func (d *Data) Delete(keys []string, UserID string) error {
	return errors.New("unable for inmemory mode")
}

// сброс всех значений в карте
func (d *Data) ClearStor() {
	d.Lock()
	defer d.Unlock()
	for k := range d.ShortUrls {
		delete(d.ShortUrls, k)
	}
}
