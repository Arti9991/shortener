package storage

type Data struct {
	ShortUrls map[string]string
}

func NewData() Data {
	dt := make(map[string]string)
	return Data{ShortUrls: dt}
}

func (d *Data) AddValue(key string, value string) {
	_, ok := d.ShortUrls[key]
	if !ok {
		d.ShortUrls[key] = value
	}
}

func (d *Data) GetURL(val string) string {
	for k, v := range d.ShortUrls {
		if v == val {
			//delete(d.ShortUrls, k)
			return k
		}
	}
	return ""
}
