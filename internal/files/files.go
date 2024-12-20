package files

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage"
	"go.uber.org/zap"
)

type FileStor struct {
	ID       int    `json:"uuid"`
	Shorturl string `json:"short_url"`
	Origurl  string `json:"original_url"`
}
type FileData struct {
	ID   int
	stor *storage.Data
	path string
}

func NewFiles(path string, stor *storage.Data) *FileData {
	return &FileData{ID: 0, stor: stor, path: path}
}

// функция для чтения всех данных в файле и сохранения их в карту
func (d *FileData) FileRead() {
	var id int
	logger.Log.Info("INFO reading file")
	file, err := os.OpenFile(d.path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Info("Error in creating file!", zap.Error(err))
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		var fl FileStor
		buff, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			logger.Log.Info("Error in reading data!", zap.Error(err))
			break
		}
		err = json.Unmarshal(buff, &fl)
		if err != nil {
			logger.Log.Info("Error in unmarshalling data!", zap.Error(err))
		}
		d.stor.AddValue(fl.Shorturl, fl.Origurl)
		id = fl.ID
	}
	d.ID = id
}

// функция сохранения исходного и укороченного URL в файл
func (d *FileData) FileSave(key string, val string) {
	logger.Log.Info("INFO Saving file")
	file, err := os.OpenFile(d.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Log.Info("Error in creating file!", zap.Error(err))
	}
	writer := bufio.NewWriter(file)
	defer file.Close()
	d.ID += 1
	var fl FileStor
	fl.ID = d.ID
	fl.Origurl = val
	fl.Shorturl = key
	fl.ID = d.ID
	data, err := json.Marshal(&fl)
	if err != nil {
		logger.Log.Info("Error in marshalling data!", zap.Error(err))
	}

	//запись строки
	if _, err := writer.Write(data); err != nil {
		logger.Log.Info("Error in writing data!", zap.Error(err))
	}
	//запись символа переноса строки
	if err := writer.WriteByte('\n'); err != nil {
		logger.Log.Info("Error in writing data!", zap.Error(err))
	}

	if err := writer.Flush(); err != nil {
		logger.Log.Info("Error in flashing data!", zap.Error(err))
	}
}
