package inmemory

import (
	"sync"
)

type Data struct {
	sync.Mutex
	ShortUrls map[string]string
}

// инициализация карты для хранения пар:
// ключ (сокращенный URL) - значение (исходный URL)
func NewData() *Data {
	dt := make(map[string]string)
	return &Data{ShortUrls: dt}
}

// добавление пары ключ (сокращенный URL) - значение (исходный URL)
func (d *Data) AddValue(key string, value string) {
	d.Lock()
	defer d.Unlock()
	_, ok := d.ShortUrls[key]
	if !ok {
		d.ShortUrls[key] = value
	}

}

// получение оригнального URL по сокращенному
func (d *Data) GetURL(key string) string {
	d.Lock()
	defer d.Unlock()
	return d.ShortUrls[key]
}

// сброс всех значений в карте
func (d *Data) ClearStor() {
	d.Lock()
	defer d.Unlock()
	for k := range d.ShortUrls {
		delete(d.ShortUrls, k)
	}
}
