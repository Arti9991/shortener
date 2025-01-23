package files

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage"
	"go.uber.org/zap"
)

type FileStor struct {
	ID       int    `json:"uuid"`
	Shorturl string `json:"short_url"`
	Origurl  string `json:"original_url"`
}

type FileData struct {
	ID       int
	stor     storage.StorFunc
	Path     string
	InMemory bool //флаг для типа работы с памятью (файл или временная)
}

// конструктор структуры для работы с файлами. Также он создает/проверяет сам файл
func NewFiles(Path string, stor storage.StorFunc) (*FileData, error) {
	file, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil || Path == "" {
		return &FileData{InMemory: true}, err
	}
	file.Close()
	return &FileData{ID: 0, stor: stor, Path: Path, InMemory: false}, nil
}

// тестовый инциализатор с отключенным флагом
func FilesTest() *FileData {
	return &FileData{InMemory: true}
}

// // функция для чтения всех данных в файле и сохранения их в базу или память
// func (d *FileData) FileRead() error {
// 	// проверка флага на хранение данных в памяти
// 	if d.InMemory {
// 		return nil
// 	}
// 	var id int
// 	logger.Log.Info("INFO reading file")
// 	file, err := os.OpenFile(d.Path, os.O_RDONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		logger.Log.Info("Error in reading file! Setting in memory mode!", zap.Error(err))
// 		d.InMemory = true
// 		return err
// 	}
// 	defer file.Close()
// 	reader := bufio.NewReader(file)
// 	for {
// 		var fl FileStor
// 		buff, err := reader.ReadBytes('\n')
// 		if err == io.EOF {
// 			break
// 		} else if err != nil && err != io.EOF {
// 			logger.Log.Info("Error in reading data!", zap.Error(err))
// 			return err
// 		}
// 		err = json.Unmarshal(buff, &fl)
// 		if err != nil {
// 			logger.Log.Info("Error in unmarshalling data!", zap.Error(err))
// 			return err
// 		}
// 		d.stor.Save(fl.Shorturl, fl.Origurl)
// 		id = fl.ID
// 	}
// 	d.ID = id
// 	return nil
// }

// функция сохранения исходного и укороченного URL в файл
func (d *FileData) FileSave(key string, val string) error {
	// проверка флага на хранение данных в памяти
	if d.InMemory {
		return nil
	}
	logger.Log.Info("INFO Saving file")
	file, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Log.Info("Error in saving file! Setting in memory mode!", zap.Error(err))
		d.InMemory = true
		return err
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
	return nil
}

// функция сохранения множества URL в файл при чтении их из JSON
func (d *FileData) FileSaveTx(InURLs models.InBuff, BaseAdr string) error {
	// проверка флага на хранение данных в памяти
	if d.InMemory {
		return nil
	}

	logger.Log.Info("INFO Saving file")
	file, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Log.Info("Error in saving file! Setting in memory mode!", zap.Error(err))
		d.InMemory = true
		return err
	}
	writer := bufio.NewWriter(file)
	defer file.Close()

	for _, income := range InURLs {

		hashStr := income.Hash
		d.ID += 1
		var fl FileStor
		fl.ID = d.ID
		fl.Origurl = income.URL
		fl.Shorturl = hashStr

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
	return nil
}
