package storage

import (
	"sync"
)

type Data struct {
	sync.Mutex
	ShortUrls map[string]string
}

func NewData() *Data {
	dt := make(map[string]string)
	return &Data{ShortUrls: dt}
}

func (d *Data) AddValue(key string, value string) {
	d.Lock()
	defer d.Unlock()
	_, ok := d.ShortUrls[key]
	if !ok {
		d.ShortUrls[key] = value
	}

}

func (d *Data) GetURL(key string) string {
	d.Lock()
	defer d.Unlock()
	return d.ShortUrls[key]
}

func (d *Data) ClearStor() {
	d.Lock()
	defer d.Unlock()
	for k := range d.ShortUrls {
		delete(d.ShortUrls, k)
	}
}
