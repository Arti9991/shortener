package inmemory

import (
	"errors"
	"sync"

	"github.com/Arti9991/shortener/internal/models"
)

// структура с картой для хранения всех URL.
type Data struct {
	ShortUrls map[string]string
	UserKeys  map[string][]string
	sync.Mutex
}

// инициализация карты для хранения пар:
// ключ (сокращенный URL) - значение (исходный URL).
func NewData() *Data {
	dt := make(map[string]string)
	us := make(map[string][]string)
	return &Data{ShortUrls: dt, UserKeys: us}
}

// Save добавление пары ключ (сокращенный URL) - значение (исходный URL).
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

// Get получение оригнального URL по сокращенному.
func (d *Data) Get(key string) (string, error) {
	d.Lock()
	defer d.Unlock()
	if d.ShortUrls[key] == "" {
		return "", models.ErrorNoURL
	} else {
		return d.ShortUrls[key], nil
	}
}

// GetOrig получение хэша по оригиналу URL.
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

// GetUser получение всех сокращенных и оригинальных URL для конкретного пользователя.
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

// SaveTx множестевнное сохранение во внутреннюю память.
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

		var OutURL models.BatchOutURL
		OutURL.ShortURL = BaseAdr + "/" + hashStr
		OutURL.CorrID = income.CorrID

		OutBuff = append(OutBuff, OutURL)
	}
	return OutBuff, nil
}

// Ping заглушка функции ping для реализации DuckType.
func (d *Data) Ping() error {
	return nil
}

// Delete заглушка для реализации интерфейса хранилища.
func (d *Data) Delete(keys []string, UserID string) error {
	return errors.New("unable for inmemory mode")
}

// ClearStor сброс всех значений в карте
func (d *Data) ClearStor() {
	d.Lock()
	defer d.Unlock()
	for k := range d.ShortUrls {
		delete(d.ShortUrls, k)
	}
}

// CloseDB очищает встроенное хранилище.
func (d *Data) CloseDB() error {
	d.ClearStor()
	return nil
}

// Stats функция для получения количества сохраненых
// URL в памяти и количества пользователей
func (d *Data) Stats() (models.URLStats, error) {

	var stats models.URLStats
	stats.NumUrls = len(d.ShortUrls)
	stats.NumUsers = len(d.UserKeys)
	return stats, nil
}
